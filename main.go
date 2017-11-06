package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	"os/user"

	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/xandout/hdbcli/config"
	"github.com/xandout/hdbcli/database"
	"github.com/xandout/hdbcli/hana-shortcuts"
)

var defPrompt = ">>> "
var multiPrompt = "... "

var special = hana_shortcuts.Commands

func handler(db *sql.DB, in string) {

	if strings.ToLower(in) == "help;" {
		fmt.Println("Help")
		for _, sc := range special {
			fmt.Printf("%s  \n\tUsage %s\n", sc.Name, sc.Help)
		}
		return
	}
	splitCommand := strings.SplitN(in, " ", 2)
	for _, sc := range special {

		if strings.HasPrefix(splitCommand[0], sc.Name) {
			if len(splitCommand) >= 2 {
				in = sc.Build(strings.Replace(splitCommand[1], ";", "", -1))

			} else {
				in = sc.Build()
			}

			log.Printf("Running %s\n", in)
			break
		}
	}

	// Looks like a query
	if strings.HasPrefix(strings.ToLower(in), "select ") {
		rows, qErr := db.Query(in)
		if qErr != nil {
			log.Printf("%v\n", qErr)
			return
		}

		printNoRowsErr := database.PrintRows(rows)
		if printNoRowsErr != nil {
			log.Printf("%v\n", printNoRowsErr)
		}

	} else {
		res, execErr := db.Exec(in)
		if execErr != nil {
			log.Printf("%v\n", execErr)
			return
		}
		ra, raErr := res.RowsAffected()
		li, liErr := res.LastInsertId()

		if raErr != nil {
			log.Printf("%v\n", raErr)
		}
		log.Printf("%d Rows Affected\n", ra)
		if liErr != nil {
			log.Printf("%v\n", liErr)
		}
		log.Printf("%v", li)
	}

}

func main() {

	u, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}
	conf, err := config.LoadConfiguration(filepath.Join(u.HomeDir, ".hdbcli_config.json"))

	if err != nil {
		log.Fatal(err)
	}
	dbErr := database.Initialize(*conf)

	if dbErr != nil {
		log.Println("dberr")
		log.Fatal(err.Error())
	}

	pingErr := database.DBCon.Ping()
	if pingErr != nil {
		log.Printf("%\n", pingErr)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 defPrompt,
		HistoryFile:            filepath.Join(u.HomeDir, ".hdbcli_history"),
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	//var inMulti bool

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)

		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt(multiPrompt)
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt(defPrompt)
		rl.SaveHistory(cmd)
		handler(database.DBCon, cmd)
	}
}
