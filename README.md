# zipred

[![Build Status](https://travis-ci.org/gofunky/zipred.svg)](https://travis-ci.org/gofunky/zipred)
[![GoDoc](https://godoc.org/github.com/gofunky/zipred?status.svg)](https://godoc.org/github.com/gofunky/zipred)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofunky/zipred)](https://goreportcard.com/report/github.com/gofunky/zipred)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/7664447e93c742219959e310a1d3f2d9)](https://www.codacy.com/app/gofunky/zipred?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gofunky/zipred&amp;utm_campaign=Badge_Grade)

ZIP file operations can get costly, especially for large files. This library allows you to filter and extract an online zip file on the fly.

In contrast to a conventional zip parser, it has the following benefits: 
* There is less latency since data is processed directly from the buffer on the fly.
* The download can be stopped once the metadata or target file has been found. Hence, less data is transferred.
* Irrelevant data is directly discarded without memory allocation.

This library gives you an efficient and idiomatic way for indexing zip files on the web.

For examples, check the corresponding folder.
