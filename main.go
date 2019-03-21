package main

import (
	"log"

	"./httptool"
)

func main() {
	httpClient := httptool.NewHttpClient()
	httpClient.Header["Content-Type"] = "application/json"
	result, err := httpClient.Post("https://localhost:44316/api/values/PostData", `{
		"name":"tom"
	}`)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
