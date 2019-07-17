package handlers

import (
	"fmt"
	"log"
	"pwdbot/structs"
	"pwdbot/utils"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var CmdHandler map[string]func(m *tb.Message)

func init() {
	CmdHandler = make(map[string]func(m *tb.Message))
	CmdHandler["/start"] = start
	CmdHandler["/new"] = new
	CmdHandler["/save"] = save
	CmdHandler["/update"] = update
	CmdHandler["/get"] = get
	CmdHandler["/search"] = search
	CmdHandler["/list"] = list
}

func start(m *tb.Message) {
	help := []string{
		"Hellow, " + m.Sender.Username + "",
		"This is a password management bot.Commands below are now available:",
		"- /new `length` return a random string in specified length",
		"- /save `password` `account name` save password of the account",
		"- /update `acount name` `property name``::``property value` update the specified property of the account name",
		"- /get `account name` get detail of the account",
		"- /list list all accounts",
	}
	_, err := structs.UniqBot.Send(m.Sender, strings.Join(help, "\n\n"), tb.ModeMarkdown)
	if err != nil {
		panic(err)
	}
}

// `/new 32`
func new(m *tb.Message) {
	msgText := m.Text
	splitArr := strings.Split(msgText, " ")
	length := 32 // default length
	if len(splitArr) == 2 {
		if l, err := strconv.ParseInt(splitArr[1], 10, 64); err == nil {
			length = int(l)
		}
	}

	pwd := fmt.Sprintf("new passwd:\n\n\t`%s`", utils.Generator(length))
	msg, err := structs.UniqBot.Send(m.Sender, pwd, tb.ModeMarkdown)
	if err != nil {
		log.Println(msg)
		panic(err)
	}
}

// `/save 456yu#$%^ github`
func save(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 3 {
		structs.UniqBot.Send(m.Sender, "wrong format of parameter")
		return
	}

	pwd := strings.Replace(splitArr[1], " ", "", -1)
	accountName := strings.Replace(splitArr[2], " ", "", -1)

	if pwd == "" || accountName == "" {
		structs.UniqBot.Send(m.Sender, "invalid update pwd or accountName")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if _, ok := userData.Accounts[accountName]; ok && len(userData.Accounts) > 0 {
		structs.UniqBot.Send(m.Sender, fmt.Sprintf("account `%s` already exists, just update it", accountName), tb.ModeMarkdown)
		return
	}

	newAccount := structs.Account{Name: accountName, PWD: pwd, Info: make(map[string]string)}
	userData.Accounts[accountName] = newAccount
	userData.Save()
	log.Println("save account: ", accountName)
	structs.UniqBot.Send(m.Sender, "saved account", tb.ModeMarkdown)
}

// `/update github [username::hulb]`
func update(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 3 {
		structs.UniqBot.Send(m.Sender, "wrong format of parameter")
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	updateInfo := strings.Replace(splitArr[2], " ", "", -1)
	splitArrUpdate := strings.Split(updateInfo, "::")
	if len(splitArrUpdate) != 2 {
		structs.UniqBot.Send(m.Sender, "wrong format of update info")
		return
	}

	updateKey := splitArrUpdate[0]
	updateValue := splitArrUpdate[1]

	if updateKey == "" || updateValue == "" {
		structs.UniqBot.Send(m.Sender, "invalid update key or value")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if account, ok := userData.Accounts[accountName]; ok {
		history := structs.ChangeHistory{ChangeTime: time.Now(), Old: make(map[string]string)}
		switch updateKey {
		case "name":
			history.Old["name"] = account.Name
			account.Name = updateValue
		case "pwd":
			history.Old["pwd"] = account.PWD
			account.PWD = updateValue
		case "username":
			history.Old["username"] = account.UserName
			account.UserName = updateValue
		case "email":
			history.Old["email"] = account.Email
			account.Email = updateValue
		default:
			if v, ok := account.Info[updateKey]; ok {
				history.Old[updateKey] = v
			}

			account.Info[updateKey] = updateValue
		}

		userData.Accounts[accountName] = account
		userData.Save()
		structs.UniqBot.Send(m.Sender, "updated")
		log.Println("update account: ", accountName)
		return
	}

	structs.UniqBot.Send(m.Sender, "nothing updated", tb.ModeMarkdown)
}

// `/get github`
func get(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "wrong format of parameter")
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	if accountName == "" {
		structs.UniqBot.Send(m.Sender, "invalid update key or value")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if account, ok := userData.Accounts[accountName]; ok {
		res := []string{
			fmt.Sprintf("- name: `%s`", account.Name),
			fmt.Sprintf("- username: `%s`", account.UserName),
			fmt.Sprintf("- password: `%s`", account.PWD),
			fmt.Sprintf("- email: `%s`", account.Email),
		}

		for k, v := range account.Info {
			res = append(res, fmt.Sprintf("- %s: `%s`", k, v))
		}

		structs.UniqBot.Send(m.Sender, strings.Join(res, "\n"), tb.ModeMarkdown)
		log.Println("query account: ", accountName)
		return
	}

	structs.UniqBot.Send(m.Sender, "account not found, save it first.")
}

// `/search key`
func search(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "wrong format of parameter")
		return
	}

	searchKey := strings.Replace(splitArr[1], " ", "", -1)
	if searchKey == "" {
		structs.UniqBot.Send(m.Sender, "invalid serach key")
		return
	}

	userData := structs.GetUserData(m.Sender)
	matchAccountName := []string{}
	for _, account := range userData.Accounts {
		switch {
		case strings.Contains(account.Name, searchKey):
			matchAccountName = append(matchAccountName, account.Name)
			continue
		case strings.Contains(account.UserName, searchKey):
			matchAccountName = append(matchAccountName, account.Name)
			continue
		case strings.Contains(account.Email, searchKey):
			matchAccountName = append(matchAccountName, account.Name)
			continue
		default:
			findMatch := false
			for k, v := range account.Info {
				if strings.Contains(k, searchKey) || strings.Contains(v, searchKey) {
					findMatch = true
					break
				}
			}

			if findMatch {
				matchAccountName = append(matchAccountName, account.Name)
				continue
			}
		}
	}

	if len(matchAccountName) > 0 {
		findAccountNames := []string{}
		for _, name := range matchAccountName {
			findAccountNames = append(findAccountNames, fmt.Sprintf("-\t`%s`", name))
		}

		structs.UniqBot.Send(m.Sender, strings.Join(findAccountNames, "\n\n"), tb.ModeMarkdown)
		return
	}

	structs.UniqBot.Send(m.Sender, "not found related account")
}

// `/list`
func list(m *tb.Message) {
	userData := structs.GetUserData(m.Sender)
	allAccountNames := []string{}
	for _, account := range userData.Accounts {
		allAccountNames = append(allAccountNames, fmt.Sprintf("- `%s`\n", account.Name))
	}

	res := strings.Join(allAccountNames, "\n")
	structs.UniqBot.Send(m.Sender, res, tb.ModeMarkdown)
	log.Println("list accounts of user ", userData.User.Username)

}

// `/rm github`
func rm(m *tb.Message) {

}
