package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
	"github.com/howeyc/gopass"
	"log"
	"strconv"
	"strings"
)

/*
@TODO Output Total lines plus total repos, and many more statistics
*/

var count int64 = 0

var bow *browser.Browser = surf.NewBrowser()

//GetLines takes input as a url and the line channel
//gives the output in line channel
//Calculates lines in each file
func GetLines(uri string, line chan int64) {
	bow.Open(uri)

	val := bow.Find("div.file-info").Text()
	ind := strings.Index(strings.TrimSpace(val), " lines")
	//fmt.Println(strings.TrimSpace(val)[:ind])
	if ind == -1 {
		line <- int64(0)
	} else {
		str, _ := strconv.Atoi(strings.TrimSpace(val)[:ind])
		line <- int64(str)
	}
}

/*
@param Url of the repo, Goroutine channel, BaseUrl (default=https://github.com)
return void
About: Sends retrieval functions for fetching files
*/
func RepoRet(uri string, c chan int, baseuri string) {
	in_c := make(chan int)
	line := make(chan int64)
	bow.Open(uri)

	bow.Find("tr.js-navigation-item").Each(func(i int, s *goquery.Selection) {
		s.Find("td.content").Each(func(j int, sel *goquery.Selection) {
			val, _ := sel.Prev().Find("svg").Attr("class")
			//var filetype string
			if val == "octicon octicon-file-text" {
				val, _ = sel.Find("a").Attr("href")
				//fmt.Println(val)
				go GetLines(baseuri+val, line)
				count = count + <-line
			} else {
				if sel.Find("a").Text() != "Godeps" {
					val, _ = sel.Find("a").Attr("href")
					go RepoRet(baseuri+val, in_c, baseuri)
					<-in_c
				}
			}
			fmt.Println(sel.Find("a").Text())
		})
	})
	c <- 1

}

func ExampleScrape() {
	c := make(chan int)
	baseuri := "https://github.com"
	/*
	  Login Script
	  @param Username Password
	  @Url https://github.com/login
	*/
	var (
		user     string
		password []byte
	)
	fmt.Println(`Login to proceed:
    If Login fails, total public repositories lines will be shown,
    otherwise Public plus Private of the user.
    Thank You, May the Force be with you!`)
	fmt.Println("Enter username")
	fmt.Scanf("%s", &user)
	fmt.Println("Enter password")
	password, _ = gopass.GetPasswdMasked()
	bow.Open("https://github.com/login")
	fm, err := bow.Form("form")
	if err != nil {
		fmt.Println(err)
	}

	fm.Input("login", user)
	fm.Input("password", string(password))
	err = fm.Submit()
	if err != nil {
		log.Fatal(err)
	}
	//Login End

	bow.Open("https://github.com/" + user + "?tab=repositories")

	bow.Find("h3.repo-list-name").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Find("a").Attr("href")
		fmt.Println(val)
		go RepoRet(baseuri+val, c, baseuri)
		<-c
		fmt.Println(count)
	})

}

func main() {
	// out_c := make(chan int)
	// go RepoRet("https://github.com/shubhodeep9/AnalyticsWeekly",out_c,"https://github.com")
	// <-out_c
	ExampleScrape()
	fmt.Println(count)
}
