package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
	"google.golang.org/api/youtube/v3"
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
	Instrument{"Woodwinds", "Oboe", regexp.MustCompile(`(?i)oboe\b`)},
	Instrument{"Woodwinds", "English Horn", regexp.MustCompile(`(?i)english horn`)},
	Instrument{"Woodwinds", "Bassoon", regexp.MustCompile(`(?i)bassoon`)},
	Instrument{"Woodwinds", "Clarinet", regexp.MustCompile(`(?i)clarinet`)},
	Instrument{"Woodwinds", "Eb Clarinet", regexp.MustCompile(`(?i)e(b|-flat) clarinet`)},
	Instrument{"Woodwinds", "Bass Clarinet", regexp.MustCompile(`(?i)bass clarinet`)},
	Instrument{"Woodwinds", "Soprano Saxophone", regexp.MustCompile(`(?i)soprano sax`)},
	Instrument{"Woodwinds", "Alto Saxophone", regexp.MustCompile(`(?i)alto sax`)},
	Instrument{"Woodwinds", "Tenor Saxophone", regexp.MustCompile(`(?i)tenor sax`)},
	Instrument{"Woodwinds", "Baritone Saxophone", regexp.MustCompile(`(?i)bari(tone)? sax`)},
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

var openInstrumentationRegex = regexp.MustCompile(`(?i)\bopen instrumentation\b`)

// isOpenInstrumentation detects pieces with open instrumentation (i.e., every instrument can submit).
func isOpenInstrumentation(text string) bool {
	return openInstrumentationRegex.MatchString(text)
}

// Example: The final date to submit is November 24th.
var deadlineRegex = regexp.MustCompile(`(?mi)^.*(?:final date|due date|due on|last day).*?(jan(?:uary)?|feb(?:ruary)?|mar(?:ch)?|apr(?:il)?|may|june?|july?|aug(?:ust)?|sep(?:tember)?|oct(?:ober)?|nov(?:ember)?|dec(?:ember)?)\s+(\d+)`)

// findDeadline finds the deadline from the text. It returns a zero timestamp
// if no deadline is found.
func findDeadline(text string, created int64) time.Time {
	var t time.Time
	m := deadlineRegex.FindStringSubmatch(text)
	if m == nil {
		return t
	}
	format := "January 2"
	if len(m[1]) == 3 {
		format = "Jan 2"
	}
	t, err := time.Parse(format, fmt.Sprintf("%s %s", m[1], m[2]))
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

// findUpdateComment finds the latest update comment for the given project by matching author and URL.
func findUpdateComment(post *reddit.Post, updates []reddit.Comment) *reddit.Comment {
	for _, comment := range updates {
		if comment.Author == post.Author && strings.Contains(comment.Body, post.ID) {
			return &comment
		}
	}
	return nil
}

// ProjectTag is a tag attached to a project depending on keywords in the post.
type ProjectTag struct {
	Name  string
	regex *regexp.Regexp
}

var projectTags = []ProjectTag{
	ProjectTag{"beginner-friendly", regexp.MustCompile(`(?i)\*\*beginner-friendly\*\*`)},
}

// findProjectTags scans the markdown text for tags.
func findProjectTags(text string) []string {
	var result []string
	for _, tag := range projectTags {
		if tag.regex.FindString(text) != "" {
			result = append(result, tag.Name)
		}
	}
	return result
}

// Video matching support

var similarityBadwordRegex *regexp.Regexp = regexp.MustCompile(`(?i)(/?r/)?theredditsymphony|rso|community|project|performed by|composition|symphonic movement|orchestra`)

// withoutSimilarityBadwords removes common words in project titles.
func withoutSimilarityBadwords(s string) string {
	return similarityBadwordRegex.ReplaceAllString(s, "")
}

// lcs calulates the longest common substring of two strings.
// https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Longest_common_substring#Go
func lcs(s1 string, s2 string) string {
	var m = make([][]int, 1+len(s1))
	for i := 0; i < len(m); i++ {
		m[i] = make([]int, 1+len(s2))
	}
	longest := 0
	xLongest := 0
	for x := 1; x < 1+len(s1); x++ {
		for y := 1; y < 1+len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] > longest {
					longest = m[x][y]
					xLongest = x
				}
			}
		}
	}
	return s1[xLongest-longest : xLongest]
}

// findMatchingVideo finds the release video for a post.
func findMatchingVideo(post *reddit.Post, videos []youtube.PlaylistItem) *youtube.PlaylistItem {
	postTitle := withoutSimilarityBadwords(post.Title)
	var bestVideo int
	var lcstr string

	for i, v := range videos {
		// TODO: Filter out videos published before the project's deadline
		videoTitle := withoutSimilarityBadwords(v.Snippet.Title)
		l := lcs(videoTitle, postTitle)
		if len(l) > len(lcstr) {
			bestVideo = i
			lcstr = l
		}
	}
	if len(lcstr) > 10 {
		//fmt.Printf("  LCS(%d)(%s): %s\n", len(lcstr), lcstr, client.Videos[bestVideo].Title)
		return &videos[bestVideo]
	}
	return nil
}
