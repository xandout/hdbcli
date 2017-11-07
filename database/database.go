package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"time"

	"errors"

	// Import the HDB driver
	_ "github.com/SAP/go-hdb/driver"
	"github.com/xandout/hdbcli/config"
)

var (
	// DBCon is the shared connection object
	DBCon *sql.DB
)

func getDbURL(configuration config.Configuration) string {

	dsn := &url.URL{
		Scheme: "hdb",
		User:   url.UserPassword(configuration.Username, configuration.Password),
		Host:   fmt.Sprintf("%s:%d", configuration.Hostname, configuration.Port),
	}
	return dsn.String()
}

// Initialize accepts a Configuration struct and builds a connection
func Initialize(configuration config.Configuration) error {

	var err error
	DBCon, err = sql.Open("hdb", getDbURL(configuration))

	if err != nil {
		return err
	}
	return DBCon.Ping()
}

// SimpleRows contains an array of columns and multi-dimensional array of rows in string format
type SimpleRows struct {
	Columns []string
	Rows    [][]string
	Length  int
}

// ConvertRows takes *sql.Rows and converts it to a SimpleRows
func ConvertRows(rows *sql.Rows) (simpleRows *SimpleRows, err error) {
	simpleRows = new(SimpleRows)
	columns, err := rows.Columns()
	colLen := len(columns)
	if err != nil {
		return simpleRows, err
	}
	if colLen > 0 {
		simpleRows.Columns = columns
	} else {
		return simpleRows, errors.New("got 0 columns")
	}
	values := make([]interface{}, colLen)
	scanArgs := make([]interface{}, colLen)
	if err != nil {
		return simpleRows, err
	}
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		colVals := make([]string, colLen)
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		for i, v := range values {
			switch t := v.(type) {
			case float64:
				colVals[i] = fmt.Sprintf("%f", t)
			case bool:
				colVals[i] = fmt.Sprintf("%t", t)
			case string:
				colVals[i] = fmt.Sprintf("%s", t)
			case int64:
				colVals[i] = fmt.Sprintf("%d", t)
			case nil:
				colVals[i] = "NULL"
			case []uint8:
				colVals[i] = fmt.Sprintf("%s", []byte(t[:]))
			case time.Time:
				colVals[i] = fmt.Sprintf("%s", t)
			default:
				colVals[i] = fmt.Sprintf("%v", t)

			}
		}
		simpleRows.Rows = append(simpleRows.Rows, colVals)
		simpleRows.Length++
	}
	return simpleRows, nil

}

// PrintRows is a convenience method to print $COL_NAME : $COL_VALUE  Very ugly but it works
func PrintRows(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(columns) > 0 {
		count := len(columns)
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		iter := 0

		for rows.Next() {
			err := rows.Scan(scanArgs...)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Row %d\n", iter)
			for i, v := range values {

				switch t := v.(type) {
				case float64:
					fmt.Printf("\t%s : %f\n", columns[i], t)
				case bool:
					fmt.Printf("\t%s : %t\n", columns[i], t)
				case string:
					fmt.Printf("\t%s : %s\n", columns[i], t)
				case int64:
					fmt.Printf("\t%s : %d\n", columns[i], t)
				case nil:
					fmt.Printf("\t%s : %s\n", columns[i], "NULL")
				case []uint8:
					fmt.Printf("\t%s : %s\n", columns[i], []byte(t[:]))
				case time.Time:
					fmt.Printf("\t%s : %v\n", columns[i], t)
				default:
					fmt.Printf("\t%s : TYPE::%T\n", columns[i], t)

				}

			}
			iter++
		}
		rowCountLine := fmt.Sprintf("Got %d Rows\n", iter)
		log.Println(rowCountLine)
	}
	return nil
}
