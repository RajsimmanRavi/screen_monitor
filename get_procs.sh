#!/bin/bash

function get_procs {
    proc_id_list=( $(ps fx | grep 'SCREEN' | grep '?' | awk '{print $1}' | sort -rn) )

    pid_list_len=${#proc_id_list[@]}

    if [ -n "$proc_id_list" ]; then
        for pid in "${proc_id_list[@]}"
        do
            ps --forest $(ps -e --no-header -o pid,ppid|awk -vp="${pid}" 'function r(s){print s;s=a[s];while(s){sub(",","",s);t=s;sub(",.*","",t);sub("[0-9]+","",s);r(t)}}{a[$2]=a[$2]","$1}END{r(p)}')| grep -P '\d'
        done
    fi
}

if [[ "$1" = "get_procs" ]]; then 
    get_procs
else
    #diff /tmp/past_sessions.log /tmp/curr_sessions.log
    cat <(grep -vFf /tmp/curr_sessions.log /tmp/past_sessions.log) <(grep -vFf /tmp/past_sessions.log /tmp/curr_sessions.log)
fi



