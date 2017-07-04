#!/bin/bash

#Find child processes of screen example: Ref > https://superuser.com/questions/363169/ps-how-can-i-recursively-get-all-child-process-for-a-given-pid

function get_procs {
    proc_id_list=( $(ps fx | grep 'SCREEN' | grep '?' | awk '{print $1}' | sort -rn) )  # Get the screen sessions process IDs
    pid_list_len=${#proc_id_list[@]}                                                    # Get the length of the list

    if [ -n "$proc_id_list" ]; then                                                     # If not empty, do the following
        for pid in "${proc_id_list[@]}"                                                 # Go through each screen process ID
        do  
            list_offspring $pid                                                         # Get the info of it's children
        done
    fi
}

function list_offspring {
  tp=`pgrep -P $1`                      # Get the children of the parent screen process
  for i in $tp; do                      # Go through all the children (/bin/bash)
    if [ ! -z $i ]; then                # If not empty, do the following
      tp_2=`pgrep -P $i`                    # Get the first child of that /bin/bash process
      if [ ! -z "$tp_2" ]; then                 # If not empty, do the following
        ps -p $tp_2 --no-header -o cmd              # Get the command of that process 
      fi;
    fi;
  done
}

if [[ "$1" = "get_procs" ]]; then                   # If argument mentions "get_procs", then do the following
    get_procs                                       # Execute get_procs function
else                                                # else, do the following
    diff=`diff $1 $2 | awk '{if (NR > 1) print }'`      # Get the difference between 2 files. Also, remove the first line from diff output 
    echo $diff
fi

