// Running man
package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func polling(update chan bool) {
	db, err := sql.Open("mysql", "")
	defer db.Close()
	if err != nil {
		panic("Cannot connect to db")
	}

	count := 0
	err = db.QueryRow(`SELECT COUNT(*) FROM table1
	WHERE updated_datetime > DATE_SUB(now(), INTERVAL 1 MINUTE)`).Scan(&count)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Got: %d\n", count)

	if count > 0 {
		// action
		update <- true
	}

}

func main() {
	// polling db
	ticker := time.NewTicker(time.Second)
	update := make(chan bool)
	go func() {
		for _ = range ticker.C {
			fmt.Println("Polling")
			polling(update)
		}
	}()

	// listen
	fmt.Println("Main")
	for _ = range update {
		fmt.Println("New restaurant")
	}
}
