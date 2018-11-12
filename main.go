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
	Shortcode       string // emoji name
	StaticUrl       string // doesn't seem to be used?
	Url             string // emoji url
	VisibleInPicker bool   // always returns false?
}

func main() {

	flag.Parse()            // parse user input
	var resp *http.Response // API call to mastoAPI emoji
	var err2 error          // error for resp
	arg := flag.Args()[0]   // should be a url
	arg2 := arg

	exe, err1 := os.Executable()

	if hasError(err1) { // fail if can't find executable path
		return
	}

	exe = strings.TrimSuffix(exe, "gomoji") // hardcoded name, bad

	if !strings.HasPrefix(arg, "https://") &&
		!strings.HasPrefix(arg, "http://") {
		// add https:// if arg doesn't have it
		arg = "https://" + arg
	} else if strings.HasPrefix(arg, "http://") {
		// change http:// to https:// in arg
		arg = strings.Replace(arg, "http://", "https://", 1)
	}

	resp, err2 = http.Get(arg + "/api/v1/custom_emojis")

	if hasError(err2) { // fail if api endpoint does not exist
		return
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)

	if hasError(err3) { // fail if resp.Body cannot be read
		return
	}

	var emoji *[]Emoji
	err4 := json.Unmarshal(body, &emoji)

	if hasError(err4) { // fail if body is not json
		return
	}

	// create a path to store images
	err5 := os.MkdirAll(exe+"/"+arg2, os.ModePerm)

	if hasError(err5) { // fail if unable to create path
		return
	}

	for i := 0; i < len(*emoji); i++ {

		// (*emoji)[i] is used instead of *emoji[i] because the compiler reads
		// *emoji[i] as *(emoji[i]), which doesn't exist

		// all images are always .png files, assumption may break in future
		file, err6 := os.Create(exe + arg2 + "/" + (*emoji)[i].Shortcode + ".png")

		if hasError(err6) { // fail if cannot create image
			os.Exit(1) // crash program
		}
		defer file.Close()

		img, err7 := http.Get((*emoji)[i].Url)

		if hasError(err7) { // fail if emoji[i].Url is a not a real URL
			os.Exit(1) //crash program
		}

		defer img.Body.Close()

		// write image to file created earlier
		_, err8 := io.Copy(file, img.Body)

		if hasError(err8) { // fail if img.Body does not download image
			os.Exit(1) // crash program
		}

		fmt.Println(exe + arg2 + "/" + (*emoji)[i].Shortcode + ".png created")
	}

	fmt.Println("Finished!")

}

// helper function
func hasError(err error) bool {
	if err != nil {
		fmt.Println(err.Error)
		return true
	}
	return false
}
