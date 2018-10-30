# classic_actions_to_html
Go program to write action descriptions (the content that displays when viewing an action) to HTML.

# Background
From time to time, Salsa's clients want to retrieve their actions for reference.  Unfortunately, that's something that Salas Classic can't do.

This app solves that difficulty by writing the action descriptions to HTML files.  Note that only the description is written.  The output does not contain

* Salsa templates,
* Letters or petition text, or
* Targets.

# Domain fixes
Salsa used to have a domain named `democracyinaction.org`.  That domain was turned off in favor of using `salsalabs.com`.

Clients that uploaded and used images and files when `democracyinaction.org` was alive still have actions that reference that domain.


This app solves that problem by automatically modifying URLs for the old domain(s) to `salsalabs.com`.

## Login credentials

The `classic_actions_to_html` application looks for your login credentials in a YAML file.  You provide the filename as part of the execution.

  The easiest way to get started is to  copy the `sample_login.yaml` file and edit it.  Here's an example.
```yaml
host: salsa4.salalabs.com
email: chuck@chew.cheese
password: extra-super-secret-password!
```
The `email` and `password` are the ones that you normally use to log in. The `host` can be found by using [this page](https://help.salsalabs.com/hc/en-us/articles/115000341773-Salsa-Application-Program-Interface-API-#api_host) in Salsa's documentation.

Save the new login YAML file to disk.  We'll need it when we  run the `classic_actions_to_html` app.

# Pre requisites

1. A [recent version of Go](https://golang.org/doc/install) installed.  You should use the googles
if you are installing on Windows.  Really. Trust me on this.

1. The correct Go directory structure.  Believe it or not, this is _very_ inportant.  Here's a sample.

```text
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```

## Installing `classic_acdtions_to_html`
```bash
go get "github.com/salsalabs/classic_actions_to_html"
go install "github.com/salsalabs/classic_actions_to_html"
```

# Usage
```text
go run main.go --help

usage: classic_actions_to_html --login=LOGIN [<flags>]

A command-line app to read actions, correct DIA URLs and write contents as HTML.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with login credentials
  --summary      Show action dates, keys and titles. Does not write HTML
```
# Output

Use `--summary` in the command line returns a list of actions.  
```text
2017-05-18 - 24864 - SALSA TEST - 180644.html
2017-05-19 - 24866 - Tell Congress Breastfeeding Families Need the FAM Act!.html
2017-05-23 - 24887 - Tell Congress to Support ALL Breastfeeding & Working Moms!.html
2017-06-08 - 24952 - Tell Congress to Support Breastfeeding Moms!.html
2018-08-13 - 26434 - Tell Us About Your Experience with Breastfeeding and Emergencies.html
```

Leave `--summary` off to create HTML.  The ouput appears in the `html` directory.  Each
action is stored in a file.  The filename provides quite a lot of information about the
action.

Pleaes be reminded that the output does not contain

* Salsa templates,
* Letters or petition text, or
* Targets.

Here's a sample.

```html

<!DOCTYPE html>
<html>
  <head>
    <title><center> SALSA TEST - 180644 </center></title>
  </head>
  <body>
    <div>
      <h1><center> SALSA TEST - 180644 </center></title>
    </div>
    <div>
      <p>This is a test action by SalsaLabs.</p>
    </div>
  </body>
</html>
```

# Questions?  Comments?
Use the [Issues link](https://github.com/salsalabs/classic_actions_to_html/issues) in the repository.  Don't waste your time by contacting Salsa support.
