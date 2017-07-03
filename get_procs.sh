#!/bin/bash

#Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid

function get_procs {
    proc_id_list=( $(ps fx | grep 'SCREEN' | grep '?' | awk '{print $1}' | sort -rn) )

    pid_list_len=${#proc_id_list[@]}

    if [ -n "$proc_id_list" ]; then
        for pid in "${proc_id_list[@]}"
        do
            list_offspring $pid
        done
    fi
}

function list_offspring {
  tp=`pgrep -P $1`
  for i in $tp; do
    if [ ! -z $i ]; then
      tp_2=`pgrep -P $i`
      if [ ! -z "$tp_2" ]; then
        ps -p $tp_2 --no-header -o cmd
      fi;
    fi;
  done
}

if [[ "$1" = "get_procs" ]]; then 
    get_procs
else
    cat <(grep -vFf $1 $2) <(grep -vFf $2 $1)
fi

