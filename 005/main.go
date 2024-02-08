package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	PROG_NAME = "booking.ir"
	MAX_TICKET = 255
)

var (
	AVAIL_TICKET = 32
)

func main() {
	for AVAIL_TICKET > 0 {
		var f_name = []string{"John", "Joe", "Jane"}
		var e_mail = []string{"jo@hn.com", "Jo@e.com", "Ja@ne.com"}
		var ticket = []int{7, 2, 4}
		fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
		var req_fname string
		fmt.Printf("Please enter your first name\t")
		fmt.Scanf("%s", &req_fname)
		var req_email string
		fmt.Printf("Please enter your email address\t")
		fmt.Scanf("%s", &req_email)
		var req_ticket int= 1
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
				f_name = append(f_name, req_fname)
				e_mail = append(e_mail, req_email)
				ticket = append(ticket, req_ticket)
				fmt.Print(resp)
				break
			}
		}
		fmt.Println("These people `already` took the ticket")
		for i := 0; i < len(ticket); i++ {
			fmt.Printf("[%v, %v, %v ]\n",f_name[i], e_mail[i], ticket[i])
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

