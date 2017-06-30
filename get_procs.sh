#!/bin/bash

#Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid

function get_procs {
    proc_id_list=( $(ps fx | grep 'SCREEN' | grep '?' | awk '{print $1}' | sort -rn) )

    pid_list_len=${#proc_id_list[@]}

    if [ -n "$proc_id_list" ]; then
        for pid in "${proc_id_list[@]}"
        do
            # Ignoring sessions that have "timeout" because they change all the time (created/exited)
            ps --no-header -o pid,cmd --forest $(ps -e --no-header -o pid,ppid|awk -vp="${pid}" 'function r(s){print s;s=a[s];while(s){sub(",","",s);t=s;sub(",.*","",t);sub("[0-9]+","",s);r(t)}}{a[$2]=a[$2]","$1}END{r(p)}') | grep -v "timeout"
        done
    fi
}

if [[ "$1" = "get_procs" ]]; then 
    get_procs
else
    cat <(grep -vFf $1 $2) <(grep -vFf $2 $1)
fi

