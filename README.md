# Markdown to HTML cmd-line tool

If you have Go installed, install with:

    go get -u github.com/gomarkdown/mdtohtml

To run:

    mdtohtml [options] inputfile [outputfile]

Run `mdtohtml` to see all options.

This is also an example of how to use [gomarkdown/markdown](https://github.com/gomarkdown/markdown) library.

# Forked project

Added 2 functions to the project :

* [Possibility to generate the html with the github design (-cssgh option)](https://github.com/sindresorhus/github-markdown-css)
* Possibility to generate a PDF

The code will be improved in the future...

You can use this README to test the utility.

```
mdtohtml -cssgh README.md README.pdf
mdtohtml -cssgh README.md README.html
```

# Examples for the tests

Header1   | Header2              | Header3
--------- | -------------------- | --------------------
12        | with newline<br>here | Value
Value1    | Value2               | Value3
Value1    | Value2               | Value3
Value1    | Value2               | Value3
Value1    | Value2               | Value3

Here a list :

* Point1
* Point2
    * Subpoint 1
    * Subpoint 2
    * SubPoint 3
    * Sub **Point** 4
* ~~Point3~~


