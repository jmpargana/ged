# ged
[![GoDoc](https://godoc.org/github.com/jmpargana/ged?status.svg)](https://godoc.org/github.com/jmpargana/ged)
[![Build
Status](https://travis-ci.org/jmpargana/ged.svg?branch=master)](https://travis-ci.org/jmpargana/ged)
[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

ged is my version of the *sed* command line utility to find and replace text 
with a simpler syntax and prettier output that runs concurrently as is written
in go.

![](https://s3.eu-central-1.amazonaws.com/jmpargana.github.io/ged.gif)

With ged you can search and replace any number of occurrences in a line, you can give line ranges or just one line with the commands -o, -lr and -l respectfully.
You can also use expandable regular expression when using the flag -r and -v will give a colorful output of the changes made, file by file.

Currently when running the program with very large files (2M+ lines), the *sed*
program beats *ged* in terms of performance for about 1 second, but when running with multiple files concurrently (5+) they average out, and eventually *ged* becames a faster solution.

## Installation

```sh
go get github.com/jmpargana/ged
cd $GOPATH/src/github.com/jmpargana/ged
go install
```




## License


[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

- **[MIT License](https://opensource.org/licenses/mit-license.php)**
- Copyright 2020 © João Pargana
