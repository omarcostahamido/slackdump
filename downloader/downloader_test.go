package downloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rusq/slackdump/v2/internal/mocks/mock_downloader"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

var (
	file1 = slack.File{ID: "f1", Name: "filename1.ext", URLPrivateDownload: "file1_url", Size: 100}
	file2 = slack.File{ID: "f2", Name: "filename2.ext", URLPrivateDownload: "file2_url", Size: 200}
	file3 = slack.File{ID: "f3", Name: "filename3.ext", URLPrivateDownload: "file3_url", Size: 300}
	file4 = slack.File{ID: "f4", Name: "filename4.ext", URLPrivateDownload: "file4_url", Size: 400}
	file5 = slack.File{ID: "f5", Name: "filename5.ext", URLPrivateDownload: "file5_url", Size: 500}
	file6 = slack.File{ID: "f6", Name: "filename6.ext", URLPrivateDownload: "file6_url", Size: 600}
	file7 = slack.File{ID: "f7", Name: "filename7.ext", URLPrivateDownload: "file7_url", Size: 700}
	file8 = slack.File{ID: "f8", Name: "filename8.ext", URLPrivateDownload: "file8_url", Size: 800}
	file9 = slack.File{ID: "f9", Name: "filename9.ext", URLPrivateDownload: "file9_url", Size: 900}
)

func TestSlackDumper_SaveFileTo(t *testing.T) {
	tmpdir := t.TempDir()

	type fields struct {
		l       *rate.Limiter
		retries int
		workers int
	}
	type args struct {
		ctx context.Context
		dir string
		f   *slack.File
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expectFn func(mc *mock_downloader.MockDownloader)
		want     int64
		wantErr  bool
	}{
		{
			"ok",
			fields{
				l:       rate.NewLimiter(defLimit, 1),
				retries: defRetries,
				workers: defNumWorkers},
			args{
				context.Background(),
				tmpdir,
				&file1,
			},
			func(mc *mock_downloader.MockDownloader) {
				mc.EXPECT().
					GetFile("file1_url", gomock.Any()).
					Return(nil)
			},
			int64(file1.Size),
			false,
		},
		{
			"getfile rekt",
			fields{
				l:       rate.NewLimiter(defLimit, 1),
				retries: defRetries,
				workers: defNumWorkers},
			args{
				context.Background(),
				tmpdir,
				&file2,
			},
			func(mc *mock_downloader.MockDownloader) {
				mc.EXPECT().
					GetFile("file2_url", gomock.Any()).
					Return(errors.New("rekt"))
			},
			int64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mc := mock_downloader.NewMockDownloader(ctrl)

			tt.expectFn(mc)

			sd := &Client{
				client:  mc,
				limiter: tt.fields.l,
				retries: tt.fields.retries,
				workers: tt.fields.workers,
			}
			got, err := sd.SaveFile(tt.args.ctx, tt.args.dir, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackDumper.SaveFileTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlackDumper.SaveFileTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackDumper_saveFile(t *testing.T) {
	tmpdir := t.TempDir()

	type fields struct {
		l       *rate.Limiter
		retries int
		workers int
	}
	type args struct {
		ctx context.Context
		dir string
		f   *slack.File
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expectFn func(mc *mock_downloader.MockDownloader)
		want     int64
		wantErr  bool
	}{
		{
			"ok",
			fields{
				l:       rate.NewLimiter(defLimit, 1),
				retries: defRetries,
				workers: defNumWorkers,
			},
			args{
				context.Background(),
				tmpdir,
				&file1,
			},
			func(mc *mock_downloader.MockDownloader) {
				mc.EXPECT().
					GetFile("file1_url", gomock.Any()).
					Return(nil)
			},
			int64(file1.Size),
			false,
		},
		{
			"getfile rekt",
			fields{
				l:       rate.NewLimiter(defLimit, 1),
				retries: defRetries,
				workers: defNumWorkers,
			},
			args{
				context.Background(),
				tmpdir,
				&file2,
			},
			func(mc *mock_downloader.MockDownloader) {
				mc.EXPECT().
					GetFile("file2_url", gomock.Any()).
					Return(errors.New("rekt"))
			},
			int64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mc := mock_downloader.NewMockDownloader(ctrl)

			tt.expectFn(mc)

			sd := &Client{
				client:  mc,
				limiter: tt.fields.l,
				retries: tt.fields.retries,
				workers: tt.fields.workers,
			}
			got, err := sd.saveFile(tt.args.ctx, tt.args.dir, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackDumper.saveFileWithLimiter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlackDumper.saveFileWithLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filename(t *testing.T) {
	type args struct {
		f *slack.File
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"file1", args{&file1}, "f1-filename1.ext"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filename(tt.args.f); got != tt.want {
				t.Errorf("filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackDumper_newFileDownloader(t *testing.T) {
	t.Parallel()
	tl := rate.NewLimiter(defLimit, 1)
	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	t.Run("ensure file downloader is running", func(t *testing.T) {
		mc := mock_downloader.NewMockDownloader(gomock.NewController(t))
		sd := Client{
			client:  mc,
			limiter: tl,
			retries: 3,
			workers: 4,
		}

		mc.EXPECT().
			GetFile(file9.URLPrivateDownload, gomock.Any()).
			Return(nil).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
		defer cancel()
		filesC := make(chan *slack.File, 1)
		filesC <- &file9
		close(filesC)

		done, err := sd.AsyncDownloader(ctx, tmpdir, filesC)
		require.NoError(t, err)

		<-done
		filename := filepath.Join(tmpdir, filename(&file9))
		assert.FileExists(t, filename)

	})
}

func TestSlackDumper_worker(t *testing.T) {
	t.Parallel()
	tl := rate.NewLimiter(defLimit, 1)
	tmpdir := t.TempDir()

	t.Run("sending a single file", func(t *testing.T) {
		mc := mock_downloader.NewMockDownloader(gomock.NewController(t))
		sd := Client{
			client:  mc,
			limiter: tl,
			retries: defRetries,
			workers: defNumWorkers,
		}

		mc.EXPECT().
			GetFile(file1.URLPrivateDownload, gomock.Any()).
			Return(nil).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		reqC := make(chan FileRequest, 1)
		reqC <- FileRequest{Directory: tmpdir, File: &file1}
		close(reqC)

		sd.worker(ctx, reqC)
		assert.FileExists(t, filepath.Join(tmpdir, filename(&file1)))
	})
	t.Run("getfile error", func(t *testing.T) {
		mc := mock_downloader.NewMockDownloader(gomock.NewController(t))
		sd := Client{
			client:  mc,
			limiter: tl,
			retries: defRetries,
			workers: defNumWorkers,
		}

		mc.EXPECT().
			GetFile(file1.URLPrivateDownload, gomock.Any()).
			Return(errors.New("rekt")).
			Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		reqC := make(chan FileRequest, 1)
		reqC <- FileRequest{Directory: tmpdir, File: &file1}
		close(reqC)

		sd.worker(ctx, reqC)
		_, err := os.Stat(filepath.Join(tmpdir, filename(&file1)))
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("cancelled context", func(t *testing.T) {
		mc := mock_downloader.NewMockDownloader(gomock.NewController(t))
		sd := Client{
			client:  mc,
			limiter: tl,
			retries: defRetries,
			workers: defNumWorkers,
		}

		reqC := make(chan FileRequest, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		cancel()

		sd.worker(ctx, reqC)
	})
}

func Test_mkdir(t *testing.T) {
	dir := t.TempDir()

	testFile := filepath.Join(dir, "existing_file")
	if err := os.WriteFile(testFile, []byte("I should not be moved"), 0666); err != nil {
		t.Fatalf("failed to create a test file: %s", err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty name",
			args{},
			true,
		},
		{
			"new directory",
			args{filepath.Join(dir, "test1")},
			false,
		},
		{
			"directory already exists",
			args{filepath.Join(dir, "test1")},
			false,
		},
		{
			"another dir",
			args{filepath.Join(dir, "test2")},
			false,
		},
		{
			"object with such name exists, and is a file",
			args{testFile},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mkdir(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_startWorkers(t *testing.T) {
	t.Run("check that start actually starts workers", func(t *testing.T) {
		const qSz = 10

		ctrl := gomock.NewController(t)
		dc := mock_downloader.NewMockDownloader(ctrl)
		cl := Client{
			client:  dc,
			limiter: rate.NewLimiter(5000, 1),
			workers: defNumWorkers,
		}

		dc.EXPECT().GetFile(gomock.Any(), gomock.Any()).Times(qSz).Return(nil)

		fileQueue := makeFileReqQ(qSz, t.TempDir())
		fileChan := slice2chan(fileQueue, defFileBufSz)
		wg := cl.startWorkers(context.Background(), fileChan)

		wg.Wait()
	})
}

// slice2chan takes the slice of []T, create a chan T and sends all elements of
// []T to it.  It closes the channel after all elements are sent.
func slice2chan[T any](input []T, bufSz int) <-chan T {
	output := make(chan T, bufSz)
	go func() {
		defer close(output)
		for _, v := range input {
			output <- v
		}
	}()
	return output
}

func TestClient_Start(t *testing.T) {
	t.Run("make sure structures initialised", func(t *testing.T) {
		c := clientWithMock(t)

		c.Start(context.Background())
		defer c.Stop()

		assert.True(t, c.started)
		assert.NotNil(t, c.wg)
		assert.NotNil(t, c.fileRequests)
	})
}

func TestClient_Stop(t *testing.T) {
	t.Run("ensure stopped", func(t *testing.T) {
		c := clientWithMock(t)
		c.Start(context.Background())
		assert.True(t, c.started)

		c.Stop()
		assert.False(t, c.started)
		assert.Nil(t, c.fileRequests)
		assert.Nil(t, c.wg)
	})
	t.Run("stop on stopped downloader does nothing", func(t *testing.T) {
		c := clientWithMock(t)
		c.Stop()
		assert.False(t, c.started)
		assert.Nil(t, c.fileRequests)
		assert.Nil(t, c.wg)
	})
}

func clientWithMock(t *testing.T) *Client {
	ctrl := gomock.NewController(t)
	dc := mock_downloader.NewMockDownloader(ctrl)
	c := &Client{
		client:  dc,
		limiter: rate.NewLimiter(5000, 1),
		workers: defNumWorkers,
	}
	return c
}

func TestClient_DownloadFile(t *testing.T) {
	dir := filepath.Join(t.TempDir())
	t.Run("returns error on stopped downloader", func(t *testing.T) {
		c := clientWithMock(t)
		err := c.DownloadFile(dir, slack.File{ID: "xx", Name: "tt"})
		assert.ErrorIs(t, err, ErrNotStarted)
	})
	t.Run("ensure that file is placed on the queue", func(t *testing.T) {
		c := clientWithMock(t)
		c.Start(context.Background())

		c.client.(*mock_downloader.MockDownloader).EXPECT().
			GetFile(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		err := c.DownloadFile(dir, file1)
		assert.NoError(t, err)

		c.Stop()
	})
}
