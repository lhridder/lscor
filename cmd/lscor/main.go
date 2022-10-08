package main

import (
	"log"
	"lscor/api"
	"lscor/config"
	"lscor/corero"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %s", err)
		return
	}

	err = corero.FetchToken()
	if err != nil {
		log.Printf("Failed to fetch auth token: %s", err)
		return
	}

	api.Start()

	err = corero.Logout()
	if err != nil {
		log.Printf("Failed to log out: %s", err)
		return
	}
}
