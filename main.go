package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// Instrument defines an RSO instrument for matching in descriptions.
type Instrument struct {
	Name  string
	regex *regexp.Regexp
}

var instruments = []Instrument{
	Instrument{"Flute", regexp.MustCompile(`(?i)flute`)},
	Instrument{"Piccolo", regexp.MustCompile(`(?i)piccolo`)},
	Instrument{"Recorder", regexp.MustCompile(`(?i)recorder`)},
	Instrument{"Oboe", regexp.MustCompile(`(?i)oboe`)},
	Instrument{"English Horn", regexp.MustCompile(`(?i)english horn`)},
	Instrument{"Bassoon", regexp.MustCompile(`(?i)bassoon`)},
	Instrument{"Clarinet", regexp.MustCompile(`(?i)clarinet`)},
	Instrument{"Eb Clarinet", regexp.MustCompile(`(?i)e(b|-flat) clarinet`)},
	Instrument{"Bass Clarinet", regexp.MustCompile(`(?i)bass clarinet`)},
	Instrument{"Soprano Saxophone", regexp.MustCompile(`(?i)soprano sax`)},
	Instrument{"Alto Saxophone", regexp.MustCompile(`(?i)alto sax`)},
	Instrument{"Tenor Saxophone", regexp.MustCompile(`(?i)tenor sax`)},
	Instrument{"Baritone Saxophone", regexp.MustCompile(`(?i)baritone sax`)},
	Instrument{"Cornet", regexp.MustCompile(`(?i)cornet`)},
	Instrument{"Trumpet", regexp.MustCompile(`(?i)trumpet`)},
	Instrument{"Horn", regexp.MustCompile(`(?i)horn in`)},
	Instrument{"Trombone", regexp.MustCompile(`(?i)trombone`)},
	Instrument{"Tuba", regexp.MustCompile(`(?i)tuba`)},
	Instrument{"Euphonium", regexp.MustCompile(`(?i)euphonium`)},
	Instrument{"Violin", regexp.MustCompile(`(?i)violin`)},
	Instrument{"Viola", regexp.MustCompile(`(?i)viola`)},
	Instrument{"Cello", regexp.MustCompile(`(?i)cello`)},
	Instrument{"Double Bass", regexp.MustCompile(`(?i)double bass`)},
	Instrument{"Harp", regexp.MustCompile(`(?i)harp`)},
	Instrument{"Keyboard", regexp.MustCompile(`(?im)(keyboard|piano$)`)},
	Instrument{"Percussion", regexp.MustCompile(`(?i)(percussion|drum|triangle|cymbal)`)},
	Instrument{"Timpani", regexp.MustCompile(`(?i)timpani`)},
}

// findInstruments returns all instruments that the project with the given description needs.
func findInstruments(text string) []string {
	// TODO: First, match Markdown list item (^\* )
	var result []string
	for _, instr := range instruments {
		if instr.regex.FindString(text) != "" {
			result = append(result, instr.Name)
		}
	}
	return result
}

// Example: The final date to submit is November 24th.
var deadlineRegex = regexp.MustCompile(`(?mi)^.*(?:final date|due date|due on|last day).*?(january|february|march|april|may|june|july|august|september|october|november|december)\s+(\d+)`)

// findDeadline finds the deadline from the text. It returns a zero timestamp
// if no deadline is found.
func findDeadline(text string, created int64) time.Time {
	var t time.Time
	m := deadlineRegex.FindStringSubmatch(text)
	if m == nil {
		return t
	}
	t, err := time.Parse("January 2", fmt.Sprintf("%s %s", m[1], m[2]))
	if err != nil {
		return t
	}
	// Set deadline year so that it is after the post creation.
	ctime := time.Unix(created, 0)
	t = t.AddDate(ctime.Year(), 0, 0)
	if t.Month() < ctime.Month() {
		t = t.AddDate(1, 0, 0)
	}
	return t
}

// TODO: Newest published videos from
// https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=PLAl3fvW4KndjZDMFIs7w-f6Cm7Bp49gPA&maxResults=50&key=[YOUR_API_KEY]'
// Use nextPageToken to walk to last page

func printProjects(search *RedditSearch) {
	for i, project := range search.Data.Children {
		fmt.Printf("Project %d: %s\n%s\n", i, project.Data.Title, project.Data.URL)
		fmt.Printf("Due: %s\n", findDeadline(project.Data.Selftext, int64(project.Data.CreatedUtc)).Format("2006-01-02"))
		for _, instr := range findInstruments(project.Data.Selftext) {
			fmt.Printf(" - %s\n", instr)
		}
	}
}

func printGanttChartData(search *RedditSearch) {
	for _, project := range search.Data.Children {
		title := project.Data.Title
		created := time.Unix(int64(project.Data.CreatedUtc), 0).Format("2006-01-02")
		deadline := findDeadline(project.Data.Selftext, int64(project.Data.CreatedUtc))
		if !deadline.IsZero() {
			fmt.Printf("%s\t%s\t%s\n", title, created, deadline.Format("2006-01-02"))
		}
	}
}

// https://www.reddit.com/r/TheRedditSymphony/search.json?restrict_sr=1&sort=new&q=flair:%22Approved%20Project%22&limit=100
// next page with &count=100&after=xyz

func main() {
	searchjson, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("couldn't read search.json: %s\n", err)
		return
	}
	var search RedditSearch
	if err := json.Unmarshal(searchjson, &search); err != nil {
		fmt.Printf("couldn't parse json: %s\n", err)
		return
	}

	//printProjects(&search)
	printGanttChartData(&search)
}
