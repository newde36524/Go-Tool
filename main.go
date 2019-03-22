package main

import (
	"fmt"

	"./filetool"
	"./redistool"
)

func main() {

}
func ReadFile(index, pagnum int, filePath string) {
	data, err := filetool.ReadPagingFile(index, pagnum, filePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestRedis() {
	client := new(redistool.RedisClient)
	fmt.Println("连接redis服务端")

	client.Login("ip:port", &redistool.RedisClientOption{
		Password: "password",
	})

	res, err := client.Set("a", "hello")
	fmt.Println("Set", res, err)

	res, err = client.Get("a")
	fmt.Println("Get", res, err)
}
