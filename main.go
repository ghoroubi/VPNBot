package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"os"
	"strconv"
)

//import "fmt"

func main() {
	InitConfig()
	DBConnect()
	var err error
	// CodeGenerator()
	 //x:=GetCodes()

	bot,err=tgbotapi.NewBotAPI(Token)

	if err != nil {
log.Println(err.Error())
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err.Error())
	}
	for update := range updates {
		if update.Message != nil {
			log.Println("TID:",update.Message.From.ID,"User Input: "+update.Message.Text)
			fmt.Println("User Input: "+update.Message.Text)
			switch update.Message.Text {
			case "/start":
				SendTextMessage(update.Message.Chat.ID,Welcome,GetHomeKeys)
			case "/closebyadmin":
				os.Exit(0)
			case GetVPN:
				if IsForbiddens(update.Message.From.ID) {
					fmt.Println("Frbidden User.",update.Message.From.ID)
					log.Println("Frbidden User.",update.Message.From.ID)
					msg:=tgbotapi.NewMessage(update.Message.Chat.ID,YouAreForbidden)
					msg.ReplyMarkup=GetHomeKeys()
					bot.Send(msg)

				}else {
					SendForceReply(update.Message.Chat.ID, EnterCode)
				}
			case "/b":
				SendForceReply(update.Message.Chat.ID,EnterForbiddenID)
			case "/u":
				SendForceReply(update.Message.Chat.ID,EnterIDToUnblock)
			case "/list":
				BlockedUsersList(update.Message.Chat.ID)
			}//end of switch
		}
		if update.Message.ReplyToMessage != nil {
			switch update.Message.ReplyToMessage.Text {

			case EnterCode:
				CodeReview(update.Message.Chat.ID,update.Message.Text)
			case CodeIsInvalid:
				CodeReview(update.Message.Chat.ID,update.Message.Text)
			case EnterForbiddenID:
				btid,_:=strconv.Atoi(update.Message.Text)
				if BlockUser(btid){
					msg:=tgbotapi.NewMessage(update.Message.Chat.ID,"OK.👍")
					msg.ReplyMarkup=GetHomeKeys()
					bot.Send(msg)
				}else {msg:=tgbotapi.NewMessage(update.Message.Chat.ID,"Failed.Try Again. 😔")
				msg.ReplyMarkup=GetHomeKeys()
				bot.Send(msg)
				}
			case EnterIDToUnblock:
				ubtid,_:=strconv.Atoi(update.Message.Text)
				if Unblock(ubtid){
					msg:=tgbotapi.NewMessage(update.Message.Chat.ID,"OK.👍")
					msg.ReplyMarkup=GetHomeKeys()
					bot.Send(msg)
				}else {msg:=tgbotapi.NewMessage(update.Message.Chat.ID,"Failed.Try Again. 😔")
					msg.ReplyMarkup=GetHomeKeys()
					bot.Send(msg)
				}
		}//end of switch
		}//end of else if
	}//end of for

//	fmt.Println(GetCodes())
	defer DB.Close()
	defer LogFile.Close()
}
