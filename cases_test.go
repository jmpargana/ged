package main

var readTests = map[string]struct {
	before, after, search, replace string
}{
	"empty text with flags":    {"", "", "something", "somethingelse"},
	"empty text with no flags": {"", "", "", ""},
	"no changes":               {"hello there", "hello there", "something", "no change"},
	"no changes with existing text": {
		before:  "hello there",
		after:   "hello there",
		search:  "something",
		replace: "no change",
	},
	"no changes multiple lines": {
		before:  "hello there\nand there",
		after:   "hello there\nand there",
		search:  "something",
		replace: "no change",
	},
	"no changes with empty flags": {
		before:  "hello there\nand there",
		after:   "hello there\nand there",
		search:  "",
		replace: "",
	},
	"preserve multiple empty lines at end with flags": {
		before:  "\nhello there\nand there\n\n",
		after:   "\nhello there\nand there\n",
		search:  "this",
		replace: "that",
	},
	"preserve empty line at end with flags": {
		before:  "\nhello there\nand there\n",
		after:   "\nhello there\nand there",
		search:  "this",
		replace: "that",
	},
	"preserve multiple empty lines at end with not flags": {
		before:  "\nhello there\nand there\n\n",
		after:   "\nhello there\nand there\n",
		search:  "",
		replace: "",
	},
	"preserve multiple empty lines in middle with not flags": {
		before:  "\nhello there\nand there\n\nand\n\n\nthat",
		after:   "\nhello there\nand there\n\nand\n\n\nthat",
		search:  "",
		replace: "",
	},
	"preserve multiple empty lines in middle with flags": {
		before:  "\nhello there\nand there\n\n\n\n\n\none",
		after:   "\nhello there\nand there\n\n\n\n\n\none",
		search:  "something",
		replace: "doesn't exist",
	},
	"preserve white spaces in beginning and end of line": {
		before:  "\n   hello there\nand there  \n\n\n   \n  \n\none",
		after:   "\n   hello there\nand there  \n\n\n   \n  \n\none",
		search:  "this",
		replace: "that",
	},
	"change word in beginning": {
		before:  "this is a test",
		after:   "that is a test",
		search:  "this",
		replace: "that",
	},
	"change word in middle": {
		before:  "test this is",
		after:   "test that is",
		search:  "this",
		replace: "that",
	},
	"change word in end": {
		before:  "test is this",
		after:   "test is that",
		search:  "this",
		replace: "that",
	},
	"change word in all lines": {
		before:  "test is this\nthis is test\ntest this is",
		after:   "test is that\nthat is test\ntest that is",
		search:  "this",
		replace: "that",
	},
	"change word in all lines where occurs": {
		before:  "test is this\nelse is test\ntest this is",
		after:   "test is that\nelse is test\ntest that is",
		search:  "this",
		replace: "that",
	},
	"change word with spaces": {
		before:  "test is this\nelse is test\ntest this is",
		after:   "test is this\nelse is test\ntestthatis",
		search:  " this ",
		replace: "that",
	},
	"change word with spaces multiple lines": {
		before:  "test is this\nelse is test\ntest this is",
		after:   "test isthat\nelse is test\ntestthat is",
		search:  " this",
		replace: "that",
	},
	"change multiple words": {
		before:  "hello and goodbye from out",
		after:   "hello and goodbye tuo morf",
		search:  "from out",
		replace: "tuo morf",
	},
	"match entire sentence": {
		before:  "hello \nworld\nword hello\nhello world",
		after:   "hello \nworld\nword hello\nthere",
		search:  "hello world",
		replace: "there",
	},
}
