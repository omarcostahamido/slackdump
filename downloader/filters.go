package downloader

import (
	"log"
)

// fltSeen filters the files from filesC to ensure that no duplicates
// are downloaded.
func fltSeen(filesC <-chan FileRequest) <-chan FileRequest {
	dlQ := make(chan FileRequest)
	go func() {
		// closing stop will lead to all worker goroutines to terminate.
		defer close(dlQ)

		// seen contains file ids that already been seen,
		// so we don't download the same file twice
		seen := make(map[string]bool)
		// files queue must be closed by the caller (see DumpToDir.(1))
		for f := range filesC {
			id := f.File.ID + f.Directory
			if _, ok := seen[id]; ok {
				log.Printf("already seen %s, skipping", filename(f.File))
				continue
			}
			seen[id] = true
			dlQ <- f
		}
	}()
	return dlQ
}
