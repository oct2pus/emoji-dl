package emoji

type batch struct {
	Size      int
	Downloads bool
}

// Batch controls how batching downloads works.
//
// Batch.Size determines how many emoji images to try and download at once.
//
// Batch.Downloads enables or disables batch downloading.
var Batch = batch{Size: 75, Downloads: true}
