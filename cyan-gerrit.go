package main

import(
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main(){
	type Gerrit_s struct { //struct for unmarshalling JSON
		Project string
		Subject string
		Created string
		Updated string
		Sortkey string
		Number int
	}

	lastrunData, err := ioutil.ReadFile("lastrun") //file containing time program was last run
	var lastrun time.Time //time that program was last run
	if os.IsNotExist(err) {
		lastrun = time.Now().AddDate(0, 0, -1) //if lastrun file doesn't exist, set to 1 day ago
	} else if err != nil {
		fmt.Println("Error reading lastrun: ", err)
		os.Exit(1)
	} else {
		lastrun, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(lastrunData)) //parse time last run from file
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
	body = bytes.TrimPrefix(body, []byte(")]}'\n")) //get rid of odd delimiters at beginning of response

	var changes []Gerrit_s
	err = json.Unmarshal(body, &changes)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data: ", err)
		os.Exit(5)
	}

	fmt.Println("Project\t\tSubject\t\tTime") //print changes
	for _, change := range changes {
		changeTime, err := time.Parse("2006-01-02 15:04:05.000000000", change.Updated) //parse last updated time for change
		if err == nil && changeTime.After(lastrun) { //if there was no error and updated time is after last run time
			if (strings.HasPrefix(change.Project, "CyanogenMod/android_device") && //skip any device commits that aren't an oppo device
				!strings.HasPrefix(change.Project, "CyanogenMod/android_device_oppo")) ||
				(strings.HasPrefix(change.Project, "CyanogenMod/android_kernel") && //skip any kernel commits that aren't oneplus
				!strings.HasPrefix(change.Project, "CyanogenMod/android_kernel_oneplus")) {
					continue
				}
			fmt.Printf("%s\t\t%s\t\t%s\n", change.Project, change.Subject, changeTime.Format("01-02 15:04")) //print change project, subject, and updated time
		}
	}

	ioutil.WriteFile("lastrun", []byte(time.Now().String()), 0644) //update lastrun file
}
