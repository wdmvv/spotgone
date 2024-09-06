package cmd

import (
	"flag"
	"fmt"
	"net/url"
	"regexp"
	"spg/internal/vault"
	"strings"
)

func CmdArgs() error {
	var id *string = flag.String("id", "", "album or playlist id, links work too")
	var downType *string = flag.String("t", "playlist", "download type of spotify item, either album or playlist")
	var downPath *string = flag.String("d", "Downloads", "download path, by default creates Download folder")
	var fileType *string = flag.String("f", "mp3", "file type of downloads, mp3 by default")
	var ytdlpBin *string = flag.String("b", "", "path to the ytdlp binary to use")
	var routines *int = flag.Int("r", 5, "number of routines for downloader, 5 by default")

	flag.Parse()

	// TODO: add path validations - requires sticking path together in first place

	// i dont know whether this is a good pattern, will keep until further notice

	vault.Settings.Net.APIendpoint = "https://api.spotify.com/v1"
	vault.Settings.Net.TokenEndpoint = "https://accounts.spotify.com/api/token"

	vault.Settings.Net.AlbumTracksRoute = func(id string) string {
		return vault.Settings.Net.APIendpoint +
			fmt.Sprintf("/albums/%s/tracks?limit=50", id)
	}
	vault.Settings.Net.AlbumInfoRoute = func(id string) string {
		return vault.Settings.Net.APIendpoint +
			fmt.Sprintf("/albums/%s", id)
	}
	vault.Settings.Net.PlaylistTracksRoute = func(id string, offset int, limit int) string {
		return vault.Settings.Net.APIendpoint +
			fmt.Sprintf("/playlists/%s/tracks?offset=%d&limit=%d", id, offset, limit)
	}

	vault.Settings.Cmd.DownTypeRaw = *downType
	vault.Settings.Cmd.DownPath = *downPath
	vault.Settings.Cmd.FileType = *fileType
	vault.Settings.Cmd.YtdlpBin = *ytdlpBin
	vault.Settings.Cmd.Routines = *routines
	vault.Settings.Cmd.ID = *id

	return validateCmd()
}

// one could say that validating after setting values is unwise (and makes no sense)
// and i agree but in this case there is nothing sensitive happening
// plus it'd mess with cmd function
func validateCmd() error {
	if err := validateId(&vault.Settings.Cmd.ID); err != nil {
		return err
	}
	if err := validateDownType(vault.Settings.Cmd.DownTypeRaw,
		&vault.Settings.Cmd.DownType); err != nil {
		return err
	}
	if err := validateFileType(vault.Settings.Cmd.FileType); err != nil {
		return err
	}

	return nil
}

func validateId(id *string) error {
	*id = strings.Trim(*id, " ")
	if *id == "" {
		return fmt.Errorf("error: download id must not be empty")
	}
	u, err := url.Parse(*id)
	if err != nil {
		return err
	}
	if u.Path != "" {
		pth := strings.Split(u.Path, "/")
		*id = pth[len(pth)-1]
	}
	return nil
}

func validateDownType(tp string, t *int) error {
	rgalb := regexp.MustCompile("(?i)(album)\\s*(?-i)|([aA]{1}(\\s*)$)")
	rgpls := regexp.MustCompile("(?i)(playlist)\\s*(?-i)|([pP]{1}(\\s*)$)")

	if !rgalb.Match([]byte(tp)) && !rgpls.Match([]byte(tp)) {
		return fmt.Errorf("invalid download type %s", tp)
	}
	*t = 1

	if rgalb.Match([]byte(tp)) {
		*t = 0
	}
	return nil
}

func validateFileType(tp string) error {
	re := regexp.MustCompile("(?i)(aac|alac|flac|m4a|mp3|opus|vorbis|wav)(?-i)")
	if !re.Match([]byte(tp)) {
		return fmt.Errorf("invalid download type %s", tp)
	}
	return nil
}
