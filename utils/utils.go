package utils

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	currentPath := strings.Replace(dir, "\\", "/", -1)

	return currentPath
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
