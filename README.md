erc-checktls
============

[![Build Status](https://travis-ci.org/keltia/erc-checktls.svg?branch=master)](https://travis-ci.org/keltia/erc-checktls)
[![GoDoc](http://godoc.org/github.com/keltia/erc-checktls?status.svg)](http://godoc.org/github.com/keltia/erc-checktls)

This is a small utility which will provide summary & diff-like operations for the reports generated by [ssllabs-scan](https://github.com/ssllabs/ssllabs-scan).

In addition the grade checked by [Imirhil](https://tls.imirhil.fr/) will be checked as well and displayed.

## Requirements

* Go >= 1.8
* jq (optional)

## Usage

SYNOPSIS
```
erc-checktls [-IvV] [-t csv|text] [-S site] <json file>
  
  -I	Do not fetch tls.imirhil.fr grade
  -S string
    	Display that site
  -V	More verbose mode
  -t string
    	Type of report (default "text")
  -v	Verbose mode
```

The json file needs to be generated by running `ssllabs-scan` post-processed by `jq` like this:
 
```
ssllabs-scan -hostfile <host list> | jq -j -c . | sort > <json file>
```

OPTIONS

| Option  | Default | Description|
| ------- |---------|------------|
| -I      | false   | Do not fetch tls.imirhil.fr grade |
| -S      | none    | Displays that site info only |
| -t      | text    | Output plain text or csv |
| -v      | false   | Be verbose |
| -V      | false   | More verbose: displays ciphers info |

## Using behind a web Proxy

Linux/Unix:
```
    export HTTP_PROXY=[http://]host[:port] (sh/bash/zsh)
    setenv HTTP_PROXY [http://]host[:port] (csh/tcsh)
```

Windows:
```
    set HTTP_PROXY=[http://]host[:port]
```

The rules of Go's `ProxyFromEnvironment` apply (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, lowercase variants allowed).

## License

The [BSD 2-Clause license][bsd].

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
