#!/bin/bash

BACKUP_FILENAME=/tmp/created_backup.tgz

rm -f ${BACKUP_FILENAME}

tar -czf ${BACKUP_FILENAME} ./

echo ${BACKUP_FILENAME}
