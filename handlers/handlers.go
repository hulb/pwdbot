package handlers

import (
	"fmt"
	"log"
	"pwdbot/structs"
	"pwdbot/utils"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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
	CmdHandler["/saerch"] = search
	CmdHandler["/list"] = list
	CmdHandler["/rm"] = rm
	CmdHandler["/addtotp"] = addTOTP
	CmdHandler["/gettotp"] = getTOTP
}

func start(m *tb.Message) {
	help := []string{
		"Hellow, " + m.Sender.Username + "",
		"This is a password management bot.Commands below are now available:",
		"- /new `[length](optional)` return a random string in specified length",
		"- /save `[password]` `[account name]` save password of the account",
		"- /update `[acount name]``.``[property name]``=``[value]` update the specified property of the account name",
		"- /get `[account name]` get detail of the account",
		"- /search `[search key]` fuzzy search accounts that match the key",
		"- /rm `[account name]` delete the account",
		"- /addtotp `[account name]` add TOTP key for the account",
		"- /gettotp `[account name]` get a TOTP code for the account",
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

	pwd := fmt.Sprintf("I generate a random password for you:\n\n\t`%s`", utils.Generator(length))
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
		structs.UniqBot.Send(m.Sender,
			"Incorrect format of command /save.\nYou should use it like: /save `[password]` `[account name]`.\nA blank is between `[password]` and `[account name]`.\nFor example, /save `mypassword` `google`",
			tb.ModeMarkdown)
		return
	}

	pwd := strings.Replace(splitArr[1], " ", "", -1)
	accountName := strings.Replace(splitArr[2], " ", "", -1)

	if pwd == "" || accountName == "" {
		structs.UniqBot.Send(m.Sender, "invalid pwd or accountName")
		return
	}

	if strings.Contains(accountName, ".") {
		structs.UniqBot.Send(m.Sender, "Invalid account name. The character `.` can not exists in account name.", tb.ModeMarkdown)
		return
	}

	userData := structs.GetUserData(m.Sender)
	if _, ok := userData.Accounts[accountName]; ok && len(userData.Accounts) > 0 {
		structs.UniqBot.Send(m.Sender, fmt.Sprintf("Account `%s` already exists! You can just use /update ot update it or /get to overview it.", accountName), tb.ModeMarkdown)
		return
	}

	newAccount := structs.Account{Name: accountName, PWD: pwd, Info: make(map[string]string)}
	userData.Accounts[accountName] = newAccount
	userData.Save()
	log.Printf("An account named %s has been saved.", accountName)
	structs.UniqBot.Send(m.Sender, fmt.Sprintf("An account named `%s` has been saved.\nYou can use /get to overview it.", accountName), tb.ModeMarkdown)
}

// `/update github.username=hulb`
func update(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "wrong format of parameter")
		return
	}

	params := strings.Replace(splitArr[1], " ", "", -1)
	paramsArr := strings.Split(params, ".")
	accountName := paramsArr[0]

	propertyArr := strings.Split(paramsArr[1], "=")
	if len(propertyArr) != 2 {
		structs.UniqBot.Send(m.Sender, "update parameters should in format like `[property]``=``[value]`", tb.ModeMarkdown)
		return
	}

	updateKey := propertyArr[0]
	updateValue := propertyArr[1]
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

		// keep old value as history
		if len(history.Old) > 0 {
			allHistory := userData.Accounts[accountName].Hisotry
			allHistory = append(allHistory, history)
			account.Hisotry = allHistory
		}

		userData.Accounts[accountName] = account
		userData.Save()
		structs.UniqBot.Send(m.Sender, fmt.Sprintf("Property `%s` of account `%s` has been updated.", updateKey, accountName), tb.ModeMarkdown)
		log.Println(fmt.Sprintf("Property `%s` of account `%s` has been updated.", updateKey, accountName))
		return
	}

	structs.UniqBot.Send(m.Sender, "nothing updated", tb.ModeMarkdown)
}

// `/get github`
func get(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "Incorrect format of command /get.\nYou should use it like: /get `[account name]`\n", tb.ModeMarkdown)
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	if accountName == "" {
		structs.UniqBot.Send(m.Sender, "invalid account name")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if account, ok := userData.Accounts[accountName]; ok {
		structs.UniqBot.Send(m.Sender, account.String(), tb.ModeMarkdown)
		log.Println("query account: ", accountName)
		return
	}

	structs.UniqBot.Send(m.Sender, "Account not found, save it first.")
}

// `/search key`
func search(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "Incorrect format of command /search.\nYou should use it like: /search `[search key]`\n")
		return
	}

	searchKey := strings.Replace(splitArr[1], " ", "", -1)
	if searchKey == "" {
		structs.UniqBot.Send(m.Sender, "Invalid serach key")
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

	switch len(matchAccountName) {
	case 0:
		structs.UniqBot.Send(m.Sender, "No account found")
	case 1:
		account := userData.Accounts[matchAccountName[0]]
		structs.UniqBot.Send(m.Sender, account.String(), tb.ModeMarkdown)
	default:
		findAccountNames := []string{}
		for _, name := range matchAccountName {
			findAccountNames = append(findAccountNames, fmt.Sprintf("-\t`%s`", name))
		}

		structs.UniqBot.Send(m.Sender, strings.Join(findAccountNames, "\n\n"), tb.ModeMarkdown)
	}
}

// `/list`
func list(m *tb.Message) {
	userData := structs.GetUserData(m.Sender)
	allAccountNames := []string{}
	for _, account := range userData.Accounts {
		allAccountNames = append(allAccountNames, fmt.Sprintf("- `%s`\n", account.Name))
	}

	if len(allAccountNames) == 0 {
		structs.UniqBot.Send(m.Sender, "No account found for current user, you can /save one first.", tb.ModeMarkdown)
		return
	}
	res := strings.Join(allAccountNames, "\n")
	structs.UniqBot.Send(m.Sender, res, tb.ModeMarkdown)
	log.Println("list accounts of user ", userData.User.Username)

}

// `/rm github`
func rm(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "Incorrect format of command /rm.\nYou should use it like: /rm `[account name]`\n", tb.ModeMarkdown)
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	if accountName == "" {
		structs.UniqBot.Send(m.Sender, "Invalid account name")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if _, ok := userData.Accounts[accountName]; ok {
		delete(userData.Accounts, accountName)
		userData.Save()
		log.Println(fmt.Sprintf("Account %s has been deleted.", accountName))
	}

	structs.UniqBot.Send(m.Sender, fmt.Sprintf("Account `%s` has been deleted.", accountName), tb.ModeMarkdown)
}

// `/addtotp github xxx`
func addTOTP(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 3 {
		structs.UniqBot.Send(m.Sender,
			"Incorrect format of command /addtotp.\nYou should use it like: /addtotp `[account name]` `[totp key uri]`.\nA blank is between `[account name]` and `[totp key uri]`.",
			tb.ModeMarkdown)
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	totpKeyURI := strings.Replace(splitArr[2], " ", "", -1)

	if accountName == "" || totpKeyURI == "" {
		structs.UniqBot.Send(m.Sender, "invalid account name or totp key uri")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if account, ok := userData.Accounts[accountName]; ok {
		if account.KeyUri != "" {
			history := structs.ChangeHistory{ChangeTime: time.Now(), Old: make(map[string]string)}
			history.Old["keyuri"] = account.KeyUri
			account.Hisotry = append(account.Hisotry, history)
		}

		account.KeyUri = totpKeyURI
		userData.Accounts[accountName] = account
		userData.Save()

		structs.UniqBot.Send(m.Sender, fmt.Sprintf("The totp key URI of account `%s` has been updated.", accountName), tb.ModeMarkdown)
		log.Println(fmt.Sprintf("The totp key URI of account `%s` has been updated.", accountName))
		return
	}

	structs.UniqBot.Send(m.Sender, "Account not found")
}

// `/gettotp github`
func getTOTP(m *tb.Message) {
	splitArr := strings.Split(m.Text, " ")
	if len(splitArr) != 2 {
		structs.UniqBot.Send(m.Sender, "Incorrect format of command /gettotp.\nYou should use it like: /gettotp `[account name]`\n", tb.ModeMarkdown)
		return
	}

	accountName := strings.Replace(splitArr[1], " ", "", -1)
	if accountName == "" {
		structs.UniqBot.Send(m.Sender, "invalid account name")
		return
	}

	userData := structs.GetUserData(m.Sender)
	if account, ok := userData.Accounts[accountName]; ok {
		if otpKey, err := otp.NewKeyFromURL(account.KeyUri); err != nil {
			log.Printf("Error occures when trying to get OTP key of account %s, err: %s\n", accountName, err)
			structs.UniqBot.Send(m.Sender, "Get TOTP code fail! Please try to reset the TOTP key URI of the account.", tb.ModeMarkdown)
			return
		} else {
			if code, err := totp.GenerateCode(otpKey.Secret(), time.Now()); err != nil {
				log.Printf("Error occures when trying to get TOTP code of account %s, err: %s\n", accountName, err)
				structs.UniqBot.Send(m.Sender, "Get TOTP code fail! Please try to reset the TOTP key URI of the account.", tb.ModeMarkdown)
				return
			} else {
				structs.UniqBot.Send(m.Sender, fmt.Sprintf("Code below is expired in 30s:\n\n\t`%s`", code), tb.ModeMarkdown)
				return
			}
		}
	}

	structs.UniqBot.Send(m.Sender, "Account not found, save it first.")
}
