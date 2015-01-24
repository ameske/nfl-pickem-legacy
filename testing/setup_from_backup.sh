#!/bin/bash

#***************************
# setup_from_backup.sh
#
# Sets up a go_nfl databse from a Postgres .bak file
#
# Usage:
#   setup_from_backup.sh <BACKUP LOCATION>
#
# Author: Kyle Ames
#***************************

BACKUP=$1
PSQL='/Applications/Postgres.app/Contents/Versions/9.3/bin/psql'

$PSQL < clear_db.sql
$PSQL -c "CREATE ROLE nfl WITH LOGIN"
$PSQL -c "CREATE DATABASE nfl_app"
$PSQL nfl_app < $BACKUP
