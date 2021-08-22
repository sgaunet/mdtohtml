# Markdown to HTML cmd-line tool

If you have Go installed, install with:

    go get -u github.com/sgaunet/mdtohtml

To run:

    mdtohtml [options] inputfile [outputfile]

Run `mdtohtml` to see all options.

This is also an example of how to use [gomarkdown/markdown](https://github.com/gomarkdown/markdown) library.

# Forked project

I clean some code, remove some options and add [the github-markdown CSS](https://github.com/sindresorhus/github-markdown-css)

You can use this README to test the app.

```
mdtohtml README.md README.html
```


# Docker Image

Now there is a docker image to integrate the binary into your own docker image for example.

For example, the Dockerfile should look like :

```
FROM sgaunet/mdtohtml:0.3.1 AS mdtohtml

FROM <BASE-IMAGE:VERSION>
...
COPY --from=mdtohtml /usr/bin/mdtohtml /usr/bin/mdtohtml
...

```

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

## Page break

You can add page-break by adding this html code :

```
<div style = "display:block; clear:both; page-break-after:always;"></div>
```

It will be interpreted when using wkhtmltopdf to generate a PDF.

Check with README-with-page-break.md file.

```
mdtohtml ../README-with-page-break.md ../README-with-page.break.html
wkhtmltopdf ../README-with-page.break.html ../README-with-page.break.pdf
```