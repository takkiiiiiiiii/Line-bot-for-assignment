package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const url  = "https://elms.u-aizu.ac.jp/login/index.php"

var assignments []string
var dates []string
var schedule string


func main() {
	username  := os.Getenv("LMS_ID")
	password  := os.Getenv("LMS_PASS")
	payload_username  := "username=" + username
	payload_password  := "password=" + password
	payload_rememberusername := "rememberusername=1"
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_ACCESS_SECRET_ASSIGNMENT"),
        os.Getenv("LINE_CHANNEL_ACCESS_TOKEN_ASSIGNMENT"),
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
							fmt.Printf("Failed to request : %v", err)
							os.Exit(1)
						}
						loginPage, err := os.Open("html/login.html")
						if err != nil {
							fmt.Printf("Failed to open %v", err)
							os.Exit(2)
						}
						defer loginPage.Close()

						auth_doc, err := goquery.NewDocumentFromReader(loginPage)
						if err != nil {
							fmt.Printf("Failed to read %v", err)
							os.Exit(3)
						}
						var val []string
						auth_doc.Find("input").Each(func(index int, item *goquery.Selection) {
							nameElement, _ := item.Attr("value")
							val = append(val, nameElement)
						})
						payload_loginToken := "logintoken=" + val[0]
						fmt.Println(val[0])

						err = exec.Command("curl", "-X", "POST", url, "-s", "-L",
							"-F", "anchor=", "-F", payload_username, "-F", payload_password, "-F",
							payload_loginToken, "-F", payload_rememberusername,
							"-o", "html/mypage.html", "-b", "cookie/login_cookie.txt", "-c", "cookie/mypage_cookie.txt").Run()
						if err != nil {
							fmt.Printf("Failed to request : %v", err)
							os.Exit(4)
						}

						mypage, err := os.Open("html/mypage.html")
						if err != nil {
							fmt.Printf("Failed to open: %v", err)
							os.Exit(5)
						}
						defer mypage.Close()
						doc, err := goquery.NewDocumentFromReader(mypage)
						if err != nil {
							fmt.Printf("Failed to read: %v", err)
							os.Exit(6)
						}
						// ここから修正
				
						selection := doc.Find("div.event").Find("div.overflow-auto")
						selection.Each(func(index int, item *goquery.Selection) {
							assignment := item.Find("a.text-truncate").Text()
							date := item.Find("div.date").Text()
							assignments = append(assignments, assignment)
							dates = append(dates, date)
						})
						
						if len(assignments) == 0 {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("サーバメンテナンスのため、現在、LMSは利用できません。")).Do(); err != nil {
								log.Fatal(err)
							}
						} else if len(assignments) / 2 < 3 {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(assignments[1])).Do(); err != nil {
								log.Fatal(err)
							}
						} else {
							for i := range assignments {
								schedule += "+--------------------------+\n" + strconv.Itoa(i+1) + " " + assignments[i] + "\n" + dates[i] + "\n"
							}
							schedule += "+--------------------------+\n"
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(schedule)).Do(); err != nil {
								log.Fatal(err)
							}
						}
						// }
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
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Print(err)
	}
	time.Sleep(10 * time.Second)
}

