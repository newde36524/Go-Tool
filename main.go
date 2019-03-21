package main

import (
	"fmt"

	"./filetool"
)

func main() {
	// httpClient := httptool.NewHttpClient()
	// httpClient.Header["Content-Type"] = "application/json"
	// result, err := httpClient.Post("https://localhost:44316/api/values/PostData", `{
	// 	"name":"tom"
	// }`)
	// if err == nil {
	// 	log.Println(result)
	// } else {
	// 	log.Println(err)
	// }

	filePath := `a`
	fmt.Println(getFileSize(filePath))
	// {
	// 	file, err := os.Open(filePath)
	// 	if err != nil {

	// 	}
	// 	defer file.Close()

	// 	bufReader := bufio.NewReader(file)
	// 	for {
	// 		b, e := bufReader.ReadByte()
	// 		if e != nil {
	// 			fmt.Println(e)
	// 			fmt.Println(b)
	// 			// break
	// 		}
	// 		fmt.Println(b)
	// 	}
	// }
	// for i := 0; i < 25; i++ {
	// 	ReadFile(i, 1024, filePath)
	// }

	// ReadFile(1, 3, filePath)
	// ReadFile(0, 699999999999, filePath)
}
func ReadFile(index, pagnum int, filePath string) {
	data, err := filetool.ReadPagingFile(index, pagnum, filePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
