package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Emoji struct {
	Shortcode       string
	StaticUrl       string
	Url             string
	VisibleInPicker bool
}

func main() {

	flag.Parse()
	var resp *http.Response
	var err error
	arg := flag.Args()[0]
	arg2 := arg

	exe, err6 := os.Executable()

	if hasError(err6) {
		return
	}

	exe = strings.TrimSuffix(exe, "gomoji") // hardcoded name, bad

	if !strings.HasPrefix(arg, "https://") &&
		!strings.HasPrefix(arg, "http://") {
		arg = "https://" + arg
	} else if strings.HasPrefix(arg, "http://") {
		arg = strings.Replace(arg, "http://", "https://", 1)
	}

	resp, err = http.Get(arg + "/api/v1/custom_emojis")

	if hasError(err) {
		return
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)

	if hasError(err2) {
		return
	}

	var emoji *[]Emoji
	err3 := json.Unmarshal(body, &emoji)

	if hasError(err3) {
		return
	}

	err7 := os.MkdirAll(exe+"/"+arg2, os.ModePerm)

	if hasError(err7) {
		return
	}

	for i := 0; i < len(*emoji); i++ {

		//	urlReader := strings.NewReader((*emoji)[i].Url)

		file, err4 := os.Create(exe + arg2 + "/" + (*emoji)[i].Shortcode + ".png")

		if hasError(err4) {
			os.Exit(1) // crash program
		}
		defer file.Close()

		img, err8 := http.Get((*emoji)[i].Url)

		if hasError(err8) {
			os.Exit(1) //crash program
		}

		defer img.Body.Close()

		_, err5 := io.Copy(file, img.Body)

		if hasError(err5) {
			os.Exit(1) // crash program
		}

		fmt.Println(exe + arg2 + "/" + (*emoji)[i].Shortcode + ".png created")
	}

	fmt.Println("Finished!")

}

func hasError(err error) bool {
	if err != nil {
		fmt.Println(err.Error)
		return true
	}
	return false
}
