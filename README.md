# ged

ged is my version of the *sed* command line utility to find and replace text 
with a simpler syntax and prettier output that runs concurrently as is written
in go.

![](https://s3.eu-central-1.amazonaws.com/jmpargana.github.io/ged.gif)

With ged you can search and replace any number of occurrences in a line, you can give line ranges or just one line with the commands -o, -lr and -l respectfully.
You can also use expandable regular expression when using the flag -r and -v will give a colorful output of the changes made, file by file.

Installation:

```sh
go get github.com/jmpargana/ged
cd $GOPATH/src/github.com/jmpargana/ged
go install
```
