# Screen Monitor
A program to monitor screen sessions and it's window panes (running on Linux Ubuntu). The program will take a snapshot of current screen sessions (running on the system) and compares it with an "ideal" screen sessions log file. This file should contain the status of screen sessions when it's running properly. If no file, this program creates the file on the specified path. 

## Dependencies

This program does require the gomail package (https://github.com/go-gomail/gomail). Hence you can download them by using the following command:

    go get gopkg.in/gomail.v2
    
## Installation

You can create a cron job to run every minute:

    * * * * * /usr/local/go/bin/go run /path/to/screen_monitor.go

## Contact

If there is any questions, contact me: rajsimmanr@savinetwork.ca




