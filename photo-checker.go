package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout

	baseDirPtr := flag.String("base-dir", "", "base directory for checking photos")
	dryRunPtr := flag.Bool("dry-run", false, "send this flag to do a dry run without checking")

	flag.Parse()

	log.WithFields(logrus.Fields{
		"base-dir": *baseDirPtr,
		"dry-run":  *dryRunPtr,
	}).Info("Command line inputs: ")

	out, err := exec.Command("ls", "-1", *baseDirPtr).Output()
	if err != nil {
		log.Error("Error listing directory contents: ", err)
		log.Fatal(err)
	}

	t := time.Now()
	currentDate := t.Format("2006-01-02")
	auditOutputDir := fmt.Sprintf("%s/Hashdeep/Hashdeep-Audits/Hashdeep-Audit-%s", *baseDirPtr, currentDate)

	log.WithFields(logrus.Fields{
		"auditOutputDir": auditOutputDir,
	}).Info("Creating Audit output directory: ")

	if !*dryRunPtr {
		if !sh.Test("dir", auditOutputDir) {
			sh.Command("mkdir", auditOutputDir).Run()
		} else {
			log.Fatal("Directory already exists: ", auditOutputDir)
		}
	}

	directories := strings.Split(string(out), "\n")
	for _, dir := range directories {
		if (len(dir) != 0) && (!strings.Contains(dir, "Hashdeep")) {
			// sh.Command("echo", dir).Run()
			log.Info("Current directory to process: ", dir)

			logFile := fmt.Sprintf("%s/Hashdeep/Hashdeep-Logs/%s.txt", *baseDirPtr, dir)
			fullDirPath := fmt.Sprintf("%s/%s", *baseDirPtr, dir)
			auditOutputFile := fmt.Sprintf("%s/%s-Audit-%s.txt", auditOutputDir, dir, currentDate)

			log.WithFields(logrus.Fields{
				"auditOutputFile": auditOutputFile,
			}).Info("Audit output file: ")
			log.Info("hashdeep -vv -a -k", logFile, " -r ", fullDirPath)

			if !*dryRunPtr {
				sh.Command("hashdeep", "-vv", "-a", "-k", logFile, "-r", fullDirPath).WriteStdout(auditOutputFile)
			}
		}
	}
	if !*dryRunPtr {
		notify("Checking complete.")
	}
}

func notify(message string) {
	bodyWithMessage := fmt.Sprintf("token=..&user=..&message=%s", message)
	body := strings.NewReader(bodyWithMessage)

	req, err := http.NewRequest("POST", "https://api.pushover.net/1/messages.json", body)
	if err != nil {
		// don't really care if the notification is an error
	}
	req.SetBasicAuth("API_KEY", "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// don't really care if the notification is an error
	}
	defer resp.Body.Close()
}
