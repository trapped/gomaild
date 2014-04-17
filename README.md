gomaild
=======

[![Build Status](https://travis-ci.org/trapped/gomaild.svg?branch=master)](https://travis-ci.org/trapped/gomaild) [![GoDoc](https://godoc.org/github.com/trapped/gomaild?status.png)](https://godoc.org/github.com/trapped/gomaild)

A Go implementation of nowadays popular email protocols, to provide a simple yet secure and customizable email daemon.

You can obtain more specific documentation through the godoc command, or on the [godoc website](http://godoc.org/github.com/trapped/gomaild).

---

gomaild is still in development: it's not recommended to deploy it in production yet.

---

###Supported protocols

**NOTE**: Protocols might be considered "implemented" even if not completely - features described in their first RFC should be enough to fully support clients.

|	Protocol	|	Status                  	|
|	--------	|	-------------------------	|
|	SMTP		|	Implemented, to optimize	|
|	POP3		|	Implemented, to optimize	|
|	IMAP4		|	Not yet implemented     	|

####POP3

#####Supported commands

- CAPA
- APOP
- DELE
- LIST
- NOOP
- PASS
- QUIT
- RETR
- RSET
- STAT
- TOP
- UIDL
- USER

#####Supported authentication methods

- USER + PASS
- APOP

OAUTH/cookie-like authentication might get implemented.

####SMTP

#####Supported commands

- HELO
- EHLO
- MAIL (FROM)
- RCPT
- DATA
- STARTTLS
- AUTH (LOGIN PLAIN CRAM-MD5)
- QUIT
- NOOP
- RSET
- VRFY

Commands like TURN or ETRN might not get implemented (or, at least, not soon), since their use is very limited.

#####Supported authentication methods

- LOGIN
- PLAIN
- CRAM-MD5

Due to the (apparent) lack of documentation about the DIGEST-MD5 method, I haven't implemented it yet.
GSSAPI might not get implemented, since it's been standardized for the C language, not for the others: implementing it might cause conflicts with clients.
OAUTH/cookie-like authentication might get implemented.

#####Supported encryption(cipher)

- TLS

At the moment, only through the STARTTLS command. Full-connection encryption will come soon.

AES, RC4 (keys provided by the individual SMTP user) encryption is planned. TripleSec has been considered as well.

###Packages

As gomaild is built using the Go language, it is made of many small packages, the most useful (and portable) of them being the ones listed below.
Individual packages should be described in detail in their own README file, found in their folder.

| Name		                                                | Purpose	                                      |
| --------------------------------------------------------  | ----------------------------------------------- |
| `parsers/textual`	                                        | A text (telnet) parsing package.                |
| `parsers/config`	                                        | An extended JSON parsing package.               |
| `cipher`	                                                | A package wrapping most used cipher actions.    |
| [`rfc2822`](https://github.com/trapped/rfc2822)           | An email parsing package.                       |

###Configuration

gomaild reads and parses (using the `parsers/config` package) all the .conf files in its executable's root directory.
The configuration files' content is then stored in the `Configuration` variable of the `config` package.
The syntax is simple: JSON plus comments, which work as follows:
- lines starting with the hash (`#`) character will be ignored;
- text at the right side of an hash character will be ignored.
There's no escaping for the hash character yet.
