gomaild
=======

A Go implementation of nowadays popular email protocols, to provide a simple yet secure and customizable email daemon.

---

gomaild is still at an early stage of development: it's not recommended to deploy it in production yet.

---

###Supported protocols

|	Protocol	|	Status	|
|	--------	|	------	|
|	SMTP		|	Already implemented in another project, awaits porting	|
|	POP3		|	Currently implementing	|
|	IMAP4		|	Not yet implemented	|

###Packages

As gomaild is built using the Go language, it is made of many small packages, the most useful (and portable) of them being the ones listed below.
Individual packages are described in detail in their own README file, found in their folder.

|	Name		|	Purpose	|
|	----		|	-------	|
|	`parsers/textual`	|	A text (commandline, textual protocol line-like) parsing package.	|

###Configuration

gomaild reads and parses (using the `parsers/textual` package, more on its use later) all the .conf files in its executable's root directory.
The configuration files' content is then stored in the `Settings` variable of the `config` package: the variable is a Go map (a key-value store), with strings (the various filenames without the .conf extension) as keys and maps as values; these maps have, again, strings as keys - but an array of interface{} (Go's unspecified value type, to allow for reflection casting) as values.
The configuration file syntax (which is also briefly described in the gomaild.conf file) is as follows:

- the hashtag character (`#`) opens and closes comment blocks; comment blocks are not closed until a second hashtag;
- the first argument (argument number 0), here referred to as "keyword", defines the key in the current files' second-level map for the whole statement;
- arguments are wrapped with the ````` character (ASCII backtick) if they require escaping of whitespace characters (as " ", tabs, newlines...);
- whitespace characters are trimmed away from the start and end of statements;
- statements end with the `;` character (ASCII semicolon); omitting it in a line results in a multiline statement;
- semicolons (`;`) are *prohibited* in arguments, as their occurrence would result in a premature end of the statement;

Statements such as user account definitions defined in the main file (`gomaild.conf`) are to be considered "master" or "general": servers will look in the `gomaild` map first, and then in their own.
*Note*: servers won't look in the main map for their specific settings (for example, the POP3 server will only check the existence of a user-specified greeting in the `pop3` map).
