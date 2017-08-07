package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type jsonData struct {
	Time string  `json:"time"`
	Rx   uint64  `json:"rx"`
	Tx   uint64  `json:"tx"`
	Rate float64 `json:"rate"`
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

func exportHour() {
	db, err := sql.Open("sqlite3", "../mirrorband.sqlite")
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

	file, _ := os.Create("./hour.js")

	file.WriteString("window.hourData = ")
	file.WriteString(string(jsonByteArr))
	file.WriteString(";\n")

	file.Sync()
	file.Close()
}

func exportDay() {
	db, err := sql.Open("sqlite3", "../mirrorband.sqlite")
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

	file, _ := os.Create("./day.js")

	file.WriteString("window.dayData = ")
	file.WriteString(string(jsonByteArr))
	file.WriteString(";\n")

	file.Sync()
	file.Close()
}

func exportMonth() {
	db, err := sql.Open("sqlite3", "../mirrorband.sqlite")
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

	file, _ := os.Create("./month.js")

	file.WriteString("window.monthData = ")
	file.WriteString(string(jsonByteArr))
	file.WriteString(";\n")

	file.Sync()
	file.Close()
}

func main() {
	exportHour()
	exportDay()
	exportMonth()
}
