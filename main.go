package main
import (
	"fmt"
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)
// Modified 2
type mapStringScan struct {
	colPtrs []interface{}
	row map[string]string
	colCount int
	colNames []string
}

func NewMapStringScan(columnNames []string) *mapStringScan {
	noCols := len(columnNames)
	structObj := &mapStringScan {
		colPtrs: make([]interface{}, noCols),
		row: make(map[string]string, noCols),
		colCount: noCols,
		colNames: columnNames,
	}

	for i := 0 ;i < noCols; i++ {
		structObj.colPtrs[i] = new(sql.RawBytes)
	}

	return structObj
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	err := rows.Scan(s.colPtrs...)
	if err != nil {
		return err
	}
	for i:= 0 ; i < s.colCount ; i++ {
		if rb, ok := s.colPtrs[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}

	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}

func readFromTable(db *sql.DB, table string) {
	log.Println("Extracting info from", table)
	queryString := "SELECT * FROM " + table
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err.Error())
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err.Error())
	}
	rc := NewMapStringScan(cols)
	for rows.Next() {
		err := rc.Update(rows)
		if err != nil {
			log.Fatal(err.Error())
		}
		cv := rc.Get()
		log.Printf("%#v \n\n", cv)
	}
} 


func displayTables(db *sql.DB) {
	row, err := db.Query("SELECT name FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY 1")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer row.Close()
	for row.Next() {
		var table string
		row.Scan(&table)
		log.Println("Table Name:", table)
	}
}

func main() {
	_, err := os.Stat("./database.sqlite")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("database.sqlite exists ...")

	// Open the sqlite file
	sqliteDatabase, _ := sql.Open("sqlite3", "./database.sqlite") 
	defer sqliteDatabase.Close()
	displayTables(sqliteDatabase)
	readFromTable(sqliteDatabase, "League")
}