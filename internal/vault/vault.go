// for storing any kinds of information
// it is kind of located on the lowest level so it can be written/read to/from any level anywhere
package vault

// settings but named like this so i do not repeat the name
type Stgs struct {
	APIendpoint   string
	TokenEndpoint string
	Cmd           struct {
		Routines int
		DownPath string
		YtdlpBin string
		Type     string
		ID       string
	}
}

var Settings Stgs
