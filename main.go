package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	BOTTOKEN string
	FILENAME string
)

func init() {
	BOTTOKEN = os.Getenv("BOTTOKEN")
	if BOTTOKEN == "" {
		log.Fatalln("bot token required.exit")
	}

	FILENAME = os.Getenv("FILENAME")
	if FILENAME == "" {
		FILENAME = "pwdbotpwd.pwd"
	}

	log.Println("user default file name: ", FILENAME)
}

type Account struct {
	Name     string            `json:"name"`
	PWD      string            `json:"pwd"`
	UserName string            `json:"username"`
	Email    string            `json:"email"`
	Info     map[string]string `json:"info"`
	Belong2  User              `json:"belong2"`
	Hisotry  []ChangeHistory   `json:"history"`
}

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type ChangeHistory struct {
	ChangeTime time.Time         `json:"change_time"`
	Old        map[string]string `json:"old_value"`
}

func main() {
	if !Exists(GetCurrentDirectory() + FILENAME) {
		err := os.Mkdir(GetCurrentDirectory(), os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}

		f, err := os.Create(GetCurrentDirectory() + FILENAME)
		if err != nil {
			log.Println(err.Error())
		} else {
			_, err = f.Write([]byte("[]"))
			if err != nil {
				panic(err)
			}
		}

		f.Close()
	}

	// DEBUG set http proxy
	// proxy, _ := url.Parse("http://127.0.0.1:7890")
	// tr := &http.Transport{
	// 	Proxy:           http.ProxyURL(proxy),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{
	// 	Transport: tr,
	// }

	b, err := tb.NewBot(tb.Settings{
		Token:  BOTTOKEN,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		help := []string{
			"Hellow, " + m.Sender.Username + "",
			"This is a password management bot.Commands below are now available:",
			"- /new `length` return a random string in specified length",
			"- /save `password` `account name` save password of the account",
			"- /update `acount name` `property name``::``property value` update the specified property of the account name",
			"- /get `account name` get detail of the account",
			"- /list list all accounts",
		}
		_, err := b.Send(m.Sender, strings.Join(help, "\n\n"), tb.ModeMarkdown)
		if err != nil {
			panic(err)
		}
	})

	// `/new 32`，
	b.Handle("/new", func(m *tb.Message) {
		msgText := m.Text
		splitArr := strings.Split(msgText, " ")
		length := 32 // default length
		if len(splitArr) == 2 {
			if l, err := strconv.ParseInt(splitArr[1], 10, 64); err == nil {
				length = int(l)
			}
		}

		pwd := fmt.Sprintf("new passwd:\n\n\t`%s`", Generator(length))
		msg, err := b.Send(m.Sender, pwd, tb.ModeMarkdown)
		if err != nil {
			log.Println(msg)
			panic(err)
		}
	})

	// `/save 456yu#$%^ github`
	b.Handle("/save", func(m *tb.Message) {

		user := User{
			ID:           m.Sender.ID,
			FirstName:    m.Sender.FirstName,
			LastName:     m.Sender.LastName,
			UserName:     m.Sender.Username,
			LanguageCode: m.Sender.LanguageCode,
		}

		splitArr := strings.Split(m.Text, " ")
		if len(splitArr) != 3 {
			b.Send(m.Sender, "wrong format of parameter")
			return
		}

		pwd := strings.Replace(splitArr[1], " ", "", -1)
		accountName := strings.Replace(splitArr[2], " ", "", -1)

		if pwd == "" || accountName == "" {
			b.Send(m.Sender, "invalid update pwd or accountName")
			return
		}

		exists := false
		f := ReadFile(GetCurrentDirectory() + FILENAME)
		accounts := []Account{}
		if err := json.Unmarshal(f, &accounts); err == nil {
			for _, a := range accounts {
				if a.Name == accountName && a.Belong2.ID == m.Sender.ID {
					exists = true
					b.Send(m.Sender, fmt.Sprintf("account `%s` already exists, just update it", a.Name), tb.ModeMarkdown)
				}
			}

			if !exists {
				newAccount := Account{Name: accountName, PWD: pwd, Info: make(map[string]string), Belong2: user}
				accounts = append(accounts, newAccount)
				if acbytes, err := json.Marshal(accounts); err == nil {
					WriteFile(GetCurrentDirectory()+FILENAME, acbytes)
					b.Send(m.Sender, "saved")
					log.Println("saved account: ", accountName)
				} else {
					panic(err)
				}
			}
		}
	})

	// `/update github [username::hulb]`
	b.Handle("/update", func(m *tb.Message) {
		splitArr := strings.Split(m.Text, " ")
		if len(splitArr) != 3 {
			b.Send(m.Sender, "wrong format of parameter")
			return
		}

		accountName := strings.Replace(splitArr[1], " ", "", -1)
		updateInfo := strings.Replace(splitArr[2], " ", "", -1)
		splitArrUpdate := strings.Split(updateInfo, "::")
		if len(splitArrUpdate) != 2 {
			b.Send(m.Sender, "wrong format of update info")
			return
		}

		updateKey := splitArrUpdate[0]
		updateValue := splitArrUpdate[1]

		if updateKey == "" || updateValue == "" {
			b.Send(m.Sender, "invalid update key or value")
			return
		}

		updated := false
		f := ReadFile(GetCurrentDirectory() + FILENAME)

		accounts := []Account{}
		if err := json.Unmarshal(f, &accounts); err == nil {
			for idx, a := range accounts {
				if a.Name == accountName && a.Belong2.ID == m.Sender.ID {
					history := ChangeHistory{ChangeTime: time.Now(), Old: make(map[string]string)}
					switch updateKey {
					case "name":
						history.Old["name"] = a.Name
						accounts[idx].Name = updateValue
					case "pwd":
						history.Old["pwd"] = a.PWD
						accounts[idx].PWD = updateValue
					case "username":
						history.Old["username"] = a.UserName
						accounts[idx].UserName = updateValue
					case "email":
						history.Old["email"] = a.Email
						accounts[idx].Email = updateValue
					default:
						if v, ok := accounts[idx].Info[updateKey]; ok {
							history.Old[updateKey] = v
						}

						accounts[idx].Info[updateKey] = updateValue
					}

					if len(history.Old) > 0 {
						updated = true
						accounts[idx].Hisotry = append(accounts[idx].Hisotry, history)
					}
				}
			}
		}

		if updated {
			if acbytes, err := json.Marshal(accounts); err == nil {
				WriteFile(GetCurrentDirectory()+FILENAME, acbytes)
			}
		}

		if updated {
			b.Send(m.Sender, "updated")
			log.Println("update account: ", accountName)
		} else {
			b.Send(m.Sender, "nothing updated", tb.ModeMarkdown)
		}
	})

	// `/get github`
	b.Handle("/get", func(m *tb.Message) {
		splitArr := strings.Split(m.Text, " ")
		if len(splitArr) != 2 {
			b.Send(m.Sender, "wrong format of parameter")
			return
		}

		accountName := strings.Replace(splitArr[1], " ", "", -1)
		if accountName == "" {
			b.Send(m.Sender, "invalid update key or value")
			return
		}

		var account *Account
		f := ReadFile(GetCurrentDirectory() + FILENAME)
		accounts := []Account{}
		if err := json.Unmarshal(f, &accounts); err == nil {
			for _, a := range accounts {
				if a.Name == accountName && a.Belong2.ID == m.Sender.ID {
					account = &a
				}
			}
		} else {
			panic(err)
		}

		if account != nil {
			res := []string{
				fmt.Sprintf("- name: `%s`", account.Name),
				fmt.Sprintf("- username: `%s`", account.UserName),
				fmt.Sprintf("- password: `%s`", account.PWD),
				fmt.Sprintf("- email: `%s`", account.Email),
			}

			for k, v := range account.Info {
				res = append(res, fmt.Sprintf("- %s: `%s`", k, v))
			}

			b.Send(m.Sender, strings.Join(res, "\n"), tb.ModeMarkdown)
			log.Println("query account: ", accountName)
		}
	})

	// `/search key`
	b.Handle("/search", func(m *tb.Message) {

	})

	// `/list`
	b.Handle("/list", func(m *tb.Message) {
		allAccountNames := []string{}
		f := ReadFile(GetCurrentDirectory() + FILENAME)
		accounts := []Account{}
		if err := json.Unmarshal(f, &accounts); err == nil {
			if len(accounts) == 0 {
				log.Println("no accounts")
			}
			for _, a := range accounts {
				if a.Belong2.ID == m.Sender.ID {
					allAccountNames = append(allAccountNames, fmt.Sprintf("- `%s`\n", a.Name))
				}
			}
		} else {
			panic(err)
		}

		res := strings.Join(allAccountNames, "\n")
		b.Send(m.Sender, res, tb.ModeMarkdown)
		log.Println("list accounts")
	})

	// `/rm github`
	b.Handle("/rm", func(m *tb.Message) {

	})

	b.Start()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	currentPath := strings.Replace(dir, "\\", "/", -1)

	return currentPath + "/data/"
}

func ReadFile(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return f

}

func WriteFile(path string, data []byte) {
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		panic(err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// Generator generates a random string with specified length
func Generator(length int) string {
	number := []int{}
	upper := []int{}
	lower := []int{}
	special := []int{}

	for i := 65; i <= 90; i++ {
		upper = append(upper, i)
	}

	for i := 97; i <= 122; i++ {
		lower = append(lower, i)
	}

	for i := 48; i <= 57; i++ {
		number = append(number, i)
	}

	for i := 33; i <= 47; i++ {
		special = append(special, i)
	}

	for i := 58; i <= 64; i++ {
		special = append(special, i)
	}

	for i := 91; i < 96; i++ {
		special = append(special, i)
	}

	for i := 123; i <= 126; i++ {
		special = append(special, i)
	}

	seed := [][]int{number, upper, lower, special}
	result := []string{}
	for len(result) < length {
		arr := seed[rand.Intn(len(seed))]
		result = append(result, string(arr[rand.Intn(len(arr))]))
	}

	newPWD := strings.Join(result, "")
	log.Println("generate new password:", newPWD)

	return newPWD
}
