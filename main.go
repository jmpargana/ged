package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	verbose      = flag.Bool("v", false, "output changes to user")
	regex        = flag.Bool("r", false, "search text is regex expandable")
	occurences   = flag.Int("o", -1, "occurences per line (by default all entries changes)")
	line         = flag.Int("l", -1, "change only on given line")
	lineRange    = flag.String("lr", ":", "change between range of lines")
	changedFiles = make(map[string][]changedLine)
	mutex        = new(sync.Mutex)
	start        = 0
	end          = -1
	re           *regexp.Regexp
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("\t%s [SEARCH] [REPLACE] file1 file2 dir1 ...\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	if err := run(); err != nil {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	args := loadArgs()

	if len(args) < 3 {
		return errors.New("you need to run this script with at least one file")
	}
	if err := parseLineRange(); err != nil {
		return err
	}
	search := args[0]
	replace := args[1]

	if *regex {
		re = regexp.MustCompile(search)
	}
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

func loadArgs() []string {
	args := flag.Args()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
	}
	return args
}

func splitCheck() ([]string, error) {
	lineRange := strings.Split(*lineRange, ":")
	if len(lineRange) != 2 {
		return nil, errors.New("give range like this: n:m")
	}
	return lineRange, nil
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
		if *line != -1 && *line != lineNum+1 || lineNum < start || end != -1 && end <= lineNum {
			contents = append(contents, s.Text())
			continue
		}
		t := ""

		if *regex {
			t = re.ReplaceAllString(s.Text(), replace)
		} else {
			t = strings.Replace(s.Text(), search, replace, *occurences)
		}
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
			fmt.Printf("%d:\t%s%s\t\t%s%s%s\n", line.line, RED, line.before, GREEN, line.after, BLACK)
		}
	}
}

func parseLineRange() (err error) {
	if *lineRange == ":" {
		return nil
	}

	lineRangeRunes := []rune(*lineRange)

	if lineRangeRunes[0] == rune(':') {
		lineRange, err := splitCheck()
		if err != nil {
			return err
		}
		end, err = strconv.Atoi(lineRange[1])
		return err
	} else if lineRangeRunes[len(lineRangeRunes)-1] == rune(':') {
		lineRange, err := splitCheck()
		if err != nil {
			return err
		}
		start, err = strconv.Atoi(lineRange[0])
		return err
	}

	lineRange, err := splitCheck()
	if err != nil {
		return err
	}

	start, err = strconv.Atoi(lineRange[0])
	if err != nil {
		return err
	}
	end, err = strconv.Atoi(lineRange[1])
	return err
}
