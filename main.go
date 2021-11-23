package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	siteURL := flag.String("site", "", "site URL")
	imgSelector := flag.String("img-selector", "", "img selector")
	flag.Parse()

	// Request the HTML page.
	res, err := http.Get(*siteURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(*imgSelector).Each(
		func(i int, s *goquery.Selection) {
			// For each item found, get the src
			src := s.AttrOr("src", "")
			fmt.Printf("Image %d: %s\n", i, src)

			// Get the data
			resp, err := http.Get(src)
			if err != nil {
				log.Print(err)
				return
			}
			defer resp.Body.Close()

			// Create the file
			filepath := strings.Replace(src, "/", "-", -1)
			filepath = strings.Replace(filepath, ":", "-", -1)
			file, err := os.Create(filepath)
			if err != nil {
				log.Print(err)
				return
			}
			defer file.Close()

			// Write the body to file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Print(err)
				return
			}

		},
	)

}
