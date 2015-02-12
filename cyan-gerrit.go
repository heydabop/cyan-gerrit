package main

import(
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main(){
	resp, err := http.Get("http://review.cyanogenmod.org/changes/?q=status:merged+branch:cm-12.0&n=2")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(2)
	}
	fmt.Println(string(body))
}
