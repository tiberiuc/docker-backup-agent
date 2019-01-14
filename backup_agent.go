package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	port := "9191"
	envPort := os.Getenv("BACKUP_AGENT_PORT")
	envBackupCommand := os.Getenv("BACKUP_AGENT_BACKUP_SCRIPT")
	envRestoreCommand := os.Getenv("BACKUP_AGENT_RESTORE_SCRIPT")
	if envPort != "" {
		port = envPort
	}

	intPort, errPort := strconv.Atoi(port)

	if errPort != nil {
		panic("Port must be a number")
	}

	portPtr := flag.Int("port", intPort, "port to start server on, default 9191. Also can be specified by BACKUP_AGENT_PORT environment variable")
	backupCmdPtr := flag.String("backup", envBackupCommand, "backup command to run. Also can be specified by BACKUP_AGENT_BACKUP_SCRIPT environment variable")
	restoreCmdPtr := flag.String("restore", envRestoreCommand, "restore command to run. Also can be specified by BACKUP_AGENT_RESTORE_SCRIPT environment variable")

	flag.Parse()

	port = strconv.Itoa(*portPtr)

	fmt.Println("Backup command:" + *backupCmdPtr)
	fmt.Println("Restore command:" + *restoreCmdPtr)

	http.HandleFunc("/backup", func(w http.ResponseWriter, r *http.Request) {
		var stdout, stderr bytes.Buffer

		cmd := exec.Command("/bin/sh", "-c", *backupCmdPtr)

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fileName := strings.TrimSpace(outStr)
		_, fileErr := os.Stat(fileName)
		if err != nil || fileErr != nil {

			outError := outStr + "\n" + errStr
			if err != nil {
				outError += err.Error()
			}

			if fileErr != nil {
				outError += fileErr.Error()
			}

			http.Error(w, outError, http.StatusInternalServerError)
		} else {
			log.Println("Sending backup file: ", fileName)
			http.ServeFile(w, r, fileName)
		}
	})

	http.HandleFunc("/restore", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, *restoreCmdPtr)
		tmpfile, err := ioutil.TempFile("", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		defer os.Remove(tmpfile.Name())

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		if _, err := tmpfile.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		if err := tmpfile.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		cmd := exec.Command("/bin/sh", "-c", *restoreCmdPtr+" "+tmpfile.Name())

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {

			outError := string(stdoutStderr) + err.Error()

			http.Error(w, outError, http.StatusInternalServerError)
		} else {
			log.Println("Backup restored")
			fmt.Fprintf(w, string(stdoutStderr))
		}
	})

	fmt.Println("Starting backup server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
