// for storing all kinds of information
// it is kind of located on the lowest level so it can be written/read to/from any level anywhere
// technically i do not know whether it should be in internal since pkg is using it and will break otherwise
// but ill leave it as-is for now
package vault

// settings but named like this so i do not repeat the name and it looks good outise of this file
type Stgs struct {
	Net struct {
		APIendpoint         string
		TokenEndpoint       string
		AlbumTracksRoute    func(string) string
		AlbumInfoRoute      func(string) string
		PlaylistTracksRoute func(string, int, int) string
	}
	Cmd struct {
		DownType string
		DownPath string
		FileType string
		YtdlpBin string
		Routines int
		ID       string
	}
}

var Settings Stgs
