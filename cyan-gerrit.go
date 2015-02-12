package main

import(
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"strings"
)

func main(){
	type gerrit_s struct {
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

	resp, err := http.Get("http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0&n=2")
	if err != nil {
		os.Exit(3)
	}
	defer resp.Body.Close()
	body_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(4)
	}
	body := string(body_data)
	body = strings.TrimPrefix(body, ")]}'\n")
	fmt.Println(body)

	now := time.Now().String()
	ioutil.WriteFile("lastrun", []byte(now), 0644)
}
