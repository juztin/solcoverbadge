package main

import (
	"bytes"
	"flag"
	"fmt"
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
	coverage  = flag.String("coverage", "/src/coverage/index.html", "The path to the coverage index.html")
	badgePath = flag.String("badge", "/src/coverage/coverage.svg", "The filepath for the coverage badge")
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

func main() {
	flag.Parse()

	// Get stats
	log.Println("Reading coverage")
	percent, err := getStats(*coverage)
	if err != nil {
		fmt.Println(err)
	}

	// Generate coverage badge
	log.Println("Calling badge service")
	color := coverageColor(percent)
	value := coverageValue(percent)
	url := fmt.Sprintf(SHIELD_FMT, value, color)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Writing coverage badge
	log.Println("Writing badge")
	data, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(*badgePath, data, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
