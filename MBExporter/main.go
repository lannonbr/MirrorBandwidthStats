package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type jsonData struct {
	Time string  `json:"time"`
	Rx   uint64  `json:"rx"`
	Tx   uint64  `json:"tx"`
	Rate float64 `json:"rate"`
}

type dujsonData struct {
	Distro string `json:"distro"`
	Bytes  uint64 `json:"bytes"`
}

var id int
var time string
var rx, tx uint64
var rate float64

func checkErr(str string, err error) {
	if err != nil {
		log.Fatalln(str, err)
	}
}

func exportHour() string {
	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	checkErr("Error: Failed opening database: ", err)

	rows, err := db.Query("SELECT * FROM hour ORDER BY id DESC LIMIT 24")
	checkErr("Error: query failed: ", err)

	var entries []jsonData

	for rows.Next() {
		err = rows.Scan(&id, &time, &rx, &tx, &rate)
		checkErr("Error: failed extracting data from row: ", err)

		newEntry := jsonData{Time: time, Rx: rx, Tx: tx, Rate: rate}
		entries = append(entries, newEntry)
	}

	jsonByteArr, err := json.Marshal(entries)
	checkErr("Error: Marshalling data failed: ", err)

	return "window.hourData = " + string(jsonByteArr) + ";\n"
}

func exportDay() string {
	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	checkErr("Error: Failed opening database: ", err)

	rows, err := db.Query("SELECT * FROM day ORDER BY id DESC LIMIT 7")
	checkErr("Error: query failed: ", err)

	var entries []jsonData

	for rows.Next() {
		err = rows.Scan(&id, &time, &rx, &tx, &rate)
		checkErr("Error: failed extracting data from row: ", err)

		newEntry := jsonData{Time: time, Rx: rx, Tx: tx, Rate: rate}
		entries = append(entries, newEntry)
	}

	jsonByteArr, err := json.Marshal(entries)
	checkErr("Error: Marshalling data failed: ", err)

	return "window.dayData = " + string(jsonByteArr) + ";\n"
}

func exportMonth() string {
	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	checkErr("Error: Failed opening database: ", err)

	rows, err := db.Query("SELECT * FROM month ORDER BY id DESC LIMIT 12")
	checkErr("Error: query failed: ", err)

	var entries []jsonData

	for rows.Next() {
		err = rows.Scan(&id, &time, &rx, &tx, &rate)
		checkErr("Error: failed extracting data from row: ", err)

		newEntry := jsonData{Time: time, Rx: rx, Tx: tx, Rate: rate}
		entries = append(entries, newEntry)
	}

	jsonByteArr, err := json.Marshal(entries)
	checkErr("Error: Marshalling data failed: ", err)

	return "window.monthData = " + string(jsonByteArr) + ";\n"
}

func exportTotal() string {
	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	checkErr("Error: Failed opening database: ", err)

	rows, err := db.Query("SELECT * FROM agg")
	checkErr("Error: query failed", err)

	var total int64
	var tot int64

	for rows.Next() {
		err = rows.Scan(&id, &time, &tot)
		checkErr("Error: Failed extracting data from row: ", err)

		total += tot
	}

	percentage := float64(total) / 1000000000000000.0

	return "window.pb = { total: " + strconv.FormatInt(total, 10) + ", percentage: " + strconv.FormatFloat(percentage, 'f', 5, 64) + "};\n"
}

func exportDistroUsage() string {
	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	checkErr("Error: Failed opening database: ", err)

	rows, err := db.Query("SELECT id, distro, bytes FROM distrousage ORDER BY id DESC LIMIT 41")
	checkErr("Error: Query failed", err)

	var entries []dujsonData

	var bytes uint64
	var distro string

	for rows.Next() {
		err = rows.Scan(&id, &distro, &bytes)
		checkErr("DU Error: Failed extracting data from row: ", err)

		newEntry := dujsonData{Distro: distro, Bytes: bytes}
		entries = append(entries, newEntry)
	}

	jsonByteArr, err := json.Marshal(entries)
	checkErr("Error: Marshalling data failed: ", err)

	return "window.distrousage = " + string(jsonByteArr) + ";\n"
}

func main() {
	hourStr := exportHour()
	dayStr := exportDay()
	monthStr := exportMonth()
	totalStr := exportTotal()
	distrousageStr := exportDistroUsage()

	file, _ := os.Create("./statsData.js")

	file.WriteString(hourStr)
	file.WriteString(dayStr)
	file.WriteString(monthStr)
	file.WriteString(totalStr)
	file.WriteString(distrousageStr)

	file.Sync()
	file.Close()
}
