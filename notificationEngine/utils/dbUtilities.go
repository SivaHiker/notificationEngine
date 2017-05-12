package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func writeDBRecordsIntoMachineTable(myvalue string, mykey string) {

}

func getTotalRecordsCount(mykey string, namespace string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	var count int
	if err != nil {
		fmt.Print(err.Error())
	}
	rows, err := db.Query("Select count(*) from perf_readings")
	if rows.Next() {
		err := rows.Scan(&count)
		fmt.Print(err, count)

	}
	fmt.Println(err)
}
