RSO Website
===========

This is the code for the RSO website, found at https://www.rso-music.com/

The site shows data from Reddit, YouTube and Google Sheets that it fetches
periodically to generate static pages.


Setup
-----

You need Go installed. Run `go build` to build the binary `rso-projects`.

To run, you need API keys for Reddit and YouTube.

### Reddit

Go to https://old.reddit.com/prefs/apps/ and create a "personal use script".
Create a file named `agentfile` with the following contents:

```
user_agent: "graw:rso-projects:0.1 (by /u/<your reddit name>)"
client_id: "<put your client id>"
client_secret: "<put your client secret>"
username: ""
password: ""
```

Note that username and password stay empty, we only fetch public data.

### YouTube

Go to the Google Cloud Console to create a YouTube API key: https://console.cloud.google.com/apis/api/youtube.googleapis.com/credentials

Put your API key in the `YOUTUBE_API_KEY` environment variable.


Running
-------

You need to run the binary from the root of the repository, for example:

```
YOUTUBE_API_KEY=<your key> ./rso-projects
```

It will fetch data from Reddit and YouTube and will render `template.html` to
`static/index.html`. All fetched data is also stored in JSON files in `data/`
for debugging and for development. After changing the template, run
`./rso-projects -cached` to render from these data files (fast!) instead of
re-fetching everyting.

Set up your web server to serve from `static/`.


How does it work?
-----------------

Data fetching happens in `data.go`. We fetch:

- Recent posts with flairs "Official", "Official Project" and "Approved
  Project" via a Reddit search.
- The latest three weekly update threads, including comments, also via a
  Reddit search.
- All videos from the ["RSO All Playlist"](https://www.youtube.com/playlist?list=PLAl3fvW4KndiZAQtPmFCUFD6nImDC89Gv) on YouTube.

In `projects.go`, there is code to find information from this data, including:
- project start (= date of Reddit post)
- project deadline (by searching the post for certain keywords and a date)
- instruments that may be submitted for a project, by matching from a
- pre-defined list of instruments
  latest update from the weekly update thread, by looking for a comment that
  contains a link to the original project post
- the released video (for finished projects), by comparing longest common substrings

All this information is compiled in `htmlpage.go` and provided to
`template.html` (via [Go templating][gotmpl]), which results in
`static/index.html`. For the stats page (implemented in JavaScript), all this
data is also written to `static/projects.json`.

[gotmpl]: https://golang.org/pkg/text/template/

The stats page uses a CSV export of the [all projects sheet][allpr] to show
older projects as well. This file is currently not downloaded automatically.
Fetch it from [here][csv], remove the rows before the table (so that the
header row is the first line of the file), and save to
`static/allprojects.csv`.

[allpr]: https://docs.google.com/spreadsheets/d/12njIGc2_G4uMJ8uvfq1uKvdFfzopRYdhCdRdfo3e7Hg/edit?usp=sharing
[csv]: https://docs.google.com/spreadsheets/d/12njIGc2_G4uMJ8uvfq1uKvdFfzopRYdhCdRdfo3e7Hg/gviz/tq?tqx=out:csv
