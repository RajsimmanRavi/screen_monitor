#!/bin/bash

<<<<<<< HEAD
#Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid

=======
>>>>>>> 70c0dd0125e2e56886ab6e35f027c5d9f1af304c
function get_procs {
    proc_id_list=( $(ps fx | grep 'SCREEN' | grep '?' | awk '{print $1}' | sort -rn) )

    pid_list_len=${#proc_id_list[@]}

    if [ -n "$proc_id_list" ]; then
        for pid in "${proc_id_list[@]}"
        do
<<<<<<< HEAD
            # Ignoring sessions that have "timeout" because they change all the time (created/exited)
            ps --no-header -o pid,cmd --forest $(ps -e --no-header -o pid,ppid|awk -vp="${pid}" 'function r(s){print s;s=a[s];while(s){sub(",","",s);t=s;sub(",.*","",t);sub("[0-9]+","",s);r(t)}}{a[$2]=a[$2]","$1}END{r(p)}') | grep -v "timeout"
=======
            ps --forest $(ps -e --no-header -o pid,ppid|awk -vp="${pid}" 'function r(s){print s;s=a[s];while(s){sub(",","",s);t=s;sub(",.*","",t);sub("[0-9]+","",s);r(t)}}{a[$2]=a[$2]","$1}END{r(p)}')| grep -P '\d'
>>>>>>> 70c0dd0125e2e56886ab6e35f027c5d9f1af304c
        done
    fi
}

if [[ "$1" = "get_procs" ]]; then 
    get_procs
else
<<<<<<< HEAD
    cat <(grep -vFf $1 $2) <(grep -vFf $2 $1)
fi

=======
    #diff /tmp/past_sessions.log /tmp/curr_sessions.log
    cat <(grep -vFf /tmp/curr_sessions.log /tmp/past_sessions.log) <(grep -vFf /tmp/past_sessions.log /tmp/curr_sessions.log)
fi



>>>>>>> 70c0dd0125e2e56886ab6e35f027c5d9f1af304c
