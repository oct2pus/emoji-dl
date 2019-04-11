package emoji

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
