package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"database/sql"

	"os/user"

	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/xandout/hdbcli/config"
	"github.com/xandout/hdbcli/database"

	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/xandout/hdbcli/shortcuts"
)

var special = shortcuts.Commands
var printType = "table"

func tablePrinter(simpleRows *database.SimpleRows) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(simpleRows.Columns)
	table.AppendBulk(simpleRows.Rows)
	table.SetAlignment(tablewriter.ALIGN_RIGHT) // Set Alignment
	table.Render()

}

func csvPrinter(simpleRows *database.SimpleRows, printHeader bool) {
	if printHeader {
		w := csv.NewWriter(os.Stdout)

		header := [][]string{simpleRows.Columns}
		for _, record := range header {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		for _, record := range simpleRows.Rows {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		// Write any buffered data to the underlying writer (standard output).
		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	} else {
		w := csv.NewWriter(os.Stdout)

		for _, record := range simpleRows.Rows {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		// Write any buffered data to the underlying writer (standard output).
		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}
}

func handler(db *sql.DB, in string) {

	splitCommand := strings.SplitN(in, " ", 2)

	switch {
	case strings.ToLower(in) == "help;":
		fmt.Println("Help")
		for _, sc := range special {
			fmt.Printf("%s\n\tUsage %s\n", sc.Name, sc.Help)
		}
		return
	case strings.ToLower(in) == "/out csv;":
		printType = "csv"
		return
	case strings.ToLower(in) == "/out table;":
		printType = "table"
		return
	}
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

		converted, err := database.ConvertRows(rows)
		if err != nil {
			log.Fatal(err)
		}

		switch printType {
		case "table":
			tablePrinter(converted)
		case "csv":
			csvPrinter(converted, true)
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

func startREPL(u *user.User) {
	var defPrompt = ">>> "
	var multiPrompt = "... "
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 defPrompt,
		HistoryFile:            filepath.Join(u.HomeDir, ".hdbcli_history"),
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := rl.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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
		if err := rl.SaveHistory(cmd); err != nil {
			log.Println("Couldn't save history ", err)
		}
		handler(database.DBCon, cmd)

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
		log.Fatal(pingErr)
	}
	startREPL(u)

}
