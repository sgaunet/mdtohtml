# Markdown to HTML cmd-line tool

If you have Go installed, install with:

    go get -u github.com/gomarkdown/mdtohtml

To run:

    mdtohtml [options] inputfile [outputfile]

Run `mdtohtml` to see all options.

This is also an example of how to use [gomarkdown/markdown](https://github.com/gomarkdown/markdown) library.

<div style = "display:block; clear:both; page-break-after:always;"></div>

# Forked project

Add the [Possibility to generate the html with the github design (-cssgh option)](https://github.com/sindresorhus/github-markdown-css)

The code will be improved in the future...

You can use this README to test the app.

```
mdtohtml -cssgh README.md README.html
```

<div style = "display:block; clear:both; page-break-after:always;"></div>

# Examples for the tests


## Table

Header1   | Header2              | Header3
--------- | -------------------- | --------------------
12        | with newline<br><span style="color:red">here</span> | Value
Value1    | Value2               | <span style="color:green">Value3</span>
Value1    | Value2               | Value3
Value1    | Value2               | Value3
Value1    | Value2               | Value3

## List 

**Here a list :**

* Point1
* <span style="color:blue">Point2</span>
    * Subpoint 1
    * Subpoint 2
    * SubPoint 3
    * Sub **Point** 4
* ~~Point3~~


## Image

![Example](img/Logo-Docker-.jpg)