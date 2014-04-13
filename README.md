gomaild
=======

[![Build Status](https://travis-ci.org/trapped/gomaild.svg?branch=master)](https://travis-ci.org/trapped/gomaild)

A Go implementation of nowadays popular email protocols, to provide a simple yet secure and customizable email daemon.

---

gomaild is still at an early stage of development: it's not recommended to deploy it in production yet.

---

###Supported protocols

**NOTE**: Protocols might be considered "implemented" even if not completely - features described in their first RFC should be enough to fully support clients.

|	Protocol	|	Status                  	|
|	--------	|	-------------------------	|
|	SMTP		|	Implemented, to optimize	|
|	POP3		|	Implemented, to optimize	|
|	IMAP4		|	Not yet implemented     	|

###Packages

As gomaild is built using the Go language, it is made of many small packages, the most useful (and portable) of them being the ones listed below.
Individual packages should be described in detail in their own README file, found in their folder.

| Name		                                                | Purpose	                          |
| --------------------------------------------------------  | --------------------------------    |
| `parsers/textual`	                                        | A text (telnet) parsing package.    |
| [`rfc2822`](https://github.com/trapped/rfc2822)           | An email parsing package.           |

###Configuration

gomaild reads and parses (using the `parsers/config` package) all the .conf files in its executable's root directory.
The configuration files' content is then stored in the `Configuration` variable of the `config` package.
The syntax is simple: JSON plus comments, which work as follows:
- lines starting with the hash (`#`) character will be ignored;
- text at the right side of an hash character will be ignored.
There's no escaping for the hash character yet.
