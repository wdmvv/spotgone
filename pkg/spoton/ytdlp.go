package spoton

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

func (a *Album) Download() {
	// a.toPlaylist().Download() // teehee, might change later
}

func (p *Playlist) Download() {}
