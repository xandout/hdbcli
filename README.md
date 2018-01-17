# HANA CLI

This is a work-in-progress replacement for SAP's `hdbsql`

[![Go Report Card](https://goreportcard.com/badge/github.com/xandout/hdbcli)](https://goreportcard.com/report/github.com/xandout/hdbcli)

## Current improvements

* Line editing, for example: moving the cursor left and right

* Command history

* Standard shell redirection

* Profiles for simple invocation

* Tab completion for shortcuts




## Usage

### Shell redirection

    hdbcli < test.sql
    
### Line editing

![Demo](gifs/demo.gif)

Just use your arrow keys :)

### Profile

This is currently loaded from `~/.hdbcli_config.json` and we will support command line flags in the future.

```json

{
  "hostname" : "your.hana.host.com",
  "port" : 30015,
  "username" : "USER",
  "password" : "SuperSecurePassword",
  "database" : "D00"
}
```


### Shortcuts

`/schemas` Shows all schemas in the current database.

`/describe TABLE_NAME` Describes TABLE_NAME.

`/exit` Exits the CLI

`/mode [csv|table]` Gets and sets the output mode


### Known Limitations

* LastInsertId and RowsAffected are always 0.  Maybe a bug upstream

* Some SELECT queries hang forever.  Only 2 packets return, maybe an upstream bug