package main

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

const
(
	Welcome       = "سلام ، خوش آمدید"
	GetVPN        = "دریافت کد برای رفع تحریم های آمریکا"
	EnterCode     = "لطفا کد فعال سازی را وارد نمایید:"
	CodeIsInvalid = "کد وارد شده قبلا استفاده شده یا صحیح نمی باشد. لطفا کد معتبر وارد نمایید"
	YouAreForbidden = "Sorry You haven`t permission to do this action!"
	EnterForbiddenID = "Enter Forbidden User ID"
	EnterIDToUnblock="Enter User ID to Unblock"
	)

var
(
	DBUser     string
	DBName     string
	DBPassword string
	DBHost     string
	DBEngin    string
	Token      string
	DB         *sqlx.DB
	bot        *tgbotapi.BotAPI
	LogFile *os.File
	VPNT1 string
	VPNT2 string
)

type User struct {
	ID int `json:"id"`
	Code   string `json:"code"`
	TelegramID int `json:"telegramid"`
	ChatID int64 `json:"ChatId"`
	MobileNumber string `json:"mobile_number"`
	Ip string `json:"Ip"`
	IsSuccess string `json:"IsSuccess"`
	Password string `json:"Password"`
	Username string `json:"Username"`
}

type VpnResult struct {
ErrorMessage string `json:"ErrorMessage"`
IsSuccess string `json:"IsSuccess"`
}
type VPNInfo struct {
	Ip string `json:"ip"`
	Username string `json:"username"`
	Pwd string `json:"pwd"`
	Code string `json:"code"`
}
var VPNCodeTaleSchema=`
CREATE TABLE IF NOT EXISTS public.vpncodes
(
  id integer NOT NULL ,
  code text NOT NULL,
  ctime date DEFAULT now(),
  ip text,
  username text,
  pwd text,
  telegramid text NOT NULL,
	mobilenumber text,
  CONSTRAINT vpncodes_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.vpncodes
  OWNER TO postgres;
`
var ForbiddenTableSchema =`CREATE TABLE IF NOT EXISTS public.forbidden
(
   id integer,
   tid text,
   PRIMARY KEY (id)
)
WITH (
 OIDS = FALSE
);
`