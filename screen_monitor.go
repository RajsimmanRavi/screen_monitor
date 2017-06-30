package main

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	//"time"
    "io/ioutil"
)

const CURR_TMP_FILE = "/tmp/curr_sessions.log"
const LOG_FILE = "/tmp/screen_monitor_logs"
const IDEAL_SESSIONS_FILE = "/home/xxx/go/src/screen_monitor/ideal_screen_sessions.log"
const SCRIPT = "/home/xxx/go/src/screen_monitor/get_procs.sh"
const SCRIPT_ARG = "get_procs"
const LOG_FILE_SIZE = 2000000 // 2MB

const FROM_EMAIL = "xxx@xx"
const TO_EMAIL = "xx@xx"
const PWD = "xx"

// Author: Rajsimman Ravi
// Any questions, contact: rajsimmanr@savinetwork.ca

//Execute bash command in go: Ref > https://stackoverflow.com/questions/12891294/google-golang-exec-exit-status-2-and-1
//Catch Signals: Ref > https://stackoverflow.com/questions/18106749/golang-catch-signals

// Function to check for errors. Copied from Golang example
func check(e error) {
	if e != nil {
		notify_user("exited")
		log_entry(e.Error())
		panic(e)
	}
}

// Function to check if a file exists
func check_file_exists(file_name string) bool {

    if _, err := os.Stat(file_name); os.IsNotExist(err) {
    // file does not exist
        return false
    }
    
    return true
}

//Function to check sizes of log file. 
//If it's more than certain length, refresh or delete the file
func check_log_size(){

    file, err := os.Open(LOG_FILE) 
    check(err)
    
    fi, err := file.Stat()
    check(err)

    f_size := fi.Size() 
    if f_size > LOG_FILE_SIZE {
    
        err := os.Remove(LOG_FILE)
        check(err)
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
func exec_cmd(comd string ) string {

	cmd := exec.Command("/bin/bash", SCRIPT, comd)
	var out bytes.Buffer
	var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
	err := cmd.Run()
    if err != nil {
        log_entry(fmt.Sprint(err) + ": " + stderr.String())
        check(err)
    }

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
		m.Attach(IDEAL_SESSIONS_FILE)
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
 
    // Get current screen sessions  
    curr_sessions := current_snapshot()
    write_to_file(CURR_TMP_FILE, curr_sessions)

    var past_sessions string

    if check_file_exists(IDEAL_SESSIONS_FILE){

        // Get the ideal screen sessions log. This is to compare what has changed since the screen sessions started running
        dat, err := ioutil.ReadFile(IDEAL_SESSIONS_FILE)
        check(err)
        
        past_sessions = strings.TrimSpace(string(dat))
    } else {
    
        // create file and write the current sessions
        write_to_file(IDEAL_SESSIONS_FILE, curr_sessions)
    
        past_sessions = curr_sessions
    }     


    if curr_sessions != past_sessions {

        // get difference
        var diff string = strings.TrimSpace(exec_cmd(IDEAL_SESSIONS_FILE+" "+CURR_TMP_FILE))
        
        // Get the hostname
        hostname, err := os.Hostname()
        check(err)

        var subject = "Screen Monitor Discrepancy Report for Host: " + hostname
        var content = "Hello,<br><br>"
        content += "There has been discrepancies regarding screen sessions on this host: <b>" + hostname + "</b>. <br><br>"
        content += "The main discrepancy is either some process is created or crashed behind screen sessions. The changes are shown below:<br><br>"
        content += diff + "<br><br>"
        content += "For further details, please review the attached log files.<br><br><br>"
        content += "Ideal_sessions.log shows the screen sessions that should be running.<br>"
        content += "curr_sessions.log shows the snapshot of current screen sessions. <br><br>"
        content += "It is advisable to investigate this matter further to minimize any outage time of service(s).<br><br>"
        content += "Regards,<br>"
        content += "SAVI Testbed Team"

        send_email(subject, content, true)

        log_entry("Sent Email regarding discrepancy!")

    } else {

        log_entry("No changes to screen sessions!")
    }

    // check and refresh log file
    check_log_size()
}
