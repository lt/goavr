#!/bin/bash 

for i in `echo *`
do 
  echo "Cleaning $i"
  if [ -d $i ]
    then cd $i
    make clean 
    cd ../
   fi
done
