package main

import (
  "fmt"
  "log"

  "github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
  doc, err := goquery.NewDocument("https://github.com/shubhodeep9?tab=repositories") 
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("h3.repo-list-name").Each(func(i int, s *goquery.Selection) {
    band := s.Find("a")
    fmt.Println(band.Text())
  })
  //fmt.Println(doc)
}

func main() {
  ExampleScrape()
}