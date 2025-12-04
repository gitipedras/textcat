#!/bin/bash

############
#          #
#   BASH   #
#          #
############

echo "###############"
echo "This script creates a channel in a pre-existing textcat database"
echo "I made this so u can create channels (until addons are added)"
echo "REQUIRES: sqlite3, bash/sh/zsh"
echo "###############"

read -p "Enter database file: " dbFile
read -p "Enter channel name: " chName

sqlite3 "$dbFile" << EOF
	INSERT INTO channels (name, extraData)
	VALUES('$chName', NULL)	
EOF
