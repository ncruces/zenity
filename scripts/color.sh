#!/bin/sh

COLOR=`zenity --color-selection --show-palette`

case $? in
  0)
    echo "You selected $COLOR.";;
  1)
    echo "No color selected.";;
  -1)
    echo "An unexpected error has occurred.";;
esac