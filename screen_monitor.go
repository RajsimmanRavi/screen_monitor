package main

import ( 
  "fmt"
  "os/exec"
  "os/signal"
  "os"
  "time"
  "strings"
  "syscall"
  "log"
  "gopkg.in/gomail.v2"
)


//Execute bash command in go: Ref > https://stackoverflow.com/questions/20437336/how-to-execute-system-command-in-golang-with-unknown-arguments?rq=1
//Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid
//Catch Signals: Ref > https://stackoverflow.com/questions/18106749/golang-catch-signals

// Function to check for errors. Copied from Golang example
func check(e error){
  if e!= nil{
    notify_exit()
    panic(e)
  }
}

// Function to log errors, warnings etc. to file
func log_entry(info string){

  f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(err)
  
  defer f.Close()

  log.SetOutput(f)
  log.Println()

}

// Function to execute bash command 
func exec_cmd(cmd string) string {
  //fmt.Println(cmd)
  out, err := exec.Command("bash","-c",cmd).Output()
  check(err)
  
  return string(out)
}

// Function to email user regarding program exit
func notify_exit(){

  hostname,err := os.Hostname()
  check(err)

  var subject = "Screen Monitor Program exited"
  var content = "Hello,<br><br>"
  content += "This is to notify you that a system call signal has been caught before exiting on this host: <b>"+hostname+"</b>.<br><br>"
  content += "Please restart the program if it has exited.<br><br>"
  content += "Regards,<br>"
  content += "SAVI Testbed Team"

  send_email(subject, content)

}

// Function to send email
func send_email(subject string, content string){

  m := gomail.NewMessage()
  m.SetHeader("From", "xxxx")
  m.SetHeader("To", "xxxx")
  m.SetHeader("Subject", subject)
  content = strings.Replace(content, "\n", "<br>",-1)
  //fmt.Println(content)
  m.SetBody("text/html", content)

  d := gomail.NewDialer("smtp.gmail.com", 587, "xxxx", "xxxx")

  var err = d.DialAndSend(m)
  check(err)

}

// Function to get current screen session processes and it's windows 
func current_snapshot() string {

  var screen_pid_cmd string = "ps fx | grep 'SCREEN' | grep '?' | grep -v 'bash' | cut -d \" \" -f 1 | sort -rn"
  var screen_pids []string = strings.Split(strings.TrimSpace(exec_cmd(screen_pid_cmd)), "\n")

  // array to hold the screen processes output
  var screen_trees string 

  // Traverse through each screen session pid and get it's child processes
  for _,pid := range screen_pids {

    var get_child_procs string = "ps --forest $(ps -e --no-header -o pid,ppid|awk -vp="+pid+" 'function r(s){print s;s=a[s];while(s){sub(\",\",\"\",s);t=s;sub(\",.*\",\"\",t);sub(\"[0-9]+\",\"\",s);r(t)}}{a[$2]=a[$2]\",\"$1}END{r(p)}')| grep -P '\\d';"

    var screen_procs string = strings.TrimSpace(exec_cmd(get_child_procs))

    // Append it to return string
    screen_trees = screen_trees+screen_procs+"\n\n" 

  }
  
  return screen_trees
 
}  


func main(){
  
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
    notify_exit()
    os.Exit(1)
  }()
  
  for {

    // get current snapshot of screen sessions and it's child processes
    past_sessions := current_snapshot()
    
    // wait for some time 
    time.Sleep(time.Second * 5)
    
    // Now, get another snapshot of sessions 
    curr_sessions := current_snapshot()
    
    if curr_sessions != past_sessions {

      // Get the hostname
      hostname,err := os.Hostname()
      check(err) 
      
      var subject = "Screen Monitor Discrepancy Report"
      var content = "Hello,<br><br>"
      content += "There has been discrepancies regarding screen sessions on this host: <b>"+hostname+"</b>. The details are shown below.<br><br>"       
      content += "<b>OLD SCREEN SESSIONS</b><br>"
      content += past_sessions
      content += "<b>NEW SCREEN SESSIONS</b><br>"
      content += curr_sessions
      content += "It is advisable to investigate this matter further to minimize any outage time of service(s).<br><br>"
      content += "Regards,<br>"
      content += "SAVI Testbed Team"
      
      send_email(subject, content)
      
    }

  }

  os.Exit(1)
  
}

