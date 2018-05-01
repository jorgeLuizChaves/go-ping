package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const exitByUser = 0
const numberOfChecks = 3
const logFile = "app.log"
const sitesConfigFile = "sites.txt"
const timeIntervalCheck = 1 * time.Second
const dateFormatPattern = "02/01/06 - 15:04:06"

func main() {

	for {
		showAppName()
		showAppOptions()
		userOption := userOption()
		app(userOption)
	}
}

func showAppName() {
	fmt.Println("-------------------------------")
	fmt.Println("-  SUS   Saúde Única do Site  -")
	fmt.Println("-------------------------------")
}

func showAppOptions() {
	fmt.Println("Select the options")
	fmt.Println("[1] - monitor sites")
	fmt.Println("[2] - see logs")
	fmt.Println("[0] - exit")
}

func userOption() int {
	var option int
	fmt.Print("Your option: ")

	fmt.Scan(&option)
	return option
}

func app(option int) {
	switch option {
	case 1:
		monitorWebSites()
	case 2:
		readLogFile()
	case 0:
		os.Exit(exitByUser)
	default:
		fmt.Println("Invalid Option")
	}
}

func monitorWebSites() {
	webSites := readSitesMonitoring()
	fileLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, site := range webSites {
		for i := 0; i < numberOfChecks; i++ {
			resp, err := http.Get(site)
			if err != nil {
				log.Fatal(err)
			} else {
				logLine := time.Now().Format(dateFormatPattern) + " - " + "site: " + site + " statusCode: " + strconv.FormatInt(int64(resp.StatusCode), 10) + "\n"
				fmt.Println(logLine)
				fileLog.WriteString(logLine)
			}
			time.Sleep(timeIntervalCheck)
		}
	}
	fileLog.Close()
}

func readLogFile() {
	file, err := ioutil.ReadFile(logFile)

	if err != nil {
		log.Fatal("an error occured to open the file")
		return
	}
	fmt.Printf("%s", file)
	fmt.Println()
}

func readSitesMonitoring() []string {
	var sites = []string{}

	file, err := os.Open(sitesConfigFile)
	if err != nil {
		log.Fatal(err)
		return sites
	}
	reader := bufio.NewReader(file)

	for {
		site, _, err := reader.ReadLine()
		if err == io.EOF {
			file.Close()
			break
		}
		sites = append(sites, string(site))
	}
	file.Close()
	return sites
}
