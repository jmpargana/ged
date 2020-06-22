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
}
