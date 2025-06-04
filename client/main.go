package main

import (
	"licenser/appchecker"
)

func init() {
	go appchecker.Validate("App")
}

func main() {

	select {}
	//log.Println("✅ License validated. Starting service...")

	// Start your protected service here
}
