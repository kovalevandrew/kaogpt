package config

import (
	"encoding/json"
	"log"
	"os"
)

type Tokens struct {
	Telegram string
	Gpt      string
	Edenai   string
}

func ReadJsonTokens() Tokens {
	jsonData, err := os.ReadFile("./token/tokens.json")
	if err != nil {
		log.Fatal(err)
	}

	var payload Tokens
	err = json.Unmarshal(jsonData, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}

func GetTelegramToken() string {
	tokens := ReadJsonTokens()
	return tokens.Telegram
}

func GetGPTToken() string {
	tokens := ReadJsonTokens()
	return tokens.Gpt
}

func GetEdenaiToken() string {
	tokens := ReadJsonTokens()
	return tokens.Edenai
}
