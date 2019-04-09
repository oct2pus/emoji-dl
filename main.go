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
	"sync"
)

type batch struct {
	Size      int
	Downloads bool
}

// Batch is an accessor to two configurable elements that have default values.
//
// - `Batch.Size` is how many to try and download at once
//	 Default: 75
// - `Batch.Downloads` enables or disables batch downloading;
//	 Default: true
var Batch = batch{Size: 75, Downloads: true}

// Emoji represents an Emoji object returned by a MastoAPI.
type Emoji struct {
	Shortcode       string // emoji name
	StaticURL       string // doesn't seem to be used?
	URL             string // emoji url
	VisibleInPicker bool   // always returns false?
}

// Connection represents a connection sent out, contains an Emoji and an
// indiciation if it failed.
type Connection struct {
	Emoji      Emoji
	Downloaded bool
}

// A Collection is all the Connections sent out to a server.
type Collection []Connection

// NewCollection returns a slice of Connection objects.
// Downloaded will always be initalized as false.
func NewCollection(emoji *[]Emoji) Collection {
	p := make([]Connection, 0)
	for i := range *emoji {
		p = append(p, Connection{Emoji: (*emoji)[i], Downloaded: false})
	}
	return p
}

// CollectFailed returns a version of a Collection that has stripped out all
//
func (c Collection) CollectFailed() Collection {
	p := make([]Connection, 0)
	for _, conn := range c {
		if !conn.Downloaded {
			p = append(p, conn)
		}
	}
	return p
}

// HasFailed returns if any individual Connection in the Collection has not
// downloaded. A value of 0 is always true.
func (c Collection) HasFailed() bool {
	if len(c) > 0 {
		for _, ele := range c {
			if !ele.Downloaded {
				return true
			}
		}
	}
	return false
}

func main() {

	flag.Parse() // parse user input

	var resp *http.Response   // API call to mastoAPI emoji
	var err error             // error for resp
	var emoji *[]Emoji        // contains unmarshalled data
	siteURL := flag.Args()[0] // should be a url
	pathURL := siteURL        // used in paths
	var wg sync.WaitGroup     // waitgroup used for downloading

	wd, err := os.Getwd()
	if hasError(err) { // fail if can't find working directory
		return
	}
	wd += "/"

	siteURL, pathURL = processNames(siteURL, pathURL)

	resp, err = http.Get(siteURL + "/api/v1/custom_emojis")
	if hasError(err) { // fail if api endpoint does not exist
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if hasError(err) { // fail if resp.Body cannot be read
		return
	}

	err = json.Unmarshal(body, &emoji)
	if hasError(err) { // fail if body is not json
		return
	}

	conn := NewCollection(emoji)

	// create a path to store images
	err = os.MkdirAll(wd+"/"+pathURL, os.ModePerm)
	if hasError(err) { // fail if unable to create path
		return
	}

	amount := len(conn)
	if Batch.Downloads {
		amount = Batch.Size
	}

	for conn.HasFailed() {
		if len(conn) <= amount {
			wg.Add(len(conn))
			for i := 0; i < len(conn); i++ {
				go grabImages(wd, siteURL, pathURL, &conn[i], &wg)
			}
			wg.Wait()
		} else {
			println(len(conn))
			b := (len(conn) / Batch.Size)
			count := 0
			for i := 0; i < b; i++ {
				wg.Add(amount)
				for o := 0; o < amount; o++ {
					go grabImages(wd, siteURL, pathURL, &conn[count], &wg)
					count++
				}
				wg.Wait()
			}
		}
		conn = conn.CollectFailed()
	}
	fmt.Println("Finished!")

}

// grabImages...grabs images and downloads them.
func grabImages(wd, arg, arg2 string, conn *Connection, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(wd + arg2 + "/" + conn.Emoji.Shortcode + ".png")
	if hasError(err) { // fail if cannot create image
		return
	}
	defer file.Close()

	img, err := http.Get(conn.Emoji.URL)
	if hasError(err) { // fail for a variety of reasons
		return
	}
	defer img.Body.Close()

	// write image to file created earlier
	_, err = io.Copy(file, img.Body)
	if hasError(err) { // fail if img.Body does not download image
		return
	}

	conn.Downloaded = true
	fmt.Printf("\n%v%v/%v.png created", wd, arg2, conn.Emoji.Shortcode)
}

// hasError is a helper function to print errors.
func hasError(err error) bool {
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return true
	}
	return false
}

// processNames converts the following into acceptable outputs for paths
// and file names.
func processNames(arg, arg2 string) (string, string) {
	// add https:// if arg doesn't have it

	if !strings.HasPrefix(arg, "https://") &&
		!strings.HasPrefix(arg, "http://") {
		arg = "https://" + arg
	} else if strings.HasPrefix(arg, "http://") {
		// change http:// to https:// in arg
		arg2 = strings.TrimPrefix(arg2, "http://")
		arg = strings.Replace(arg, "http://", "https://", 1)
	} else {
		// trim prefix in arg2 to avoid illegal file names
		arg2 = strings.TrimPrefix(arg2, "https://")
	}

	return arg, arg2
}
