# spotgone
wdmvv/spot-on but in go because prev version sucks!
<br>
<br><br>
grab ytdlp binary here - https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#installation<br>
<br>

<br>
honestly i'd love if someone helped me to figure out certain questions, for instance: i have 2 structures, playlist and album - is it better to make them as similar as possible<br>
or even transform to secret random third type to write methods on that third type? or is it ok if i just write methods on these separate structures even if they will be almost same...<br>
the problem is that does it even matter whether it is an album or playlist if in the end i can make album look like playlist<br>
_probably_ ill just go with this option - album is converted to the playlist under the hood and then i write methods on playlist and sometimes link one to another<br>
<br><br>
TODO:<br>
<ul>
<li>Make all paths relative to the root and not execution start</li>
<li>Replace []error with chan error for faster? results</li>
<li>Add contexts with timeouts to all requests</li>
<li>Add output formatting option in cmd args</li>
<li>Glue paths together</li>
<li>metadata inserter</li>
<li>...and some other original spoton features</li>
</ul>
