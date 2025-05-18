package main

import (
	"fmt"

	"sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger

	// TODO: init database

	// TODO: init router

	// TODO: init server

}
