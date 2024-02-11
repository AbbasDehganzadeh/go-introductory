package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	PROG_NAME  = "booking.ir"
	MAX_TICKET = 255
)

type conference struct {
	location  string
	initTime  time.Time
	duration  uint64 // is time.Duration better
	maxTicket int
	remTicket int
}

type user struct {
	name    string // TODO name should be a slice of string
	email   string
	tickets int
}

func main() {
	confShai := conference{"Shanghai", time.Date(2024, 1, 17, 19, 0, 0, 0, time.Local), 180, 56, 55}
	confChig := conference{"Chicago", time.Date(2024, 1, 17, 21, 30, 0, 0, time.Local), 120, 32, 31}
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	//TODO: check the max number of that one user can book
	for MAX_TICKET > 0 {
		user1 := user{"John", "jo@hn.com", 7}
		user2 := user{"Joe", "jo@e.hn", 2}
		user3 := user{"Jane", "ja@ne.com", 4}
		users := []user{user1, user2, user3}
		user_ := getUserInput()
		//TODO validate function/method
		// validate name & email
		val_name := len(user_.name) > 2 && user_.name != "nil" // ;-)
		val_email := strings.Contains(user_.email, "@")
		if !val_name || !val_email {
			fmt.Print("Inputs are wrong", '\n')
			if !val_name {
				fmt.Printf("[%s] is an incorrect name\n", user_.name)
			}
			if !val_email {
				fmt.Printf("[%s] is an incorrect email\n", user_.email)
			}
			continue
		}
		var conf conference
		req_ticket := 1

		var req_loc string
		for true {
			fmt.Println("Please enter your location for the conference: \n")
			// for i := 0; i<len(conferences);i++ {
			// 	fmt.Printf("%d:) '%s'\t", i, conferences[i].location)
			// }
			fmt.Printf("a:)'%s'\tb:)'%s'\t", confShai.location, confChig.location) // HARD coded
			fmt.Print(conf)                                                        // ::DEBUG
			fmt.Scanf("%s", &req_loc)
			switch {
			case req_loc == "a":
				conf = confShai
			case req_loc == "b":
				conf = confChig
			}
			fmt.Print(conf) // ::DEBUG
			fmt.Printf("Please enter the number of tickets[%v]\t", req_ticket)
			fmt.Scanf("%d", &req_ticket)
			fmt.Println(req_ticket)
			if req_ticket > conf.remTicket || req_ticket <= 0 {
				log.Printf("You have booked %d tickets,\t(ticket capacity;%d )", req_ticket, conf.maxTicket)
				ok := "No"
				fmt.Printf("Do you wanna continue? (y/n):\t")
				fmt.Scan(&ok)
				if strings.ToLower(ok)[0] != 'n' {
					break
				}
			} else {
				conf.remTicket -= req_ticket
				users = append(users, user_)
				fmt.Printf("You have booked %d tickets\n", req_ticket)
				break
			}
		}
		fmt.Print(printInfo(users))
		remaining := fmt.Sprintf("%d are remained", conf.remTicket)
		var done string
		fmt.Print("Do you want to quit! (y/n)")
		fmt.Scan(&done)
		switch strings.ToLower(done)[0] {
		case 'y', 'q':
			break
		default:
			fmt.Print(remaining)

		}
	}
	// fmt.Print(remaining)) //TODO: `remaining` should be accessible here
	fmt.Println("\n\nThanks for using our booking service")
}

func getUserInput() user {
	// It gets name, and email of user and store it to struct
	var req_name string
	fmt.Printf("Please enter your first name\t")
	fmt.Scanf("%s", &req_name)
	var req_email string
	fmt.Printf("Please enter your email address\t")
	fmt.Scanf("%s", &req_email)
	var tmp user = user{req_name, req_email, 0}
	return tmp
}

func printInfo(u []user) string {
	// It returns the nice formatting of information about the conference
	res := fmt.Sprintln("These people `already` took the ticket")
	for i, _ := range u {
		res += fmt.Sprintf(res, "[%v, %v, %v ]\n", u[i].name, u[i].email, u[i].tickets)
		//! not use Appendf because of backwards compatibility
	}
	return res
}
