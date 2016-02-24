package main

import (
  "fmt"
  "log"

  "github.com/PuerkitoBio/goquery"
)

/*
@param Url of the repo
About: Sends retrieval functions for fetching files
*/

func RepoRet(uri string){
  doc, err := goquery.NewDocument(uri)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("tr.js-navigation-item").Each(func(i int, s *goquery.Selection) {
    s.Find("td.content").Each(func(j int, sel *goquery.Selection){
      fmt.Println(sel.Find("a").Text())
    })
  })
}

func ExampleScrape() {
  baseuri := "https://github.com"
  var user string
  fmt.Println("Enter username")
  fmt.Scanf("%s",&user)
  doc, err := goquery.NewDocument("https://github.com/"+user+"?tab=repositories") 
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("h3.repo-list-name").Each(func(i int, s *goquery.Selection) {
    val, _ := s.Find("a").Attr("href")
    RepoRet(baseuri+val)
  })
  
}

func main() {
  ExampleScrape()
  //RepoRet("https://github.com/shubhodeep9/Gitlines")
}