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
	payload_username         = "username" + username
	payload_password         = "password=password" + password
	payload_rememberusername = "rememberusername=1" 
)

var data []string
var head []string
var assignment string
var i int
var s string

func main() {

	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
        os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
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
				case *linebot.TextMessage:
					replyMessage := message.Text
					if strings.Contains(replyMessage, postText) {
						err := exec.Command("curl", url, "-X", "GET", "-c", "cookie.txt", "-o", "login.html").Run()
						if err != nil {
							log.Fatalf("Failed to request : %v", err)
						}
						authPage, err := os.Open("login.html")
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

						fmt.Println(payload_loginToken)

						err = exec.Command("curl", "-X", "POST", url, "-s", "-L",
							"-F", "anchor=", "-F", payload_username, "-F", payload_password, "-F",
							payload_loginToken, "-F", payload_rememberusername,
							"-b", "cookie.txt", "-c", "cookie02.txt", "-o", "mypage.html").Run()
						if err != nil {
							log.Fatalf("Failed to request : %v", err)
						}

						myPage, err := os.Open("mypage.html")
						if err != nil {
							log.Fatalf("Failed to open %v", err)
						}
						defer myPage.Close()

						doc, err := goquery.NewDocumentFromReader(myPage)
						if err != nil {
							log.Fatalf("Failed to read %v", err)
						}
						content := doc.Find("div.card.rounded")
						content.Each(func(index int, item *goquery.Selection) {
							contents := item.Find("div.d-inline-block").Find("h3.name.d-inline-block").Text()                  // 課題の内容
							time := item.Find("div.description.card-body").Find("div.row").Find("div.col-11").Find("a").Text() //課題の教科名と締切日時
							ok := time + "\n" + "内容: " + contents
							data = append(data, ok)
						})

						fmt.Println(data) //　データがない場合
						if len(data) < 1 {
							notice := "直近の課題はありません。"
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(notice)).Do(); err != nil {
								log.Fatal(err)
							}
						} else {
							for i, s = range data {
								assignment += s + "\n" + "+-----------------------------+" + "\n"
							}

							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(assignment)).Do(); err != nil {
								log.Fatal(err)
							}
							//初期化
							data = append(data[i+1:], data[i+1:]...)
							assignment = ""
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
    if err := http.ListenAndServe(":" + os.Getenv("LINE_BOT_PORT"), nil); err != nil {
		log.Print(err)
	}
	time.Sleep(10 * time.Second)
}
