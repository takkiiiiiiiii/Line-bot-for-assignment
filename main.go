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
	url  = "https://elms.u-aizu.ac.jp/login/index.php"
	id   = "s1290077"
	pass = "takkiiiiiiiii25"
)

var data []string
var head []string
var assignment string
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
						subCmd := "https://elms.u-aizu.ac.jp/login/index.php"
						err := exec.Command("curl", subCmd, "-X", "GET", "-c", "cookie.txt", "-o", "login.html").Run()
						if err != nil {
							log.Fatalf("Failed to execute : %v", err)
						}
						auth_page, err := os.Open("/Users/yudai/Go/test/line-bot-for-assignment/login.html")
						if err != nil {
							log.Fatalf("Failed to open %v", err)
						}
						defer auth_page.Close()

						auth_doc, err := goquery.NewDocumentFromReader(auth_page)
						if err != nil {
							log.Fatalf("Failed to read %v", err)
						}
						var val []string
						auth_doc.Find("input").Each(func(index int, item *goquery.Selection) {
							nameElement, _ := item.Attr("value")
							val = append(val, nameElement)
						})
						//val[1]にloginToken

						file, err := os.Create("file.html")
						if err != nil {
							log.Fatalf("Failed to create file %v", err)
						}
						defer file.Close()
						file_open, err := os.Open("/Users/yudai/Go/test/test7/file.html")
						if err != nil {
							log.Fatalf("Failed to open %v", err)
						}
						defer file_open.Close()
						file.WriteString(html)
						//htmlをパース
						doc, err := goquery.NewDocumentFromReader(file_open)
						if err != nil {
							log.Fatalf("Failed to read %v", err)
						}
						content := doc.Find("div.card.rounded")
						content.Each(func(index int, item *goquery.Selection) {
							contents := item.Find("div.d-inline-block").Find("h3.name.d-inline-block").Text()
							time := item.Find("div.description.card-body").Find("div.row").Find("div.col-11").Find("a").Text() //課題の教科名と締切日時
							ok := time + "\n" + "内容: " + contents
							data = append(data, ok)
						})

						for i, s = range data {
							assignment += s + "\n" + "+-----------------------------+" + "\n"
						}

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(assignment)).Do(); err != nil {
							log.Fatal(err)
						}
						//初期化
						data = append(data[i+1:], data[i+1:]...)
						assignment = ""
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Fatal(err)
						}
					}
				case *linebot.StickerMessage: //スタンプの場合
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}

	http.HandleFunc("/kadai", kadai)
	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
	time.Sleep(10 * time.Second)
}
