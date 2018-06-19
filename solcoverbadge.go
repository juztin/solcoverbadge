package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const SHIELD_FMT = "https://img.shields.io/badge/coverage-%s-%s.svg?longCache=true&style=flat-square"

var (
	coverage    = flag.String("coverage", "/src/coverage/index.html", "The path to the coverage index.html")
	badgePath   = flag.String("badge", "/src/coverage/coverage.svg", "The filepath for the coverage badge")
	percentPath = flag.String("percent", "/src/coverage/coverage.txt", "The filepath for the coverage percent")
)

func getStats(file string) (float64, error) {
	// Read file
	b, err := ioutil.ReadFile(*coverage)
	if err != nil {
		return -1, err
	}

	// Parse HTML
	reader := bytes.NewReader(b)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return -1, err
	}

	// Query for percentage
	value := strings.TrimSpace(doc.Find(".wrapper .pad1y > span.strong").Last().Text())
	percent := strings.Split(value, "%")[0]
	return strconv.ParseFloat(percent, 64)
}

func coverageColor(percentage float64) string {
	switch {
	case percentage > 80:
		return "green"
	case percentage > 70:
		return "gold"
	case percentage > 50:
		return "orange"
	}
	return "red"
}

func coverageValue(percent float64) string {
	if percent == -1 {
		return "failing"
	}
	return fmt.Sprintf("%.2f%%25", percent)
}

func getShield(percent float64) (*http.Response, error) {
	color := coverageColor(percent)
	value := coverageValue(percent)
	url := fmt.Sprintf(SHIELD_FMT, value, color)
	return http.Get(url)
}

func writePercentFile(percent float64, filePath string) error {
	data := []byte(fmt.Sprintf("%.2f", percent))
	return ioutil.WriteFile(filePath, data, 0644)
}

func writeBadge(data io.ReadCloser, filePath string) error {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, b, 0644)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Get stats
	log.Println("Reading coverage")
	percent, err := getStats(*coverage)
	checkErr(err)

	// Write percent
	log.Println("Writing percent")
	err = writePercentFile(percent, *percentPath)
	checkErr(err)

	// Get badge
	log.Println("Calling badge service")
	r, err := getShield(percent)
	checkErr(err)
	defer r.Body.Close()

	// Write badge
	log.Println("Writing badge")
	err = writeBadge(r.Body, *badgePath)
	checkErr(err)
}
