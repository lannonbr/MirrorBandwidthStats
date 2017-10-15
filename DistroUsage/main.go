package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/uniplaces/carbon"

	_ "github.com/mattn/go-sqlite3"
)

func getYesterday() string {
	yesterday := carbon.Now().SubDay()

	str := yesterday.FormattedDateString()

	str = strings.Replace(str, " ", "/", -1)
	str = strings.Replace(str, ",", "", -1)

	return str
}

func extractSizeAndRequest(arr []string) (uint64, string) {
	// the value at arr[9] is the nginx log entry's size in bytes.
	size, err := strconv.ParseUint(arr[9], 10, 64)
	if err != nil {
		fmt.Println("Error parsing size", err)
	}

	// Grab the request URL, split it by a forwardslash,
	// grab the first element (which is usually either a direct link or top level directory)
	req := strings.Split(arr[6], "/")[1]

	return size, req
}

func scanFile(file *os.File, distroMap map[string]uint64) map[string]uint64 {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		nginxEntryStr := scanner.Text()
		arr := strings.Split(nginxEntryStr, " ")

		// Discard all invalid requests (Those which don't begin with "GET")
		if arr[5] != "\"GET" {
			continue
		}

		size, req := extractSizeAndRequest(arr)

		// If distroMap[req] exists, add on the size, otherwise create the entry
		if _, ok := distroMap[req]; ok {
			distroMap[req] += size
		} else {
			distroMap[req] = size
		}
	}

	return distroMap
}

func main() {
	yesterdayString := getYesterday()

	fmt.Println(yesterdayString)

	distroMap := make(map[string]uint64)

	file, err := os.Open("./Yest.log")
	if err != nil {
		fmt.Println("Error opening log file", err)
	}

	distroMap = scanFile(file, distroMap)
	file.Close()

	repoList := []string{"alpine", "archlinux", "blender", "centos", "clonezilla", "cpan", "cran", "ctan", "cygwin", "debian", "debian-cd", "debian-security", "fedora", "fedora-epel", "freebsd", "gentoo", "gentoo-portage", "gnu", "gparted", "ipfire", "isabelle", "linux", "linuxmint", "manjaro", "odroid", "openbsd", "opensuse", "parrot", "raspbian", "sabayon", "serenity", "slackware", "slitaz", "tdf", "ubuntu", "ubuntu-cdimage", "ubuntu-ports", "ubuntu-releases", "videolan", "voidlinux"}

	db, err := sql.Open("sqlite3", "./mirrorband.sqlite")
	if err != nil {
		fmt.Println("Error opening DB", err)
		os.Exit(1)
	}

	for _, repo := range repoList {
		sqlStr := fmt.Sprintf("INSERT INTO distrousage (time, distro, bytes) VALUES (\"%s\", \"%s\", %d)", yesterdayString, repo, distroMap[repo])
		if _, err = db.Exec(sqlStr); err != nil {
			fmt.Println("Error executing insert query", err)
			os.Exit(1)
		}
	}
}
