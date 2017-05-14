package utils

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func WriteDBRecordsIntoMachineTable(runid int, subrunid int, machineip string, machinetype string, cloudtype string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	if err != nil {
		fmt.Print(err.Error())
	}
	var runidInString = strconv.Itoa(runid)
	var subrunidInString = strconv.Itoa(subrunid)
	var st = "insert into perf_readings(runid,subrun_id,machine_ip,machine_type,cloud_type) values(" + runidInString + "," + subrunidInString + ",\"" + machineip + "\",\"" + machinetype + "\",\"" + cloudtype + "\")"
	println(st)
	statement, err := db.Prepare(st)
	res, err := statement.Exec()
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 1 {
		fmt.Println("Succesfully inserted records")
	}
	db.Close()
}

func UpdateDBRecordsForCPUIntoMachineTable(runid int, subrunid int, cpuavg string, cpumax string, cpumin string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	if err != nil {
		fmt.Print(err.Error())
	}
	var runidInString = strconv.Itoa(runid)
	var subrunidInString = strconv.Itoa(subrunid)
	var st = "update perf_readings set cpu_avg = " + cpuavg + ",cpu_max=" + cpumax + ",cpu_min=" + cpumin + " where runid=" + runidInString + " and subrun_id=" + subrunidInString + ""
	fmt.Println(st)
	statement, err := db.Prepare(st)
	res, err := statement.Exec()
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 1 {
		fmt.Println("Succesfully inserted records")
	}
	db.Close()
}

func UpdateDBRecordsForMemoryIntoMachineTable(runid int, subrunid int, memoryavg string, memorymax string, memorymin string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	if err != nil {
		fmt.Print(err.Error())
	}
	var runidInString = strconv.Itoa(runid)
	var subrunidInString = strconv.Itoa(subrunid)
	var st = "update perf_readings set memory_avg = " + memoryavg + ",memory_max=" + memorymax + ",memory_min=" + memorymin + " where runid=" + runidInString + " and subrun_id=" + subrunidInString + ""
	fmt.Println(st)
	statement, err := db.Prepare(st)
	res, err := statement.Exec()
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 1 {
		fmt.Println("Succesfully inserted records")
	}
	db.Close()
}

func UpdateDBRecordsForLoadIntoMachineTable(runid int, subrunid int, loadavg string, loadmax string, loadsmin string) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	if err != nil {
		fmt.Print(err.Error())
	}
	var runidInString = strconv.Itoa(runid)
	var subrunidInString = strconv.Itoa(subrunid)
	var st = "update perf_readings set load_avg = " + loadavg + ",load_max=" + loadmax + ",load_min=" + loadsmin + " where runid=" + runidInString + " and subrun_id=" + subrunidInString + ""
	statement, err := db.Prepare(st)
	res, err := statement.Exec()
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 1 {
		fmt.Println("Succesfully inserted records")
	}
	db.Close()
}

func GetTotalRecordsCount() int {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/loadAggregator")
	var runid int
	if err != nil {
		fmt.Print(err.Error())
	}
	rows, err := db.Query("Select max(runid) from perf_readings")
	if rows.Next() {
		err := rows.Scan(&runid)
		fmt.Println(err, runid)
	}
	db.Close()
	return runid
}

func UpdateRecordsIntoDB(metricName string, runid int, subrunid int, value1 string, value2 string, value3 string) {

	switch metricName {
	case "CPUUtilization":
		UpdateDBRecordsForCPUIntoMachineTable(runid, subrunid, value1, value2, value3)
	case "MemoryUsedPercent":
		UpdateDBRecordsForMemoryIntoMachineTable(runid, subrunid, value1, value2, value3)
	case "CurrentLoad":
		UpdateDBRecordsForLoadIntoMachineTable(runid, subrunid, value1, value2, value3)
	}
}
