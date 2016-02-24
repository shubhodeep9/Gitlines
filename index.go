package main

import (
  "fmt"
  "log"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "strconv"
)

/*
@TODO Make the login script
@TODO Get private repository data
*/

var count int64 = 0



func GetLines(uri string) int64 {
  doc, err := goquery.NewDocument(uri)
  if err != nil {
    log.Fatal(err)
  }
  val := doc.Find("div.file-info").Text()
  ind := strings.Index(strings.TrimSpace(val)," lines")
  //fmt.Println(strings.TrimSpace(val)[:ind])
  if ind == -1 {
    return int64(0)
  }
  str, _ := strconv.Atoi(strings.TrimSpace(val)[:ind])
  return int64(str)
}

/*
@param Url of the repo, Goroutine channel, BaseUrl (default=https://github.com)
return void
About: Sends retrieval functions for fetching files
*/
func RepoRet(uri string, c chan int,baseuri string){
  in_c := make(chan int)
  doc, err := goquery.NewDocument(uri)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("tr.js-navigation-item").Each(func(i int, s *goquery.Selection) {
    s.Find("td.content").Each(func(j int, sel *goquery.Selection){
      val, _ := sel.Prev().Find("svg").Attr("class")
      //var filetype string
      if val == "octicon octicon-file-text" {
        val, _ = sel.Find("a").Attr("href")
        //fmt.Println(val)
        count=count+GetLines(baseuri+val)
      } else {
        val, _ = sel.Find("a").Attr("href")
        go RepoRet(baseuri+val,in_c,baseuri)
        <-in_c
      }
      //fmt.Println(sel.Find("a").Text()+" "+filetype)
    })
  })
  c <- 1

}

func ExampleScrape() {
  c := make(chan int)
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
    go RepoRet(baseuri+val,c,baseuri)
    <-c
  })
  
}

func main() {
  ExampleScrape()
  // out_c := make(chan int)
  // //ExampleScrape()
  // go RepoRet("https://github.com/shubhodeep9/Gitlines",out_c,"https://github.com")
  // <-out_c
  // fmt.Println(count)
  //fmt.Println(GetLines("https://github.com/shubhodeep9/Gitlines/blob/master/.gitignore"))
}