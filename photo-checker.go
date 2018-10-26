package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"net/http"
)

func main() {
	fmt.Println("hello")
	for i, arg := range os.Args {
		fmt.Printf("arg[%d]='%s'\n", i, arg)
	}

	baseDir := os.Args[1]

	out, err := exec.Command("ls", "-1", baseDir).Output()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	currentDate := t.Format("2006-01-02")
	auditOutputDir := fmt.Sprintf("%s/Hashdeep/Hashdeep-Audit-%s", baseDir, currentDate)
	sh.Command("mkdir", auditOutputDir).Run()

	directories := strings.Split(string(out), "\n")
	for _, dir := range directories {
		// sh.Command("echo", dir).Run()
		logFile := fmt.Sprintf("%s/Hashdeep/Hashdeep-Logs/%s.txt", baseDir, dir)
		fullDirPath := fmt.Sprintf("%s/%s", baseDir, dir)
		auditOutputFile := fmt.Sprintf("%s/%s-Audit-%s.txt", auditOutputDir, dir, currentDate)
	    //fmt.Println("audit output file:")
	    //fmt.Println(auditOutputFile);		
		sh.Command("hashdeep", "-vv", "-a", "-k", logFile, "-r", fullDirPath).WriteStdout(auditOutputFile)
	}
	notify("Checking complete.")

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
