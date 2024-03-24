package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TitleCase(key string) string {

	splitKey := strings.Split(key, ".")
	newSplitKey := []string{}
	for _, key := range splitKey {
		newSplitKey = append(newSplitKey, cases.Title(language.English).String(key))
	}
	newKey := strings.Join(newSplitKey, ".")

	return newKey
}

func Stringify(value interface{}) {

	result, err := json.Marshal(value)
	if err != nil {
		fmt.Println("error when print log", err)
	}

	fmt.Println(string(result))

}

func GetRealKey(key string) (string, int) {
	splitKey := strings.Split(key, ".")
	lastIndex := len(splitKey) - 1
	lastKey := splitKey[lastIndex]

	lasKeyIndex, err := strconv.Atoi(lastKey)

	if err != nil {
		return key, 0
	}

	newKey := strings.Join(RemoveIndex(splitKey, lastIndex), ".")

	return newKey, lasKeyIndex
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
