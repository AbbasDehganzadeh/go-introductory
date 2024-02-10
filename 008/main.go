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
	duration  uint64 //
	maxTicket int
	remTicket int
}

func main() {
	confShai := conference{"Shanghai", time.Date(2024,1,17,19,0, 0, 0, time.Local), 180, 56, 55}
	confChig := conference{"Chicago", time.Date(2024,1,17,21,30,0,0,time.Local), 120, 32, 31}
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	for MAX_TICKET > 0 {
		var name = []string{"John", "Joe", "Jane"}
		var email = []string{"jo@hn.com", "Jo@e.com", "Ja@ne.com"}
		var ticket = []int{7, 2, 4}
		var req_name string
		fmt.Printf("Please enter your first name\t")
		fmt.Scanf("%s", &req_name)
		var req_email string
		fmt.Printf("Please enter your email address\t")
		fmt.Scanf("%s", &req_email)
		// validate name & email
		val_name := len(req_name) > 2 && req_name != "nil"                  // ;-)
		val_email := strings.Contains(req_email, "@") && req_email != "nil" // ;-)
		if !val_name || !val_email {
			fmt.Print("Inputs are wrong", '\n')
			if !val_name {
				fmt.Printf("[%s] is an incorrect name\n", req_name)
			}
			if !val_email {
				fmt.Printf("[%s] is an incorrect email\n", req_email)
			}
			continue
		}
		var conf conference
		req_ticket := 1
		// TODO: let user choose the number of tickets base on the location; create the struct for conference
		resp := fmt.Sprintf("You have booked %d tickets\n", req_ticket)
		var req_loc string
		for true {
			fmt.Println("Please enter your location for the conference: \n")
			// for i := 0; i<len(conferences);i++ {
			// 	fmt.Printf("%d:) '%s'\t", i, conferences[i].location)
			// }
			fmt.Printf("a:)'%s'\tb:)'%s'\t", confShai.location, confChig.location) // HARD coded
			fmt.Print(conf) // ::DEBUG
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
				fmt.Print(conf.remTicket, req_ticket)
				name = append(name, req_name)
				email = append(email, req_email)
				ticket = append(ticket, req_ticket)
				fmt.Print(resp)
				break
			}
		}
		// TODO: Print the number of tickets base on conference location
		fmt.Println("These people `already` took the ticket")
		for i := 0; i < len(ticket); i++ {
			fmt.Printf("[%v, %v, %v ]\n", name[i], email[i], ticket[i])
		}
		remaining := fmt.Sprintf("%d are remained", conf.remTicket)
		var done string
		fmt.Print("Do you want to quit! (y/n)")
		fmt.Scan(&done)
		switch strings.ToLower(done)[0] {
		case 'y', 'q':
			break
		default:
			
		}
		fmt.Print(remaining)
	}
	// fmt.Print(remaining)) //TODO: `remaining` should be accessible here
	fmt.Println("\n\nThanks for using our booking service")
}
