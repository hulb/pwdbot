package structs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pwdbot/utils"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var UniqBot *tb.Bot

func RegisterBot(bot *tb.Bot) {
	UniqBot = bot
}

type UserData struct {
	User     *tb.User           `json:"user"`
	Accounts map[string]Account `json:"accounts"`
}

func (userData UserData) GetFilePath() string {
	return utils.GetCurrentDirectory() + "/data/" + fmt.Sprintf("%d", userData.User.ID)
}

func (userData UserData) Save() {
	if userDataBytes, err := json.Marshal(&userData); err == nil {
		utils.WriteFile(userData.GetFilePath(), userDataBytes)
	} else {
		log.Println("save user data fail")
		panic(err)
	}
}

func GetUserData(user *tb.User) UserData {
	userDataFilePath := utils.GetCurrentDirectory() + "/data/"
	userDataFileName := fmt.Sprintf("%d", user.ID)
	var userData UserData
	if !utils.Exists(userDataFilePath + userDataFileName) {
		userData.User = user
		userData.Accounts = make(map[string]Account)

		err := os.Mkdir(userDataFilePath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}

		f, err := os.Create(userDataFilePath + userDataFileName)
		if err != nil {
			log.Println(err.Error())
		} else {
			if userDataBytes, err := json.Marshal(&userData); err == nil {
				_, err = f.Write(userDataBytes)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}

		f.Close()

		return userData
	}

	f := utils.ReadFile(userDataFilePath + userDataFileName)
	if err := json.Unmarshal(f, &userData); err != nil {
		panic(err)
	} else {
		return userData
	}
}

type Account struct {
	Name     string            `json:"name"`
	PWD      string            `json:"pwd"`
	UserName string            `json:"username"`
	Email    string            `json:"email"`
	Info     map[string]string `json:"info"`
	Hisotry  []ChangeHistory   `json:"history"`
}

type ChangeHistory struct {
	ChangeTime time.Time         `json:"change_time"`
	Old        map[string]string `json:"old_value"`
}
