package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/a-h/tpgit/targetprocess"
)

var url = flag.String("url", "", "Set to the root address of your TargetProcess account, e.g. https://example.tpondemand.com")
var username = flag.String("username", "", "Sets the username to use to authenticate against TargetProcess.")
var password = flag.String("password", "", "Sets the password to use to authenticate against TargetProcess.")
var entity = flag.Int("entity", 0, "Entity to add comment to.")

func main() {
	flag.Parse()

	if *url == "" {
		fmt.Println("url flag missing")
		return
	}
	if *username == "" {
		fmt.Println("username flag missing")
		return
	}
	if *password == "" {
		fmt.Println("password flag missing")
		return
	}
	if *entity == 0 {
		fmt.Println("entity flag missing")
		return
	}

	api := targetprocess.NewAPI(*url, *username, *password)
	msg := fmt.Sprintf("test message %v", time.Now())
	fmt.Printf("Adding comment to entity %d: %s\n", *entity, msg)
	err := api.Comment(*entity, msg)

	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
