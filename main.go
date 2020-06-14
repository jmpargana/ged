package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type color string

type changedLine struct {
	line          int
	before, after string
}

const (
	BLUE  color = "\u001b[34m"
	RED         = "\u001b[31m"
	GREEN       = "\u001b[32m"
	BLACK       = "\u001b[0m"
)

var (
	// all the flags come here
	verbose      = flag.Bool("v", true, "output changes to user")
	regex        = flag.Bool("r", false, "search text is regex expandable")
	matchPerLine = flag.Int("m", -1, "matches per line (by default all entries changes)")
	line         = flag.Int("l", -1, "change only on given line")
	lineRange    = flag.String("lr", ":", "change between range of lines")
	changedFiles = make(map[string][]changedLine)
	mutex        = new(sync.Mutex)
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	args := flag.Args()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
	}

	search := args[0]
	replace := args[1]

	wg := new(sync.WaitGroup)

	for _, filename := range args[2:] {
		wg.Add(1)
		go searchReplaceInFile(filename, search, replace, wg)
	}

	wg.Wait()

	if *verbose {
		printChanges()
	}

	return nil
}

func searchReplaceInFile(filename, search, replace string, wg *sync.WaitGroup) {

	files, err := fileOrDir(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v", err)
	}

	wg1 := new(sync.WaitGroup)

	for _, file := range files {
		wg1.Add(1)
		go readAndWrite(file, search, replace, wg1)
	}

	wg1.Wait()
	wg.Done()
}

func readAndWrite(file, search, replace string, wg1 *sync.WaitGroup) {

	contents, err := read(file, search, replace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v", err)
	}
	if err := write(file, contents); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v", err)
	}
	wg1.Done()
}

func fileOrDir(filename string) ([]string, error) {
	var files []string

	if err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func read(filename, search, replace string) ([]string, error) {
	var contents []string
	changedLines := make([]changedLine, 0) // maybe should allocate some space

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	for lineNum := 0; s.Scan(); lineNum++ {
		t := strings.Replace(s.Text(), search, replace, -1)
		contents = append(contents, t)

		if t != s.Text() && *verbose {
			changedLines = append(changedLines, changedLine{lineNum, s.Text(), t})
		}
	}
	file.Close()

	if len(changedLines) != 0 {
		mutex.Lock()
		changedFiles[filename] = changedLines
		mutex.Unlock()
	}

	return contents, nil
}

func write(filename string, contents []string) error {
	file, _ := os.Create(filename)
	w := bufio.NewWriter(file)

	for _, t := range contents {
		if _, err := w.WriteString(t); err != nil {
			return err
		}
		if err := w.WriteByte(byte('\n')); err != nil {
			return err
		}

		w.Flush()
	}
	file.Close()

	return nil
}

func printChanges() {
	for key, value := range changedFiles {
		fmt.Printf("\n%s%s%s\n", BLUE, key, BLACK)

		for _, line := range value {
			fmt.Printf("%d:\t%s%s\t\t\t%s%s%s\n", line.line, RED, line.before, GREEN, line.after, BLACK)
		}
	}
}
