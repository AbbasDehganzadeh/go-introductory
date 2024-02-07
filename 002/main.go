package main

import (
	"fmt"
	"log"
)

const (
	PROG_NAME    = "booking.ir"
	AVAIL_TICKET = 32
)

func main() {
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	req_ticket := 1
	fmt.Printf("Please enter the number of tickets[%v]\t", req_ticket)
	fmt.Scanf("%d", &req_ticket)
	resp := fmt.Sprintf("You have booked %d tickets\n", req_ticket)
	// TODO: validate the number of tickets
	if req_ticket > AVAIL_TICKET {
		log.Fatalf("You have booked %d tickets,\t(ticket capacity;%d )", req_ticket, AVAIL_TICKET)
	}
	fmt.Print(resp)
	fmt.Println("Thanks for using our booking service")
}
