package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	url                      = "url"
	username                 = "studentID"
	password                 = "password"
	payload_username         = "username=" + username
	payload_password         = "password=" + password
	payload_rememberusername = "rememberusername=1"
)

var head []string
var schedule string
var assignments []string // assignment name
var courses []string     // course name
var links []string       // assignment link
var i int
var s string

func main() {

	bot, err := linebot.New(
		"fd1fd5ee8ea8d5608866d25bc8f4eff8",
		"G7wLav3sSFPlO+BZtbBvtlGDeFAB0iGm5mynU0jkPXZPLFwF1PMvWXoUBYOuiM25oO4/hsLEuJVzRfxwJ6U/ZfsKnywzM850aAz4ing3oYrHl8a0KYe+ViaEdT5mH0aKFedfKPBi+6oH5zDD8WKXYAdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	kadai := func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		postText := "課題"
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				//メッセージがテキスト形式の場合
				case *linebot.TextMessage: //文字列の場合
					replyMessage := message.Text
					if strings.Contains(replyMessage, postText) {
						err := exec.Command("curl", url, "-X", "GET", "-c", "cookie/login_cookie.txt", "-o", "html/login.html").Run()
						if err != nil {
							log.Fatalf("Failed to request : %v", err)
						}
						authPage, err := os.Open("html/login.html")
						if err != nil {
							log.Fatalf("Failed to open %v", err)
						}
						defer authPage.Close()

						auth_doc, err := goquery.NewDocumentFromReader(authPage)
						if err != nil {
							log.Fatalf("Failed to read %v", err)
						}
						var val []string
						auth_doc.Find("input").Each(func(index int, item *goquery.Selection) {
							nameElement, _ := item.Attr("value")
							val = append(val, nameElement)
						})

						payload_loginToken := "logintoken=" + val[1]

						err = exec.Command("curl", "-X", "POST", url, "-s", "-L",
							"-F", "anchor=", "-F", payload_username, "-F", payload_password, "-F",
							payload_loginToken, "-F", payload_rememberusername,
							"-o", "html/mypage.html", "-b", "cookie/login_cookie.txt", "-c", "cookie/mypage_cookie.txt").Run()
						if err != nil {
							log.Fatalf("Failed to request : %v", err)
						}

						err = exec.Command("curl", "-c", "cookie/calendar_cookie.txt", "-X", "GET", "https://elms.u-aizu.ac.jp/calendar/view.php?view=upcoming",
							"-b", "cookie/mypage_cookie.txt", "-o", "html/calendar.html").Run()
						if err != nil {
							log.Fatalf("Failed to request : %v", err)
						}

						calendar, err := os.Open("html/calendar.html")
						if err != nil {
							log.Fatalf("Failed to open %v", err)
						}
						defer calendar.Close()

						doc, err := goquery.NewDocumentFromReader(calendar)
						if err != nil {
							log.Fatalf("Failed to read %v", err)
						}
						homeworks := doc.Find("div.event")
						homeworks.Each(func(index int, item *goquery.Selection) {
							assignment := item.Find("div.card.rounded").Find("div.d-inline-block").Find("h3.name.d-inline-block").Text()
							assignments = append(assignments, assignment)

						})
						subjects := doc.Find("div.description.card-body")
						subjects.Each(func(index int, item *goquery.Selection) {
							course := item.Find("div.col-11").Find("a").Text()
							link, _ := item.Find("div.row.mt-1").Find("div.col-11").Find("a").Attr("href")
							courses = append(courses, course)
							links = append(links, link)
						})
						if len(assignments) < 1 {
							notice := "直近の課題はありません。"
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(notice)).Do(); err != nil {
								log.Fatal(err)
							}
						} else {
							for i = 0; i < len(assignments); i++ {
								schedule += courses[i] + "\n" + "内容: " + assignments[i] + "\n" + "リンク: " + links[i] + "\n" + "+-----------------------------+" + "\n"
							}

							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(schedule)).Do(); err != nil {
								log.Fatal(err)
							}
							//スライスの中身を全て削除
							assignments = nil
							courses = nil
							links = nil
							schedule = ""
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Fatal(err)
						}
					}
				case *linebot.StickerMessage: //スタンプの場合
					replyMessage := fmt.Sprintf(
						"sticker username is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}

	http.HandleFunc("/kadai", kadai)
	if err := http.ListenAndServe(":port", nil); err != nil {
		log.Print(err)
	}
	time.Sleep(10 * time.Second)
}
