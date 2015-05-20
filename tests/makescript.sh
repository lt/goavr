#!/bin/bash 

for i in `echo *`
do 
  if [ -d $i ]
    then cd $i
    make
    cd ../
   fi
done
