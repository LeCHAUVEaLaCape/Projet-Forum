package config

import (
	"database/sql"
	"fmt"
	"strconv"
)

func UpdateInfoUser(database *sql.DB, input string, state string, id int) {
	tx, err := database.Begin()
	if err != nil {
		fmt.Println(err)
	}
	query := "UPDATE users SET " + state + " = ? WHERE id = " + strconv.Itoa(id)
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(input)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()
}
