# spotgone
wdmvv/spot-on but in go because prev version sucks!<br>
If you are not familiar with what spoton is in first place - last year I got annoyed by spotify and wrote a thing that takes all songs info from spotify's album/playlist and downloads from youtube via yt-dlp. I am aware that there are certain other tools that can download from spotify stream(?) directly, nonetheless, I decided to rewrite old ver with youtube downloads in go.<br>
<br><br>
grab ytdlp binary here - https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#installation<br>
In the future ill add proper installation instructions and full(ish) documentation, for now feel free to explore code yourself and figure out how it works.
<br>

TODO:<br>
<ul>
<li>Make all paths relative to the root and not execution start</li>
<li>Replace []error with chan error for faster? results</li>
<li>Add contexts with timeouts to all requests</li>
<li>Add output formatting option in cmd args</li>
<li>Glue paths together</li>
<li>metadata embedder</li>
<li>...and some other original spoton features</li>
</ul>
