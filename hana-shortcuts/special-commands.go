package hana_shortcuts

import "fmt"

// Shortcut describes a macro to run tedious commands
type Shortcut struct {
	Name string
	Sql  string
	Help string
}

// Build returns a full SQL string based on the arguments passed.
func (s Shortcut) Build(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf(s.Sql, v[0])

	}
	return s.Sql
}

var Commands = []Shortcut{
	// This is roughly synonymous to the MySQL describe command
	Shortcut{
		Name: "describe",
		Sql:  "SELECT COLUMN_NAME,DATA_TYPE_NAME,LENGTH,IS_NULLABLE FROM TABLE_COLUMNS WHERE TABLE_NAME = '%s';",
		Help: "describe TABLE_NAME;  Describes TABLE_NAME",
	},
	// This will show all the schemas in the current database
	Shortcut{
		Name: "schemas",
		Sql:  "SELECT * FROM SCHEMAS;",
		Help: "schemas; Show all schemas in database.",
	},
}
