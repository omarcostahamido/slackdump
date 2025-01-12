package downloader

import (
	"math/rand"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func Test_fltSeen(t *testing.T) {
	t.Run("ensure that we don't get dup files", func(t *testing.T) {
		source := []FileRequest{
			{Directory: "x", File: &file1},
			{Directory: "x", File: &file2},
			{Directory: "a", File: &file2}, // this should appear, different dir
			{Directory: "x", File: &file3},
			{Directory: "x", File: &file3}, // duplicate
			{Directory: "x", File: &file3}, // duplicate
			{Directory: "x", File: &file4},
			{Directory: "x", File: &file5},
			{Directory: "y", File: &file5}, // this should appear, different dir
		}
		want := []FileRequest{
			{Directory: "x", File: &file1},
			{Directory: "x", File: &file2},
			{Directory: "a", File: &file2},
			{Directory: "x", File: &file3},
			{Directory: "x", File: &file4},
			{Directory: "x", File: &file5},
			{Directory: "y", File: &file5},
		}

		filesC := make(chan FileRequest)
		go func() {
			defer close(filesC)
			for _, f := range source {
				filesC <- f
			}
		}()
		dlqC := fltSeen(filesC)

		var got []FileRequest
		for f := range dlqC {
			got = append(got, f)
		}
		assert.Equal(t, want, got)
	})
}

func BenchmarkFltSeen(b *testing.B) {
	const numReq = 100_000
	input := makeFileReqQ(numReq, b.TempDir())
	inputC := make(chan FileRequest)
	go func() {
		defer close(inputC)
		for _, req := range input {
			inputC <- req
		}
	}()

	outputC := fltSeen(inputC)

	for n := 0; n < b.N; n++ {
		for out := range outputC {
			_ = out
		}
	}

}

func makeFileReqQ(numReq int, dir string) []FileRequest {
	reqQ := make([]FileRequest, numReq)
	for i := 0; i < numReq; i++ {
		reqQ[i] = randomFileReq(dir)
	}
	return reqQ
}

func randomFileReq(dirname string) FileRequest {
	return FileRequest{Directory: dirname, File: &slack.File{ID: randString(8), Name: randString(12)}}
}

func randString(sz int) string {
	const (
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		chrstSz = len(charset)
	)
	var ret = make([]byte, sz)
	for i := 0; i < sz; i++ {
		ret[i] = charset[rand.Int()%chrstSz]
	}
	return string(ret)
}
