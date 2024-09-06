package spoton

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"spg/internal/vault"
	"sync"
)

// i will be using cmd tag to mark fields that are required as a part of processing
// fields should already contain formatted keys i.e --some-key <value>
type YTdlp struct {
	// bools are declared in a way that if they are true, then you have to use them in cmd
	// this way it is a simple if(bool) instead of manual param check

	NoProgress   bool `cmd:"--no-progress"`
	NoOverwrites bool `cmd:"--no-overwrites"`
	ExtractAudio bool `cmd:"-x"`

	// search engine
	DefaultSearch string `cmd:"--default-search"`
	// search format
	Format string `cmd:"--format"`
	// what type should save files be
	FileType string `cmd:"--audio-format"`
	// audio quality
	Quality string `cmd:"--audio-quality"`
}

var ytdlpStg YTdlp = YTdlp{
	NoProgress:   true,
	NoOverwrites: true,
	ExtractAudio: true,

	DefaultSearch: "ytsearch",
	Format:        "bestaudio/best",
	FileType:      vault.Settings.Cmd.FileType,
	Quality:       "192k",
}

var mx sync.Mutex
var wg sync.WaitGroup

func (a *Album) Download() []error {
	return a.ToPlaylist().Download()
}

func (p *Playlist) Download() []error {
	// TODO: figure out why it does not set itself in var declaration
	ytdlpStg.FileType = vault.Settings.Cmd.FileType

	errs := make([]error, 0, len(p.Tracks)+1)
	sem := make(chan struct{}, vault.Settings.Cmd.Routines)

	err := os.Mkdir(vault.Settings.Cmd.DownPath, 775)
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

func (pt *PlaylistTrack) downloadInaccurate() error {
	args := ytdlpStg.toCmd(pt.Name + "." + vault.Settings.Cmd.FileType)
	args = append(args, fmt.Sprintf("\"%s - %s\"", pt.ArtistStr(), pt.Name))
	var cmd *exec.Cmd

	// someone has to test this on windows, works on my machine though
	if runtime.GOOS == "windows" {
		cmd = exec.Command(vault.Settings.Cmd.YtdlpBin, args...)
	} else {
		cmd = exec.Command("./"+vault.Settings.Cmd.YtdlpBin, args...)
	}
	_, err := cmd.Output()
	return err
}

// function for converting ytdlp structure into cmd arguments string
// output is how output file is going to be named
// for now it'd be Song.format
// in the future should add options to the cmd args
func (y *YTdlp) toCmd(output string) []string {
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
			continue
			// tmp := cmdstr + " "

			// iter := v.Field(i).MapRange()
			// for iter.Next() {
			// 	tmp += fmt.Sprintf("%s:%s ", iter.Key(), iter.Value())
			// }
			// out = append(out, tmp)
		default:
			out = append(out, cmdstr)
			out = append(out, fmt.Sprintf("\"%v\"", v.Field(i)))
		}
	}

	// adding this manually otherwise i'd have to check somewhere else
	out = append(out, "--output")

	// ugly as hell
	out = append(out,
		fmt.Sprintf("\".%s%s\"", string(os.PathSeparator), path.Join(".", vault.Settings.Cmd.DownPath, output)))
	out = append(out, fmt.Sprintf("\"%s\"", output))

	return out
}
