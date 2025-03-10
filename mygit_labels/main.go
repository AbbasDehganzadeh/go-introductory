package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	BASE_API_URL  = "https://api.github.com/"
	BASE_AUTH_URL = "https://github.com/"
	TOKEN_FILE    = "./.token"
	OPTIONS_FILE  = "./options.yml"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Err: unable to config file,")
	}
	CLIENT_ID := os.Getenv("CLIENT_ID")

	// get username, otherwise login user
	accessToken, refreshToken, err := GetTokens(TOKEN_FILE)
	if err != nil {
		log.Panicln(accessToken, refreshToken, "\t", err)
	}
	username, err := GetUsername(accessToken)
	if username != "" {
		fmt.Printf("Hello %s\n", username)
	} else if err != nil {
		log.Panicln(err)
	} else {
		accessToken, refreshToken, err = Refresh(CLIENT_ID, refreshToken)
		if err != nil {
			log.Fatalf("Unable to Login, pls try later! %s\n", err.Error())
		}
		if accessToken == "" {
			accessToken, refreshToken, err = Login(CLIENT_ID)
			if err != nil {
				log.Fatalf("Unable to Login, pls try later! %s\n", err.Error())
			}
		}
		fmt.Println(TOKEN_FILE, accessToken, refreshToken)
		_ = SaveTokens(TOKEN_FILE, accessToken, refreshToken)
	}

	args := &Arguments{Items: "issues",
		Limit: 32, Max: 5,
		Output: "result.data", Format: Normal,
	}
	args = ParseArguments(args)

	//get issues, and prs > output the result
	var repos []Repo
	repos, err = RetrieveRepos(accessToken, *args)
	if err != nil {
		log.Fatalf("Could not retrieve topics: %s\n", err.Error())
	}
	for _, repo := range repos {
		repo.Display(args.Format)
	}
	if args.Output == "" {
		SaveFile(repos, args.Output)
	}
}
