package emoji

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// connection represents a connection sent out, contains an Emoji and an
// indiciation if it failed.
type connection struct {
	Emoji      api
	Downloaded bool
}

// Download grabs the emoji image from the server and saves it to a file.
func (c *connection) download(dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(dir + "/" + c.Emoji.Shortcode + ".png")
	if hasError(err) { // fail if cannot create image
		return
	}
	defer file.Close()

	img, err := http.Get(c.Emoji.URL)
	if hasError(err) { // fail for a variety of reasons
		return
	}
	defer img.Body.Close()

	// write image to file created earlier
	_, err = io.Copy(file, img.Body)
	if hasError(err) { // fail if img.Body does not download image
		return
	}

	if Verbose.Downloads {
		fmt.Printf("%v/%v.png created.\n", dir, c.Emoji.Shortcode)
	}

	c.Downloaded = true
}
