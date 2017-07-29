package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/lannonbr/MirrorBandwidthStats/models"
)

func cleanupBytes(str string) string {
	return strings.TrimSuffix(str, ".0")
}

func loadBandwidthCSV(filename string) []models.BandwidthEntry {
	csvfile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvfile))
	reader.Comment = '"'

	entries := []models.BandwidthEntry{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error reading csv line", err)
		}

		recv, _ := strconv.ParseUint(cleanupBytes(line[0]), 10, 64)
		send, _ := strconv.ParseUint(cleanupBytes(line[1]), 10, 64)

		entries = append(entries, models.BandwidthEntry{
			Recv:      recv,
			Send:      send,
			Timestamp: line[2],
		})
	}

	return entries
}

func humanizeBits(bits uint64) string {
	str := humanize.Bytes(bits)
	strArr := strings.Split(str, " ")
	str = strArr[0] + string([]rune(strArr[1])[0]) + "b"
	return str
}

func analyzeFile(filename string, fullPrint bool) (uint64, uint64, uint64, uint64) {
	entries := loadBandwidthCSV(filename)

	var totalRecv, totalSend uint64

	for _, entry := range entries {
		totalRecv += entry.Recv
		totalSend += entry.Send
	}

	totalOverall := totalRecv + totalSend

	rate := ((totalRecv + totalSend) * 8) / 3550

	if fullPrint {
		humanRecv := humanize.Bytes(totalRecv)
		humanSend := humanize.Bytes(totalSend)
		humanOverall := humanize.Bytes(totalOverall)

		fmt.Println("Filename:", filename)

		fmt.Println("Total Received:", humanRecv)
		fmt.Println("Total Sent:", humanSend)
		fmt.Println("Total overall:", humanOverall)
		fmt.Println("Rate:", humanizeBits(rate)+"/sec")
		fmt.Println("------")
	}

	return totalRecv, totalSend, totalOverall, rate
}

func prettyPrint(date string, totalRecv, totalSend, totalOverall, avgRate uint64) {
	fmt.Printf("Timestamp: %v\n", date)
	fmt.Println("------------")
	fmt.Println("Recieved:", humanize.Bytes(totalRecv))
	fmt.Println("Sent:", humanize.Bytes(totalSend))
	fmt.Println("Overall:", humanize.Bytes(totalOverall))
	fmt.Println("Rate:", humanizeBits(avgRate)+"/sec")
}

func csvPrint(date string, totalRecv, totalSend, totalOverall, avgRate uint64) {
	fmt.Printf("%s,%s,%s,%s,%s\n", date, humanize.Bytes(totalRecv), humanize.Bytes(totalSend), humanize.Bytes(totalOverall), humanizeBits(avgRate)+"/sec")
}

func csvPrintRaw(date string, totalRecv, totalSend, totalOverall, avgRate uint64) {
	fmt.Printf("%s,%d,%d,%d,%d", date, totalRecv, totalSend, totalOverall, avgRate)
}

func main() {

	format := os.Args[1]
	files := os.Args[2:]

	dateRegex, _ := regexp.Compile("([A-Z][a-z]{2,3})-([\\d]{1,2})-([\\d]{4})_([\\d]{2})")

	arr := dateRegex.FindStringSubmatch(files[0])

	month, day, year, hour := arr[1], arr[2], arr[3], arr[4]

	date := fmt.Sprintf("%s/%s/%s", month, day, year)
	dateWithHour := fmt.Sprintf("%s/%s/%s %s:00", month, day, year, hour)
	dateMonth := fmt.Sprintf("%s/%s", month, year)

	var counter int

	var totalRecv, totalSend, totalOverall, totalRate, avgRate uint64

	for _, arg := range files {
		recv, send, overall, rate := analyzeFile(arg, false)
		counter++

		totalRecv += recv
		totalSend += send
		totalOverall += overall

		totalRate += rate
	}

	avgRate = totalRate / uint64(counter)

	switch format {
	//Pretty Printed Formats
	case "pretty_month":
		prettyPrint(dateMonth, totalRecv, totalSend, totalOverall, avgRate)
	case "pretty_day":
		prettyPrint(date, totalRecv, totalSend, totalOverall, avgRate)
	case "pretty_hour":
		prettyPrint(dateWithHour, totalRecv, totalSend, totalOverall, avgRate)
	// CSV Formats
	case "csv_month":
		csvPrint(dateMonth, totalRecv, totalSend, totalOverall, avgRate)
	case "csv_day":
		csvPrint(date, totalRecv, totalSend, totalOverall, avgRate)
	case "csv_hour":
		csvPrint(dateWithHour, totalRecv, totalSend, totalOverall, avgRate)
	// Raw CSV (No humanized entries)
	case "csv_month_raw":
		csvPrintRaw(dateMonth, totalRecv, totalSend, totalOverall, avgRate)
	case "csv_day_raw":
		csvPrintRaw(date, totalRecv, totalSend, totalOverall, avgRate)
	case "csv_hour_raw":
		csvPrintRaw(dateWithHour, totalRecv, totalSend, totalOverall, avgRate)
	}

}
