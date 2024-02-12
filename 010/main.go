package main

import (
	"errors"
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
	confShai := conference{"shanghai", time.Date(2024, 1, 17, 19, 0, 0, 0, time.Local), 180, 56, 55}
	confChig := conference{"chicago", time.Date(2024, 1, 17, 21, 30, 0, 0, time.Local), 120, 32, 31}
	confs := []conference{confShai, confChig}
	fmt.Printf("Welcome to%s booking service\n", PROG_NAME)
	//TODO: check the max number of that one user can book
	for MAX_TICKET > 0 {
		user1 := user{"John", "jo@hn.com", 7}
		user2 := user{"Joe", "jo@e.hn", 2}
		user3 := user{"Jane", "ja@ne.com", 4}
		users := []user{user1, user2, user3}
		var conf conference
		for {
			user_, error_ := getUserInput(confs)
			if error_ != "" {
				log.Println(error_)
				ok := "No"
				fmt.Printf("Do you wanna continue? (y/n):\t")
				fmt.Scan(&ok)
				if strings.ToLower(ok)[0] != 'n' {
					continue
				}
			} else {
				conf.remTicket -= user_.tickets
				users = append(users, user_)
				fmt.Printf("You have booked %d tickets\n", user_.tickets)
				continue
			}
		}
		fmt.Print(printInfo(users))
		fmt.Printf("%d are remained", conf.remTicket)
		done := "n"
		fmt.Print("Do you want to quit! (y/n)")
		fmt.Scanf("%s", &done)
		done = strings.ToLower(done)
		if done[0] == 'y' || done[0] == 'q' {
			break
		}
	}
	//  ) //TODO: `remaining` should be accessible here
	fmt.Println("\n\nThanks for using our booking service")
}

func getUserInput(confs []conference) (user, string) {
	// It gets name, and email of user and store it to struct
	var req_name string
	fmt.Printf("Please enter your first name\t")
	fmt.Scanf("%s", &req_name)
	var req_email string
	fmt.Printf("Please enter your email address\t")
	fmt.Scanf("%s", &req_email)
	var conf conference
	req_ticket := 1
	var req_loc string
	fmt.Println("Please enter your location for the conference: \n")
	fmt.Printf("a:)'Shanghai'\tb:)'Chicago'\t") //, confShai.location, confChig.location) // HARD coded
	fmt.Scanf("%s", &req_loc)
	for _, c := range confs {
		if strings.ToLower(req_loc) == c.location { //TODO: Check for choosing option
			conf = c
		}
	}
	fmt.Printf("Please enter the number of tickets[%v-%v]\t", req_ticket, conf.maxTicket)
	fmt.Scanf("%d", &req_ticket)
	var tmp user = user{req_name, req_email, req_ticket}
	errs := tmp.Validate(conf)
	err := fmt.Sprint("The Inputs are wrong\n")
	for i, e := range errs {
		err += fmt.Sprintf("\t %v) %v\n", i, e)
	}
	return tmp, err
}

func (u *user) Validate(c conference) []error {
	// validate name, email & tickets number
	val_name := len(u.name) > 2 && u.name != "nil" // ;-)
	val_email := strings.Contains(u.email, "@")
	val_ticket := u.tickets > c.remTicket || u.tickets <= 0
	if val_name && val_email && val_ticket {
		return nil
	}
	var errs []error
	strfName := fmt.Sprintf("%v is incorrect", u.name)
	errs = appendIfNot(errs, val_name, strfName)
	strfEmail := fmt.Sprintf("%v is incorrect", u.email)
	errs = appendIfNot(errs, val_email, strfEmail)
	strfTicket := fmt.Sprintf("You have booked %d tickets,\t(ticket capacity;%d )", u.tickets, c.maxTicket)
	errs = appendIfNot(errs, val_ticket, strfTicket)
	return errs
}

func appendIfNot(e []error, b bool, s string) []error {
	// Appends the error to error list if the cond` is false
	var res []error
	if !b {
		res = append(res, errors.New(s))
	}
	return res
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
