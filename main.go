package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	// importing colly
	"github.com/gocolly/colly"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	url := os.Getenv("URL")
	targetCssClass := os.Getenv("CSS_CLASS")
	fmt.Println("class", targetCssClass)
	postDates := []string{""}
	// instantiate a new collector object
	c := colly.NewCollector(
	//colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)

	// adjust header so it won't timeout and works with cloudflare
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("Got: ", r.Request.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error: ", e.Error(), r.StatusCode)
	})

	// grab the date of the latest post
	c.OnHTML(targetCssClass, func(e *colly.HTMLElement) {
		date := e.ChildText(".entry-date")
		// Clean up the extracted data
		date = strings.TrimSpace(date)
		postDates = append(postDates, date)
	})

	// Visit the URL and start scraping
	pgErr := c.Visit(url)
	if pgErr != nil {
		log.Fatal(pgErr)
	}

	_, postDates = postDates[0], postDates[1:]

	// to get the current time
	currentTime := time.Now() //time.Parse("2006-01-02", "2025-01-22") //
	latestPostDate, tErr := time.Parse("2006-01-02", postDates[0])
	if tErr != nil {
		log.Fatal(tErr)
	}

	if currentTime.Equal(latestPostDate) {
		fmt.Println("new post updated")
	} else {
		fmt.Println("latest post time:", latestPostDate, ", current date:", currentTime)
	}
}
