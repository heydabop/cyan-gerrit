package main

import(
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
)

func main(){
	lastrunData, err := ioutil.ReadFile("lastrun")
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
	}

	request := fmt.Sprintf("http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0+age:%dm", uint64(math.Ceil(time.Since(lastrun).Minutes())))
	fmt.Println(request)
	resp, err := http.Get(request)
	if err != nil {
		os.Exit(3)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(4)
	}
	fmt.Println(string(body))

	now := time.Now().String()
	ioutil.WriteFile("lastrun", []byte(now), 0644)
}
