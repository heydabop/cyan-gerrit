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
	/*lastrunData, err := ioutil.ReadFile("lastrun")
	var lastrun time.Time
	if os.IsNotExist(err) {
		lastrun = time.Now().AddDate(0, 0, -1)
	} else if err != nil {
		os.Exit(1)
	} else {
		lastrun, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(lastrunData))
		if err != nil {
			os.Exit(2)
		}
	}*/

	resp, err := http.Get("http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0&n=3")
	if err != nil {
		os.Exit(3)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(4)
	}
	body = bytes.TrimPrefix(body, []byte(")]}'\n"))
	fmt.Printf("%s\n", body)
	var changes []Gerrit_s
	err = json.Unmarshal(body, &changes)
	if err != nil {
		os.Exit(5)
	}
	for _, change := range changes {
		fmt.Printf("%+v\n", change)
	}

	now := time.Now().String()
	ioutil.WriteFile("lastrun", []byte(now), 0644)
}
