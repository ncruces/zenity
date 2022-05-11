#!/bin/sh

ENTRY=`zenity --password --username`

case $? in
  0)
    echo "User Name: `echo $ENTRY | cut -d'|' -f1`"
    echo "Password : `echo $ENTRY | cut -d'|' -f2`"
    ;;
  1)
    echo "Stop login.";;
  -1)
    echo "An unexpected error has occurred.";;
esac