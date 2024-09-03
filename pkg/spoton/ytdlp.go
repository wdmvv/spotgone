package spoton

import (
	"fmt"
	"os"
	"reflect"
	"spg/internal/vault"
	"strings"
	"sync"
)

// i will be using cmd tag to mark fields that are required as a part of processing
// fields should already contain formatted keys i.e --some-key <value>
// cmd: formatted option for ytdlp Binary
type YTdlp struct {
	Binary        string
	DefaultSearch string `cmd:"--default-search"`
	Format        string `cmd:"--format"`
	// wait this does not work as expected?? what do i do
	Postprocessors map[string]string `cmd:"--postprocessor-args"`

	// Quiet bool
	// bools are declared in a way that if they are true, then you have to use them in cmd
	// this way it is a simple if(bool) instead of manual param check

	NoProgress   bool `cmd:"--no-progress"`
	NoOverwrites bool `cmd:"--no-overwrites"`

	// basically savepath
	Outtmpl string `cmd:"--output"`
}

var ytdlpStg YTdlp = YTdlp{
	Binary:        vault.Settings.Cmd.YtdlpBin,
	DefaultSearch: "auto",
	Format:        "bestaudio/best",
	Postprocessors: map[string]string{
		"key":              "FFmpegExtractAudio",
		"preferredcodec":   "mp3",
		"preferredquality": "192",
	},
	NoProgress:   true,
	NoOverwrites: true,
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

	fmt.Println("DEBUG:", ytdlpStg.toCmd())

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
	return nil
}

// function for converting ytdlp structure into cmd arguments string
func (y *YTdlp) toCmd() string {
	out := make([]string, 0)
	v := reflect.ValueOf(*y)
	for i := 0; i < v.NumField(); i++ {
		cmdstr := reflect.TypeOf(*y).Field(i).Tag.Get("cmd")

		switch reflect.TypeOf(*y).Field(i).Type.Kind() {
		case reflect.Bool:
			if fmt.Sprintf("%v", v.Field(i)) == "true" {
				out = append(out, cmdstr)
			}
		case reflect.Map:
			// since postprocessor stuff is not supported just yet
			continue

			// tmp := cmdstr + " "

			// iter := v.Field(i).MapRange()
			// for iter.Next() {
			// 	tmp += fmt.Sprintf("%s:%s ", iter.Key(), iter.Value())
			// }
			// out = append(out, tmp)
		default:
			out = append(out, fmt.Sprintf("%s %v", cmdstr, v.Field(i)))
		}
	}
	return strings.Join(out, " ")
}
