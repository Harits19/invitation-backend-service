package util

import (
	"fmt"
	"log"
	"net/http"
)

func Connected() {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("internet status connected")
}
