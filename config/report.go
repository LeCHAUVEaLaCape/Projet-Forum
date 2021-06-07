package config

import (
	"database/sql"
	"log"
)

// add the report to the database
func Report(nameUser string, nameReporter string, reasonReport string) {
	// Open the database
	var nameReported string
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()
	//verify if the user is already in the table
	rows_double, err := database.Query("SELECT nameUser FROM report")
	CheckError(err)
	defer rows_double.Close()
	for rows_double.Next() {
		err := rows_double.Scan(&nameReported)
		if err != nil {
			log.Fatal(err)
		}
		if nameReported == nameUser {
			return
		}
	}

	// add the user to the table report
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO report (nameUser, reasonReport, answerReport, reported, nameReporter) VALUES (?,?,'','',?)")
	CheckError(err)
	_, err = stmt.Exec(nameUser, reasonReport, nameReporter)
	CheckError(err)

	tx.Commit()
}

// get all the report pending
func SelectReport(data_dashboard map[string]interface{}) {
	var slicereport [3]string
	var all_Report [][3]string
	// Open the database
	database_report, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database_report.Close()
	rows_report, err := database_report.Query("SELECT nameReporter, nameUser, reasonReport FROM report WHERE reported = ?", "")
	CheckError(err)
	defer rows_report.Close()
	for rows_report.Next() {
		err := rows_report.Scan(&slicereport[0], &slicereport[1], &slicereport[2])
		if err != nil {
			log.Fatal(err)
		}
		all_Report = append(all_Report, slicereport)

	}
	data_dashboard["report"] = all_Report
}

func DeleteUserFromReport(answerReport string, nameReported string, reportAccepted string) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// if the report is accepted from the admnin
	if reportAccepted == "1" {
		// delete everything from this user
		DelAccount(nameReported)
		// Update the db report
		tx, err := database.Begin()
		CheckError(err)
		stmt, err := tx.Prepare("UPDATE report SET answerReport = ?, reported = ? WHERE nameUser = ?")
		CheckError(err)
		_, err = stmt.Exec(answerReport, "reportAccepted", nameReported)
		CheckError(err)
		tx.Commit()
	}
}
