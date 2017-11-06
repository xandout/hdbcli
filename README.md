# HANA CLI

This is a work-in-progress replacement for SAP's `hdbsql`


## Current improvements

* Line editing, for example: moving the cursor left and right

* Command history

* Standard shell redirection

* Profiles for simple invocation


## Usage

### Shell redirection

    hdbcli < test.sql
    2017/11/06 16:41:02 no RowsAffected available after DDL statement
    2017/11/06 16:41:02 0 Rows Affected
    2017/11/06 16:41:02 no LastInsertId available after DDL statement
    2017/11/06 16:41:02 0
    2017/11/06 16:41:02 1 Rows Affected
    2017/11/06 16:41:02 no LastInsertId available
    2017/11/06 16:41:02 0
    
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

`schemas;` Shows all schemas in the current database.

`describe TABLE_NAME;` Describes TABLE_NAME.

