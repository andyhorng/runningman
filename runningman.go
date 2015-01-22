// Running man
package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func polling(update chan time.Time, db *sql.DB) {
	ticker := time.NewTicker(time.Second)
	var last time.Time

	for _ = range ticker.C {

		var t time.Time

		fmt.Println("Polling")
		err := db.QueryRow(`SELECT updated_datetime FROM table1 ORDER BY updated_datetime DESC`).Scan(&t)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if t.After(last) && !last.IsZero() {
			update <- t
		}

		last = t
	}

}

func main() {

	db, err := sql.Open("mysql", "test:test@tcp(192.168.59.103:3306)/test?parseTime=true")
	defer db.Close()

	if err != nil {
		fmt.Println(err)
	}

	// polling db
	update := make(chan time.Time)
	go polling(update, db)

	for last := range update {
		rows, err := db.Query("SELECT name FROM table1 WHERE updated_datetime >= ?", last)
		defer rows.Close()
		if err != nil {
			fmt.Println("Cannot get name")
			continue
		}

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				fmt.Println(err)
				continue
			}

			// processing
			time.Sleep(time.Second)
			fmt.Printf("Updated %s\n", name)
		}
	}
}
