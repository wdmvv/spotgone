package cmd

import (
	"flag"
	"fmt"
	"spg/internal/vault"
)

func CmdArgs() error {
	var downType *string = flag.String("t", "playlist", "download type, either album or track")
	var downFld *string = flag.String("d", "Downloads", "download path, by default creates Downloads folder")
	var ytdlpBin *string = flag.String("b", "", "path to the ytdlp binary")
	var routines *int = flag.Int("r", 5, "number of routines for downloader, default 5")
	var id *string = flag.String("id", "", "mandatory album/playlist id")

	flag.Parse()

	if *id == "" {
		return fmt.Errorf("error: download id must not be empty")
	}
	// TODO: add album/track validation
	// TODO: add id parsing so it supports links along with raw id

	vault.Settings.APIendpoint = "https://api.spotify.com/v1"
	vault.Settings.TokenEndpoint = "https://accounts.spotify.com/api/token"

	vault.Settings.Cmd.Type = *downType
	vault.Settings.Cmd.DownPath = *downFld
	vault.Settings.Cmd.YtdlpBin = *ytdlpBin
	vault.Settings.Cmd.Routines = *routines
	vault.Settings.Cmd.ID = *id

	return nil
}
