package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	filename := "tmp.txt"

	setupFakeFile := func(t *testing.T, before, after, search, replace string) {
		err := ioutil.WriteFile(filename, []byte(before), 0644)
		if err != nil {
			fail(t, err)
		}

		result, err := read(filename, search, replace)
		if err != nil {
			fail(t, err)
		}

		got := strings.Join(result, "\n")
		if strings.Compare(got, after) != 0 {
			fmt.Println(result)
			fmt.Println(len(result))
			t.Errorf("got:\n%q\nexpected:\n%q\n", got, after)
		}

		if err := os.Remove(filename); err != nil {
			fail(t, err)
		}
	}

	for name, test := range readTests {
		t.Run(name, func(t *testing.T) {
			setupFakeFile(t, test.before, test.after, test.search, test.replace)
		})
	}
}

func fail(t *testing.T, err error) {
	t.Errorf("shouldn't fail here: %v", err)
}
