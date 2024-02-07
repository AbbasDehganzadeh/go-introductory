package main

import "fmt"

const (
	PROG_NAME    = "booking.ir"
	AVAIL_TICKET = 32
)

func main() {
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	req_ticket := 0
	fmt.Printf("Please enter the number of tickets\t")
	fmt.Scanf("%d", &req_ticket)
	resp := fmt.Sprintf("You have booked %d tickets\n", req_ticket)
	// TODO: validate the number of tickets

	fmt.Print(resp)
	fmt.Println("Thanks for using our booking service")
}
