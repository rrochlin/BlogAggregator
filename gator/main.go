package main

import (
	"fmt"
	"github.com/rrochlin/BlogAggregator/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config")
		fmt.Println(err)
		return
	}

	cfg.SetUser("lane")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println("error reading config")
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.CurrentUserName)
	fmt.Println(cfg.DBUrl)

}
