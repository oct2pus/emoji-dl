package emoji

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

var complete = make(chan bool)

// A Collection is all the connections sent out to a server.
type Collection []connection

// NewCollection returns a slice of connection objects.
// Downloaded will always be initalized as false.
func NewCollection(site string) (Collection, error) {
	var emoji *[]api
	resp, err := http.Get(site + "/api/v1/custom_emojis")
	if hasError(err) {
		return []connection{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if hasError(err) { // fail if resp.Body cannot be read
		return []connection{}, err
	}

	err = json.Unmarshal(body, &emoji)
	if hasError(err) { // fail if body is not json
		return []connection{}, err
	}

	p := make([]connection, 0)
	for i := range *emoji {
		p = append(p, connection{Emoji: (*emoji)[i], Downloaded: false})
	}
	return p, nil
}

// CollectFailed returns a version of a Collection that has stripped out all
//
func (c Collection) collectFailed() Collection {
	p := make([]connection, 0)
	for _, conn := range c {
		if !conn.Downloaded {
			p = append(p, conn)
		}
	}
	return p
}

// hasFailed returns if any individual connection in the Collection has not
// downloaded. A value of 0 is always true.
func (c Collection) hasFailed() bool {
	if len(c) > 0 {
		for _, ele := range c {
			if !ele.Downloaded {
				return true
			}
		}
	}
	return false
}

// DownloadAll downloads all images in the emoji collection.
func (c Collection) DownloadAll(dir string) {
	size := Batch.Size
	if !Batch.Downloads {
		size = len(c)
	}
	var wg sync.WaitGroup // waitgroup used for downloading
	for c.hasFailed() {
		complete <- false
		if len(c) <= size {
			wg.Add(len(c))
			for i := 0; i < len(c); i++ {
				go c[i].download(dir, &wg)
			}
			wg.Wait()
		} else {
			b := (len(c) / size)
			count := 0
			for i := 0; i < b; i++ {
				wg.Add(size)
				for o := 0; o < size; o++ {
					go c[count].download(dir, &wg)
					count++
				}
				wg.Wait()
			}
		}
		c = c.collectFailed()
	}
	complete <- true
}

// DownloadFinished returns if DownloadAll has completed.
// returns false if not initalized to begin with.
func (c Collection) DownloadFinished() bool {
	select {
	case result := <-complete:
		return result
	default:
		return false
	}
}
