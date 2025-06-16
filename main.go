package main

import (
	"context"
	"fmt"

	"github.com/adamwreuben/twiggasdk/twigga"
)

func main() {
	twiggaClient, err := twigga.NewTwiggaClient("./twigga/bongo.json")
	if err != nil {
		fmt.Println("error**: ", err.Error())
		return
	}

	dataToFilter := map[string]any{
		"appSecret": "twigga",
	}

	found, err := twiggaClient.DocumentExists(context.Background(), "Applications", dataToFilter)
	if err != nil {
		fmt.Println("Error*: ", err.Error())
		return
	}

	fmt.Println("FOUND: ", found)

}
