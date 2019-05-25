package main

import (
	"fmt"

	"github.com/newde36524/Go-Tool/redirectserver"
)

func main() {
	srv := new(redirectserver.Redirect)
	srv.Run("127.0.0.1:12336", "127.0.0.1:12337")
	fmt.Scanln()
}
