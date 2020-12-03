package main

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	"github.com/turnage/graw/reddit"
)

// Registers is a nicely-sorted list of instrument registers.
var Registers = []string{"Woodwinds", "Brass", "Strings", "Percussion", "Other"}

// Instrument defines an RSO instrument for matching in descriptions.
type Instrument struct {
	Register string
	Name     string
	regex    *regexp.Regexp
}

var instruments = []Instrument{
	Instrument{"Woodwinds", "Flute", regexp.MustCompile(`(?i)flute`)},
	Instrument{"Woodwinds", "Piccolo", regexp.MustCompile(`(?i)piccolo`)},
	Instrument{"Woodwinds", "Recorder", regexp.MustCompile(`(?i)recorder`)},
	Instrument{"Woodwinds", "Oboe", regexp.MustCompile(`(?i)oboe`)},
	Instrument{"Woodwinds", "English Horn", regexp.MustCompile(`(?i)english horn`)},
	Instrument{"Woodwinds", "Bassoon", regexp.MustCompile(`(?i)bassoon`)},
	Instrument{"Woodwinds", "Clarinet", regexp.MustCompile(`(?i)clarinet`)},
	Instrument{"Woodwinds", "Eb Clarinet", regexp.MustCompile(`(?i)e(b|-flat) clarinet`)},
	Instrument{"Woodwinds", "Bass Clarinet", regexp.MustCompile(`(?i)bass clarinet`)},
	Instrument{"Woodwinds", "Soprano Saxophone", regexp.MustCompile(`(?i)soprano sax`)},
	Instrument{"Woodwinds", "Alto Saxophone", regexp.MustCompile(`(?i)alto sax`)},
	Instrument{"Woodwinds", "Tenor Saxophone", regexp.MustCompile(`(?i)tenor sax`)},
	Instrument{"Woodwinds", "Baritone Saxophone", regexp.MustCompile(`(?i)baritone sax`)},
	Instrument{"Brass", "Cornet", regexp.MustCompile(`(?i)cornet`)},
	Instrument{"Brass", "Trumpet", regexp.MustCompile(`(?i)trumpet`)},
	Instrument{"Brass", "Horn", regexp.MustCompile(`(?i)horn in`)},
	Instrument{"Brass", "Trombone", regexp.MustCompile(`(?i)trombone`)},
	Instrument{"Brass", "Tuba", regexp.MustCompile(`(?i)tuba`)},
	Instrument{"Brass", "Euphonium", regexp.MustCompile(`(?i)euphonium`)},
	Instrument{"Strings", "Violin", regexp.MustCompile(`(?i)violin`)},
	Instrument{"Strings", "Viola", regexp.MustCompile(`(?i)viola`)},
	Instrument{"Strings", "Cello", regexp.MustCompile(`(?i)cello`)},
	Instrument{"Strings", "Double Bass", regexp.MustCompile(`(?i)double bass`)},
	Instrument{"Other", "Harp", regexp.MustCompile(`(?i)harp`)},
	Instrument{"Other", "Keyboard", regexp.MustCompile(`(?im)(keyboard|piano$)`)},
	Instrument{"Percussion", "Percussion", regexp.MustCompile(`(?i)(percussion|drum|triangle|cymbal)`)},
	Instrument{"Percussion", "Timpani", regexp.MustCompile(`(?i)timpani`)},
}

// findInstruments returns all instruments that the project with the given description needs.
func findInstruments(text string) []Instrument {
	// TODO: First, match Markdown list item (^\* )
	var result []Instrument
	for _, instr := range instruments {
		if instr.regex.FindString(text) != "" {
			result = append(result, instr)
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

// isProject identifies projects by their flair.
func isProject(post *reddit.Post) bool {
	return post.LinkFlairText == "Approved Project" || post.LinkFlairText == "Official Project"
}

func printProjects(posts []reddit.Post) {
	for i, project := range posts {
		if isProject(&project) {
			fmt.Printf("Project %d: %s\n%s\n", i, project.Title, project.URL)
			fmt.Printf("Due: %s\n", findDeadline(project.SelfText, int64(project.CreatedUTC)).Format("2006-01-02"))
			for _, instr := range findInstruments(project.SelfText) {
				fmt.Printf(" - %s\n", instr.Name)
			}
		}
	}
}

func printGanttChartData(posts []reddit.Post) {
	for _, project := range posts {
		if isProject(&project) {
			title := project.Title
			created := time.Unix(int64(project.CreatedUTC), 0).Format("2006-01-02")
			deadline := findDeadline(project.SelfText, int64(project.CreatedUTC))
			if !deadline.IsZero() {
				fmt.Printf("%s\t%s\t%s\n", title, created, deadline.Format("2006-01-02"))
			}
		}
	}
}

// https://www.reddit.com/r/TheRedditSymphony/search.json?restrict_sr=1&sort=new&q=flair:%22Approved%20Project%22&limit=100
// next page with &count=100&after=xyz

var cachedFlag = flag.Bool("cached", false, "use cached data")

func main() {
	flag.Parse()

	client, err := NewRedditClient()
	if err != nil {
		fmt.Printf("couldn't initialize reddit client: %s\n", err)
		return
	}

	if *cachedFlag {
		if err = client.LoadFromCache(); err != nil {
			fmt.Printf("couldn't load from cache: %s\n", err)
			return
		}
	} else {
		if err = client.FetchPosts(); err != nil {
			fmt.Printf("couldn't fetch posts: %s\n", err)
			return
		}
		if err = client.FetchWeeklyUpdates(); err != nil {
			fmt.Printf("couldn't fetch weekly updates: %s\n", err)
			return
		}
	}

	//printProjects(&search)
	//printGanttChartData(&search)
	createHTMLPage(client.Posts)
}
