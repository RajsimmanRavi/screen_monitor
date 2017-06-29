package main

import (
	"bytes"
	//"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const CURR_TMP_FILE = "/tmp/curr_sessions.log"
const PAST_TMP_FILE = "/tmp/past_sessions.log"
const LOG_FILE = "/tmp/screen_monitor_logs"
const SCRIPT = "/home/perplexednoob/Documents/go/src/screen_monitor/get_procs.sh"
const SCRIPT_ARG = "get_procs"

const FROM_EMAIL = "noreply@savinetwork.ca"
const TO_EMAIL = "rajsimmanr@savinetwork.ca"
const PWD = "WEe2>4PT"

// Author: Rajsimman Ravi
// Any questions, contact: rajsimmanr@savinetwork.ca

//Execute bash command in go: Ref > https://stackoverflow.com/questions/12891294/google-golang-exec-exit-status-2-and-1
//Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid
//Catch Signals: Ref > https://stackoverflow.com/questions/18106749/golang-catch-signals

// Function to check for errors. Copied from Golang example
func check(e error) {
	if e != nil {
		notify_user("exited")
		log_entry(e.Error())
		panic(e)
	}
}

// Function to log errors, warnings etc. to file
func log_entry(str string) {

	f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check(err)

	defer f.Close()

	log.SetOutput(f)
	log.Println(str)

}

// Function to execute bash command
func exec_cmd(comd string) string {

	cmd := exec.Command("/bin/bash", SCRIPT, comd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)

	return out.String()
}

// Function to email user regarding program exit
func notify_user(status string) {

	hostname, err := os.Hostname()
	check(err)

	var subject = "Screen Monitor Program " + status + " on host: " + hostname
	var content = "Hello,<br><br>"
	content += "This is to notify you that Screen Monitor program has " + status + " on this host: <b>" + hostname + "</b>.<br><br>"
	content += "Please review the status of the program.<br><br>"
	content += "Regards,<br>"
	content += "SAVI Testbed Team"

	send_email(subject, content, false)

}

// Function to send email
func send_email(subject string, content string, attach bool) {

	m := gomail.NewMessage()
	m.SetHeader("From", FROM_EMAIL)
	m.SetHeader("To", TO_EMAIL)
	m.SetHeader("Subject", subject)
	content = strings.Replace(content, "\n", "<br>", -1)
	m.SetBody("text/html", content)

	if attach {
		m.Attach(PAST_TMP_FILE)
		m.Attach(CURR_TMP_FILE)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, FROM_EMAIL, PWD)

	var err = d.DialAndSend(m)
	check(err)

}

//Function to write the screen sessions to file
func write_to_file(file string, data string) {

	f, err := os.Create(file)
	check(err)

	defer f.Close()

	_, err = f.WriteString(data + "\n")
	check(err)

	f.Sync()

	f.Close()
}

// Function to get current screen session processes and it's windows
func current_snapshot() string {
	var screen_pids string = strings.TrimSpace(exec_cmd(SCRIPT_ARG))
	return screen_pids
}

func main() {

	// Notify user that program has started
	notify_user("started")

	// Capture any signals before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		notify_user("exited")
		log_entry("Interrupted")
		os.Exit(1)
	}()

	for {

		// get current snapshot of screen sessions and it's child processes
		past_sessions := current_snapshot()
		write_to_file(PAST_TMP_FILE, past_sessions)

		// wait for some time
		time.Sleep(time.Second * 5)

		// Now, get another snapshot of sessions
		curr_sessions := current_snapshot()
		write_to_file(CURR_TMP_FILE, curr_sessions)

		if curr_sessions != past_sessions {

			// get difference
			var diff string = strings.TrimSpace(exec_cmd(""))

			// Get the hostname
			hostname, err := os.Hostname()
			check(err)

			var subject = "Screen Monitor Discrepancy Report for Host: " + hostname
			var content = "Hello,<br><br>"
			content += "There has been discrepancies regarding screen sessions on this host: <b>" + hostname + "</b>. <br><br>"
			content += "The main discrepancy is either some process is created or crashed behind screen sessions. The changes are shown below:<br><br>"
			content += diff + "<br><br>"
			content += "For further details, please review the attached log files.<br><br><br>"
			content += "past_sessions.log shows the snapshot of screen sessions running 20 secs before the discrepancy occurred.<br>"
			content += "curr_sessions.log shows the snapshot of screen sessions running after the discrepancy occurred.<br><br>"
			content += "It is advisable to investigate this matter further to minimize any outage time of service(s).<br><br>"
			content += "Regards,<br>"
			content += "SAVI Testbed Team"

			send_email(subject, content, true)

			log_entry("Sent Email regarding discrepancy!")

		} else {

			log_entry("No changes to screen sessions!")
		}
	}

	notify_user("exited")
	log_entry("Exited")
	os.Exit(1)
}
