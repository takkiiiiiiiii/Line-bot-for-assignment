package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sclevine/agouti"
)

const (
	url  = "url"
	id   = "my_id"
	pass = "my_password"
)

var data []string
var head []string
var assignment string
var i int
var s string

func main() {

	bot, err := linebot.New(
		"発行されたチャネルシークレット",
		"発行されたチャネルアクセストークン",
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

						driver := agouti.ChromeDriver(
							agouti.ChromeOptions("args", []string{
								"--headless",
							}),
							agouti.Debug,
						)
						if err := driver.Start(); err != nil {
							log.Fatalf("Failed to start driver:%v", err)
						}
						defer driver.Stop()

						page, err := driver.NewPage()

						if err != nil {
							log.Fatalf("Failed to open page:%v", err)
						}
						// ログインページに遷移
						if err := page.Navigate(url); err != nil {
							log.Fatalf("Failed to navigate:%v", err)
						}

						// ID, Passの要素を取得し、値を設定
						identity := page.FindByID("username")
						password := page.FindByID("password")
						identity.Fill(id)
						password.Fill(pass)
						input := page.FindByLabel("Remember username")
						if err = input.Click(); err != nil {
							log.Fatalf("Failed to navigate:%v", err)
						}
						btn := page.FindByButton("Log in")
						if err = btn.Click(); err != nil {
							log.Fatalf("Failed to navigate:%v", err)
						}

						link := page.FindByLink("カレンダーへ移動する ...")

						if err := link.Click(); err != nil {
							log.Fatalf("Failed to navigate:%v", err)
						}

						html, _ := page.HTML() //遷移先のhtmlを取得
						file, err := os.Create("file.html")
						if err != nil {
							log.Fatalf("Failed to create file %v", err)
						}
						defer file.Close()
						file_open, err := os.Open("/Users/takkiiiiiiiii/Go/test/test7/file.html")
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
							ok := time + "\n" + contents
							data = append(data, ok)
						})

						for i, s = range data {
							assignment += s + "\n"
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
	if err := http.ListenAndServe(":1111", nil); err != nil {
		log.Print(err)
	}
	time.Sleep(10 * time.Second)
}
