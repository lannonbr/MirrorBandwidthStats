package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

func extractSizeAndRequest(arr []string) (uint64, string, bool) {
	// the value at arr[9] is the nginx log entry's size in bytes.
	size, err := strconv.ParseUint(arr[9], 10, 64)
	if err != nil {
		fmt.Println("Error parsing size", err)
	}

	reqArr := strings.Split(arr[6], "/")

	if len(reqArr) < 2 {
		return 0, "", false
	}

	req := reqArr[1]

	return size, req, true
}

func scanFile(filename string, distroMap map[string]uint64, date string) map[string]uint64 {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading all", err)
	}

	contentStr := string(content[:])
	contentStrArr := strings.Split(contentStr, "\n")

	lines := []string{}

	for _, line := range contentStrArr {
		if strings.Contains(line, date) {
			lines = append(lines, line)
		}
	}

	for _, entry := range lines {
		arr := strings.Split(entry, " ")

		// Discard all invalid requests (Those which don't begin with "GET")
		if arr[5] != "\"GET" {
			continue
		}

		// Discard any unusual HTTP logs
		if !strings.Contains(arr[7], "HTTP") {
			continue
		}

		size, req, valid := extractSizeAndRequest(arr)

		if !valid {
			continue
		}

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

	ysSplit := strings.Split(yesterdayString, "/")
	dat := fmt.Sprintf("%s/%s/%s", ysSplit[1], ysSplit[0], ysSplit[2])

	distroMap := make(map[string]uint64)

	distroMap = scanFile("./Yest.log", distroMap, dat)

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
