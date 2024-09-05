# spotgone
wdmvv/spot-on but in go because prev version sucks!<br>
If you are not familiar with what spoton is in first place - last year I got annoyed by spotify and wrote a thing that takes all songs info from spotify's album/playlist and downloads from youtube via yt-dlp. I am aware that there are certain other tools that can download from spotify stream(?) directly, nonetheless, I decided to rewrite old ver with youtube downloads in go.<br>
<br>
## Installation
```
git clone https://github.com/wdmvv/spotgone
cd spotgone
go build && ./spg
```
## Usage
### -id
    Mandatory, ID of the playlist/album, can also be a link
    `./spg ... -id <id>`
### -b
    Mandatory, path to the ytdlp binary which is going to be used to download content
    `./spg ... -b "/path/to/binary"...`
### -t
    Used to specify type of the spotify download, can be either album or playlist (or any case of a/p)
    `./spg ... -t playlist ...`
### -d
    Download path, by default creates folder named "Downloads"
    `./spg ... -d "/some/path" ...`
### -f
    Downloaded file format, can be any of aac, alac, flac, m4a, mp3, opus, vorbis, wav, mp3 by default
    `./spg ... -f m4a ...`
### -r
    Number of downloader goroutines to launch, 5 by default
    `./spg ... -r 10 ...`

## TODO
<ul>
<li>Make all paths relative to the root and not execution start</li>
<li>Replace []error with chan error for faster? results</li>
<li>Add contexts with timeouts to all requests</li>
<li>Add output formatting option in cmd args</li>
<li>Glue paths together</li>
<li>metadata embedder</li>
<li>...and some other original spoton features</li>
</ul>
