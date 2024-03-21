package util

import (
	"encoding/json"
	"fmt"
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

func Log(value interface{}) {

	result, err := json.Marshal(value)
	if err != nil {
		fmt.Println("error when print log", err)
	}

	fmt.Println(string(result))

}
