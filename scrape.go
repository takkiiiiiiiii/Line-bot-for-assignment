package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)


const url = "https://elms.u-aizu.ac.jp/login/index.php"


func ScrapePage(payload_username, payload_password, payload_rememberusername string) ([]string, []string, string) {
	var schedule string
	var assignments []string
	var dates []string 
	err := exec.Command("curl", url, "-X", "GET", "-c", "cookie/login_cookie.txt", "-o", "html/login.html").Run()
	if err != nil {
		fmt.Printf("Failed to request : %v\n", err)
		os.Exit(1)
	}
	loginPage, err := os.Open("html/login.html")
	if err != nil {
		fmt.Printf("Failed to open %v\n", err)
		os.Exit(2)
	}
	defer loginPage.Close()

	auth_doc, err := goquery.NewDocumentFromReader(loginPage)
	if err != nil {
		fmt.Printf("Failed to read %v\n", err)
		os.Exit(3)
	}
	var val []string
	auth_doc.Find("input").Each(func(index int, item *goquery.Selection) {
		nameElement, _ := item.Attr("value")
		val = append(val, nameElement)
	})
	payload_loginToken := "logintoken=" + val[0]

	err = exec.Command("curl", "-X", "POST", url, "-s", "-L",
		"-F", "anchor=", "-F", payload_username, "-F", payload_password, "-F",
		payload_loginToken, "-F", payload_rememberusername,
		"-o", "html/mypage.html", "-b", "cookie/login_cookie.txt", "-c", "cookie/mypage_cookie.txt").Run()
	if err != nil {
		fmt.Printf("Failed to request : %v\n", err)
		os.Exit(4)
	}

	mypage, err := os.Open("html/mypage.html")
	if err != nil {
		fmt.Printf("Failed to open: %v\n", err)
		os.Exit(5)
	}
	defer mypage.Close()
	doc, err := goquery.NewDocumentFromReader(mypage)
	if err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		os.Exit(6)
	}
	selection := doc.Find("div.event").Find("div.overflow-auto")
	selection.Each(func(index int, item *goquery.Selection) {
		assignment := item.Find("a.text-truncate").Text()
		date := item.Find("div.date").Text()
		assignments = append(assignments, assignment)
		dates = append(dates, date)
	})
	return assignments, dates, schedule
}