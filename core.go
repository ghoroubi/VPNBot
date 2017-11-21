package main

import (
	"github.com/spf13/viper"
	"log"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gopkg.in/telegram-bot-api.v4"
_	"net/http"
_	"io/ioutil"
_	"encoding/json"
	"strconv"
	//"unicode"
	//"golang.org/x/text/transform"
//	"golang.org/x/text/unicode/norm"
	//norm2 "projects/GO.HERMES/go/src/golang.org/x/text/unicode/norm"
	"strings"
	"os"
	"time"
	"database/sql"
	"io/ioutil"
	"encoding/json"
	"net/http"
)
type fn func() tgbotapi.ReplyKeyboardMarkup
func InitConfig(){
	conf:=viper.New()
	conf.SetConfigName("config")
	conf.AddConfigPath(".")
	conf.SetConfigType("yml")

	err:=conf.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
		defer log.Panic(err.Error())
		fmt.Println("Coniguration Error")
	}
	defer fmt.Println("Configuration has been set succesfuly")
	log.Println("Configuration has been set succesfuly")

	//get Robot Token:
	Token=conf.GetString("Public.Tok")
	//Database:
	DBHost=conf.GetString("DB.Host")
	DBName=conf.GetString("DB.DBName")
	DBEngin=conf.GetString("DB.Engine")
	DBUser=conf.GetString("DB.User")
	DBPassword=conf.GetString("DB.Passwor")
	//////////////////
	VPNT1=conf.GetString("Public.VPNTitle1")
	VPNT2=conf.GetString("Public.VPNTitle2")
	LoggerInit()
	/*var fileerror error
	LogFile, err = os.Create("log.txt")
	LogFile, err = os.Create("log.txt", os.O_WRONLY, 0666)
	log.SetOutput(LogFile)
	if fileerror != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		}*/

}
func LoggerInit(){
	var err error
	LogFile, err = os.OpenFile("logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
	log.New(LogFile,time.Now().String(),0)
	log.SetOutput(LogFile)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer LogFile.Close()
}
func DBConnect(){
	var err error
	query:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",DBHost,DBUser,DBPassword,DBName)
	DB,err=sqlx.Connect(DBEngin,query)
	if err!=nil {
		log.Println(err.Error())
	}
	fmt.Println("DB init...")
	log.Println("DB init...")
	DB.MustExec(VPNCodeTaleSchema)
	DB.MustExec(ForbiddenTableSchema)
}
func CodeReview(chatId int64,code string ) {
	var user User
	fmt.Println("Code: ", code)
	log.Println("TID:", chatId, "Code: ", code)
	_code := Normalize(code)

	getRequest := fmt.Sprintf("http://37a602f845f9.sn.mynetname.net:1367/methods.svc/code/%s", _code)
	res, err := http.Get(getRequest)

	if err != nil {
		log.Println("TID:", chatId, "Error Occured!")
		fmt.Println("Error Occured!")
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(res.Body)
	json.Unmarshal(body, &user)
	switch user.IsSuccess {
	case "true":
		tid := strconv.Itoa(int(chatId))
		m := fmt.Sprintf("\n IP (Server):  %s \n Username:  %s \n Password:  %s \n Secret Key:123456 ", user.Ip, user.Username, user.Password)
		fmt.Println("Valid Code Enterd:", code)
		fmt.Println("Inserted To DB: <<", user.Ip, user.Username, user.Password, tid, " >>")
		log.Println("TID:", chatId, "Valid Code Enterd:", code)
		log.Println("TID:", chatId, "Inserted To DB: <<", user.Ip, user.Username, user.Password, tid, " >>")
		insQuery := fmt.Sprintf("INSERT INTO vpncodes(code,ip,username,pwd,telegramid,mobilenumber) VALUES(%s,%s,%s,%s,%s,%s)",
			code, user.Ip, user.Username, user.Password, tid, user.MobileNumber)
		_, err := DB.Query(insQuery)
		if err != nil {
			fmt.Println(err.Error())
			//log.Panic(err.Error())
		}
		strReq:=fmt.Sprintf("http://37a602f845f9.sn.mynetname.net:1367/methods.svc/set_mobile/%s/%s",code,user.MobileNumber)
		//strReq:=fmt.Sprintf("http://192.168.2.168:2215/methods.svc/set_mobile/%s/%s",code,tid)
		_,er:=http.Get(strReq)
		if er!=nil {
			log.Println(err.Error())
		}
		fmt.Printf("\n User ChatID is :%d \n", chatId)
		log.Printf("\n User ChatID:%d \n", chatId)
		mMessage := VPNT1 + "\n" + VPNT2 + "\n" + m
		msg := tgbotapi.NewMessage(chatId, mMessage)
		msg.ReplyMarkup = GetHomeKeys()
		bot.Send(msg)
	case "false":
		fmt.Println("Invalid Code Enterd...")
		log.Println("Invalid Code Enterd...")
		SendForceReply(chatId, CodeIsInvalid)
	}
}
func GetHomeKeys() tgbotapi.ReplyKeyboardMarkup {
	rep:=tgbotapi.ReplyKeyboardMarkup{}
	commands:=[][]tgbotapi.KeyboardButton{}
	getvpnButton:=tgbotapi.KeyboardButton{Text:GetVPN}
	row0:=[]tgbotapi.KeyboardButton{getvpnButton}
	commands=append(commands,row0)

	rep.Keyboard=commands
	rep.ResizeKeyboard=true

	return rep
}

func SendTextMessage(chatId int64, text string, keys fn) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = keys()
	bot.Send(msg)
}

func SendForceReply(chatId int64, text string) {
	fmt.Println("ForceReply: ", text)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
	bot.Send(msg)
}
func Normalize(strNum string) string{
	var x string
	for i := 0; i < len(strNum);i++  {
	temp:=charCodeAt(strNum,int(i))
		if temp >=1776 && temp<= 1785 {
			x+=string(temp-1728)

		}else{
			x+=string(temp)

		}

	}
	 u:=strings.Replace(x, "\x00", "", -1)
log.Println("User Verification Code Normalized Succesfully.Code:",u)
	return u

}
func charCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}
	return 0
}
func IsForbiddens(telegid int) bool{
	//var id int
	var telId string
	status:=false
	fmt.Println(telId)
	log.Println(telegid," Is online User`s TelegramID.")
	err:=DB.QueryRow("SELECT tid FROM forbidden where tid=$1",strconv.Itoa(telegid)).Scan(&telId) // Note: Ignoring errors for brevity
		switch {
		case err == sql.ErrNoRows:
			fmt.Printf("user %d Is Permitted",telegid)
			log.Printf("user %d Is Permitted",telegid)
		case err != nil:
		log.Fatal(err)
		default:
			fmt.Printf("telegram id  %s is forbidden User. \n", telId)
			log.Printf("telegram id  %s is forbidden User. \n", telId)
		status=true
		}
		return status
	}
func BlockUser(tid int) bool{
	result:=true
	_,err:=DB.Exec("INSERT INTO forbidden(tid) VALUES($1)",strconv.Itoa(tid))
	if err!=nil {
		log.Println("Error in inserting data.Error: ",err.Error())
		fmt.Println("Error in inserting data.Error: ",err.Error())
		result=false
	}
	log.Println("Succesfully inserted data.")
	fmt.Println("Succesfully inserted data.")
	return result
}
func Unblock(tid int)bool{
	result:=true
	_,err:=DB.Exec("DELETE FROM forbidden WHERE tid=$1",strconv.Itoa(tid))
	if err!=nil {
		log.Println("Error in Executing.Error: ",err.Error())
		fmt.Println("Error in Executing.Error: ",err.Error())
		result=false
	}
	log.Println("Succesfully Unblocked.")
	fmt.Println("Succesfully Unblocked.")
	return result
}
func BlockedUsersList(chatId int64){
	var list string
	rows,err:=DB.Query("Select tid from forbidden")
	if err!=nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		}
	for rows.Next(){
		var btId string
		if err:=rows.Scan(&btId);err!=nil{
			log.Println(err.Error())
			fmt.Println(err.Error())
			log.Fatal(err)
		}
		list=list+"\n"+btId
		if err:=rows.Err();err!=nil{
			log.Println(err.Error())
			fmt.Println(err.Error())
			log.Fatal(err)
		}
	}

	msg:=tgbotapi.NewMessage(chatId,list)
	msg.ReplyMarkup=GetHomeKeys()
	bot.Send(msg)
}





