package main

import(
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)


func ReplyKadai(w http.ResponseWriter, req *http.Request) {
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
					assignments, dates, schedule := ScrapePage(payload_username, payload_password, payload_rememberusername)
					
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