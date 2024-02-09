package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	PROG_NAME  = "booking.ir"
	MAX_TICKET = 255
)

var (
	AVAIL_TICKET = 32
)

func main() {
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	for AVAIL_TICKET > 0 {
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
				fmt.Printf("[%s] is an incorrected name\n", req_name)
			}
			if !val_email {
				fmt.Printf("[%s] is an incorrected email\n", req_email)
			}
			continue
		}
		req_ticket := 1
		resp := fmt.Sprintf("You have booked %d tickets\n", req_ticket)

		for true {
			fmt.Printf("Please enter the number of tickets[%v]\t", req_ticket)
			fmt.Scanf("%d", &req_ticket)
			if req_ticket > AVAIL_TICKET {
				log.Printf("You have booked %d tickets,\t(ticket capacity;%d )", req_ticket, AVAIL_TICKET)
				ok := "No"
				fmt.Printf("Do you wanna continue? (y/n):\t")
				fmt.Scan(&ok)
				if strings.ToLower(ok)[0] == 'n' {
					break
				}
			} else {
				AVAIL_TICKET -= req_ticket
				name = append(name, req_name)
				email = append(email, req_email)
				ticket = append(ticket, req_ticket)
				fmt.Print(resp)
				break
			}
		}
		fmt.Println("These people `already` took the ticket")
		for i := 0; i < len(ticket); i++ {
			fmt.Printf("[%v, %v, %v ]\n", name[i], email[i], ticket[i])
		}
		var done string
		fmt.Print("Do you want to quit! (y/n)")
		fmt.Scan(&done)
		if strings.ToLower(done)[0] == 'y' {
			break
		}
		fmt.Printf("%d are remained", AVAIL_TICKET)
	}
	fmt.Println("\n\nThanks for using our booking service")
}
