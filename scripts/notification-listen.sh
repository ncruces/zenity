#!/bin/sh

cat <<EOH| zenity --notification --listen
message: this is the message text
EOH