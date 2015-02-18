package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	now := time.Now()
	type Gerrit_s struct { //struct for unmarshalling JSON
		Project string
		Subject string
		Created string
		Updated string
		Sortkey string `json:"_sortkey"`
		Number  int    `json:"_number"`
	}

	lastrunData, err := ioutil.ReadFile("lastrun") //file containing time program was last run
	var lastrun time.Time                          //time that program was last run
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

	lastSortKey := ""

	fmt.Printf("%-30s\t%-100s\t%-10s\t%s\n", "Project", "Subject", "Time", "URL") //print changes header
FetchLoop:
	for true { //run until find a change with updated time before lastrun time (potential for infinite loop if time.Parse constantly fails
		request := "http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0&n=50"
		if lastSortKey != "" { //if lastSortKey isn't blank, continue query starting after last seen change
			request += "&N=" + lastSortKey
		}

		resp, err := http.Get(request)
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

		_, offset_seconds := now.Zone() //get offset of current timezone to allow printing local time
		offset, err := time.ParseDuration(strconv.FormatInt(int64(offset_seconds), 10) + "s")
		if err != nil {
			fmt.Println("Error parsing tz offset: ", err)
			os.Exit(6)
		}
		for _, change := range changes {
			changeTime, err := time.Parse("2006-01-02 15:04:05.000000000", change.Updated) //parse last updated time for change
			if err == nil {
				if changeTime.Before(lastrun) { //if updated time is before lastrun time, break loop
					break FetchLoop
				}
				if (strings.HasPrefix(change.Project, "CyanogenMod/android_device") && //skip any device commits that aren't an oppo device
					!strings.HasPrefix(change.Project, "CyanogenMod/android_device_oppo")) ||
					(strings.HasPrefix(change.Project, "CyanogenMod/android_kernel") && //skip any kernel commits that aren't oneplus
						!strings.HasPrefix(change.Project, "CyanogenMod/android_kernel_oneplus")) {
					continue
				}
				fmt.Printf("%-30.30s\t%-100.100s\t%-11.11s\thttp://review.cyanogenmod.org/#/c/%d/\n", //print change project, subject, updated time, and URL
					strings.TrimPrefix(change.Project, "CyanogenMod/android_"),
					change.Subject,
					changeTime.Add(offset).Format("01-02 15:04"), //print time in local zone
					change.Number)
			}
		}
		lastSortKey = changes[len(changes)-1].Sortkey //update last Sortkey in slice
	}

	ioutil.WriteFile("lastrun", []byte(now.String()), 0644) //update lastrun file
}
