package main

import(
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main(){
	type Gerrit_s struct {
		Project string
		Subject string
		Created string
		Updated string
		Sortkey string
		Number int
	}
	lastrunData, err := ioutil.ReadFile("lastrun")
	var lastrun time.Time
	if os.IsNotExist(err) {
		lastrun = time.Now().AddDate(0, 0, -1)
	} else if err != nil {
		fmt.Println("Error reading lastrun: ", err)
		os.Exit(1)
	} else {
		lastrun, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(lastrunData))
		if err != nil {
			fmt.Println("Error parsing lastrun: ", err)
			os.Exit(2)
		}
	}

	resp, err := http.Get("http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0&n=50")
	if err != nil {
		fmt.Println("Error getting JSON data: ", err)
		os.Exit(3)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading JSON data: ", err)
		os.Exit(4)
	}
	body = bytes.TrimPrefix(body, []byte(")]}'\n"))
	var changes []Gerrit_s
	err = json.Unmarshal(body, &changes)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data: ", err)
		os.Exit(5)
	}
	fmt.Println("Project\t\tSubject\t\tTime")
	for _, change := range changes {
		changeTime, err := time.Parse("2006-01-02 15:04:05.000000000", change.Updated)
		if err == nil && changeTime.After(lastrun) {
			fmt.Printf("%s\t\t%s\t\t%s\n", change.Project, change.Subject, changeTime.Format("01-02 15:04"))
		}
	}

	now := time.Now().String()
	ioutil.WriteFile("lastrun", []byte(now), 0644)
}
