package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/turnage/graw/reddit"
)

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
	createHTMLPage(client)
}
