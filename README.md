
Extremly customizable and easy to use backup agent server to be used inside docker. Used to backup data from a container when you don't have access to files inside that container

The backup agent start a server and expose two methods `backup` and `restore`. You can write your own backup restore scripts in any language. The backup script should return on stdin the filename of the backup and restore script accept as first parameter the name of the backup file

In case of an error both methods return an InternalServerError and the output of the script.

Parameters can be sent on command line (`-port`, `-backup`, `-restore`) or with environment variables (`BACKUP_AGENT_PORT`, `BACKUP_AGENT_BACKUP_SCRIPT`, `BACKUP_AGENT_RESTORE_SCRIPT`). Command line parameters have priority.


# Usage Example

## Backup Script

```
#!/bin/bash

BACKUP_FILENAME=/tmp/backup.tgz

rm -f ${BACKUP_FILENAME}

tar -cvf ${BACKUP_FILENAME} /data

echo ${BACKUP_FILENAME}
```

## Restore Script

```
#!/bin/bash

BACKUP_FILENAME=$2

tar -xvf ${BACKUP_FILENAME} ./

```

## Starting

```
./backup_agent -port=9191 -backup=./backup_script.sh -restore=./restore_script.sh
```

## Getting Backup
```
curl http://localhost:9191/backup  -o backup.tgz
```

### Restore Backup
```
curl -vX POST  http://localhost:9191/restore --data-binary @backup.tgz
```

## Real Examples

you can check this projects for real life usage

* [tiberiuc/docker-bitcoind](https://github.com/tiberiuc/docker-bitcoind)
* [tiberiuc/docker-ethereum](https://github.com/tiberiuc/docker-ethereum)
* [tiberiuc/docker-mongodb](https://github.com/tiberiuc/docker-mongodb)
