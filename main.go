package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	//	"time"
)

type Emoji struct {
	Shortcode       string // emoji name
	StaticUrl       string // doesn't seem to be used?
	Url             string // emoji url
	VisibleInPicker bool   // always returns false?
}

// BATCH is how many reponses gomoji sends out at once
// increase at your own peril
// well more like the web server's peril,
// higher numbers seem to make gomoji hit GOAWAY responses
const BATCH = 75

func main() {

	flag.Parse()            // parse user input
	var resp *http.Response // API call to mastoAPI emoji
	var err error           // error for resp
	arg := flag.Args()[0]   // should be a url
	arg2 := arg             // used in paths
	var wg sync.WaitGroup   // waitgroup used for downloading

	exe, err := os.Executable()
	if hasError(err) { // fail if can't find executable path
		return
	}

	exe, arg, arg2 = processNames(exe, arg, arg2)

	resp, err = http.Get(arg + "/api/v1/custom_emojis")

	if hasError(err) { // fail if api endpoint does not exist
		return
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if hasError(err) { // fail if resp.Body cannot be read
		return
	}

	var emoji *[]Emoji
	err = json.Unmarshal(body, &emoji)

	if hasError(err) { // fail if body is not json
		return
	}

	// create a path to store images
	err = os.MkdirAll(exe+"/"+arg2, os.ModePerm)

	if hasError(err) { // fail if unable to create path
		return
	}
	// sane amount of emojis
	if len(*emoji) <= BATCH {
		wg.Add(len(*emoji))
		for i := 0; i < len(*emoji); i++ {
			go grabImages(exe, arg, arg2, (*emoji)[i], &wg)
		}
	} else {
		// insane amount of emojis
		// TODO: less gross way of writing this
		i, r := 0, 0
		batches := len(*emoji) / BATCH
		if len(*emoji)%BATCH != 0 {
			r = int(math.Abs(float64(len(*emoji) - BATCH*batches)))
		}
		for o := 0; o < batches; o++ {
			wg.Add(BATCH)
			for p := 0; p < BATCH; p++ {
				go grabImages(exe, arg, arg2, (*emoji)[i], &wg)
				i++
			}
			wg.Wait()
			//			time.Sleep(5 * time.Second)
		}
		wg.Add(r)
		for o := 0; o < r; o++ {
			go grabImages(exe, arg, arg2, (*emoji)[i], &wg)
			i++
		}
	}
	wg.Wait()
	fmt.Println("Finished!")

}

// grabImages...grabs images and downloads them.
func grabImages(exe, arg, arg2 string, emoji Emoji, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(exe + arg2 + "/" + emoji.Shortcode +
		".png")
	if hasError(err) { // fail if cannot create image
		return
	}
	defer file.Close()

	img, err := http.Get(emoji.Url)
	if hasError(err) { // fail if emoji.Url is a not a real URL
		return
	}
	defer img.Body.Close()

	// write image to file created earlier
	_, err = io.Copy(file, img.Body)
	if hasError(err) { // fail if img.Body does not download image
		return
	}

	fmt.Println(exe + arg2 + "/" + emoji.Shortcode + ".png created")

}

// hasError is a helper function to print errors.
func hasError(err error) bool {
	if err != nil {
		fmt.Printf("error: %v", err)
		return true
	}
	return false
}

// processNames converts the following into acceptable outputs for paths
// and file names.
func processNames(exe, arg, arg2 string) (string, string, string) {
	exe = strings.TrimSuffix(exe, "gomoji") // hardcoded name, bad
	if !strings.HasPrefix(arg, "https://") &&
		!strings.HasPrefix(arg, "http://") {
		// add https:// if arg doesn't have it
		arg = "https://" + arg
	} else if strings.HasPrefix(arg, "http://") {
		// change http:// to https:// in arg
		arg2 = strings.TrimPrefix(arg2, "http://")
		arg = strings.Replace(arg, "http://", "https://", 1)
	} else {
		// trim prefix in arg2 to avoid illegal file names
		arg2 = strings.TrimPrefix(arg2, "https://")
	}

	return exe, arg, arg2
}
