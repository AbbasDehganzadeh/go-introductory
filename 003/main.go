package main

import (
	"fmt"
	"log"
)

const (
	PROG_NAME = "booking.ir"
)

var (
	AVAIL_TICKET = 32
)

func main() {
	var f_name = [32]string{"John", "Joe", "Jane"}
	var e_mail = [32]string{"jo@hn.com", "Jo@e.com", "Ja@ne.com"}
	var ticket = [32]int16{7, 2, 4}
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	var req_fname string
	fmt.Printf("Please enter your first name")
	fmt.Scanf("%s", &req_fname)
	var req_email string
	fmt.Printf("Please enter your email address")
	fmt.Scanf("%s", &req_email)
	req_ticket := 1
	fmt.Printf("Please enter the number of tickets[%v]\t", req_ticket)
	fmt.Scanf("%d", &req_ticket)
	resp := fmt.Sprintf("You have booked %d tickets\n", req_ticket)

	for true {
		if req_ticket > AVAIL_TICKET {
			log.Printf("You have booked %d tickets,\t(ticket capacity;%d )", req_ticket, AVAIL_TICKET)
		} else {
			AVAIL_TICKET -= req_ticket
			fmt.Print(resp)
			break
		}
	}
	for i := 0; i < 32; i++ {
		fmt.Println(f_name[i], e_mail[i], ticket[i])
	}
	fmt.Println("\n\nThanks for using our booking service")
}
