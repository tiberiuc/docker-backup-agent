#!/bin/sh

export BACKUP_AGENT_PORT=9191
export BACKUP_AGENT_BACKUP_SCRIPT=./backup_script.sh
export BACKUP_AGENT_RESTORE_SCRIPT=./restore_script.sh

mkdir -p ../build

cd .. && go build -o ./build/backup_agent && cd test
pwd

../build/backup_agent -port=9191 -backup=./backup_script.sh -restore=./restore_script.sh &

PID=$!

sleep 2

echo Server started with PID=$PID

rm -f /tmp/downloaded_backup.tgz

curl --silent http://localhost:9191/backup  -o /tmp/downloaded_backup.tgz

echo "Testing backup files"
cmp --silent /tmp/created_backup.tgz /tmp/downloaded_backup.tgz || echo "Created backup files are different"

curl --silent  http://localhost:9191/restore --data-binary @/tmp/downloaded_backup.tgz > /dev/null

echo "Testing resulted restore"
diff -r -q ./ /tmp/backup_restored

rm -rf /tmp/created_backup.tgz /tmp/downloaded_backup.tgz /tmp/backup_restored

kill -9 $PID
