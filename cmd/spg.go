package cmd

import (
	"flag"
	"fmt"
	"spg/internal/vault"
)

func CmdArgs() error {
	var spotType *string = flag.String("t", "playlist", "download type of spotify item, either album or playlist")
	var downPath *string = flag.String("d", "Downloads", "download path, by default creates Download folder")
	var fileType *string = flag.String("f", "mp3", "file type of downloads, mp3 by default")
	var ytdlpBin *string = flag.String("b", "", "path to the ytdlp binary to use")
	var routines *int = flag.Int("r", 5, "number of routines for downloader, 5 by default")
	var id *string = flag.String("id", "", "album or playlist id, links work too")

	flag.Parse()

	if *id == "" {
		return fmt.Errorf("error: download id must not be empty")
	} else if *ytdlpBin == "" {
		return fmt.Errorf("error: yt dlp binary must not be empty")

	}
	// TODO: add album/track validation
	// TODO: add path validations
	// TODO: add id parsing so it supports links along with raw id

	vault.Settings.APIendpoint = "https://api.spotify.com/v1"
	vault.Settings.TokenEndpoint = "https://accounts.spotify.com/api/token"

	vault.Settings.Cmd.DownType = *spotType
	vault.Settings.Cmd.DownPath = *downPath
	vault.Settings.Cmd.FileType = *fileType
	vault.Settings.Cmd.YtdlpBin = *ytdlpBin
	vault.Settings.Cmd.Routines = *routines
	vault.Settings.Cmd.ID = *id

	return nil
}
