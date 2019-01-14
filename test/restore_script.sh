#!/bin/bash

BACKUP_FILENAME=$1

rm -rf /tmp/backup_restored
mkdir -p /tmp/backup_restored
tar -xzf ${BACKUP_FILENAME} -C /tmp/backup_restored/
echo "Backup restored"

