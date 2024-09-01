package spoton

import (
	"fmt"
	"os"
	"spg/internal/vault"
	"sync"
)

type YTdlp struct {
	Binary         string
	DefaultSearch  string
	Format         string
	Postprocessors map[string]string

	// i honestly forgot whether these two work
	GeoBypass        bool
	GeoBypassCountry bool

	// Quiet bool
	NoProgress   bool
	NoOverwrites bool
}

var ytdlpStg YTdlp = YTdlp{
	vault.Settings.Cmd.YtdlpBin,
	"auto",
	"bestaudio/best",
	map[string]string{
		"key":              "FFmpegExtractAudio",
		"preferredcodec":   "mp3",
		"preferredquality": "192",
	},
	true,
	true,
	true,
	true,
}

var mx sync.Mutex
var wg sync.WaitGroup

func (a *Album) Download() []error {
	return a.ToPlaylist().Download()
}

func (p *Playlist) Download() []error {
	errs := make([]error, 0, len(p.Tracks)+1)
	sem := make(chan struct{}, vault.Settings.Cmd.Routines)

	err := os.Mkdir(vault.Settings.Cmd.DownPath, 777)
	if err != nil && os.IsNotExist(err) {
		errs = append(errs, err)
		return errs
	}

	for _, i := range p.Tracks {
		wg.Add(1)
		go func(i PlaylistTrack) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()

			err := i.downloadInaccurate()

			if err != nil {
				mx.Lock()
				errs = append(errs, err)
				mx.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return errs
}

// for the future
func (pt *PlaylistTrack) downloadInaccurate() error {
	return fmt.Errorf("test err for now")
}
