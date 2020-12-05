package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/turnage/graw/reddit"
)

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

var prometheusTemplate = template.Must(template.New("metrics.txt").Parse(`
# HELP rso_last_run_seconds Timestamp of the last run
# TYPE rso_last_run_seconds gauge
rso_last_run_seconds {{.LastRunTS}}

# HELP rso_data_items_total Number of items fetched from Reddit or YouTube
# TYPE rso_data_items_total gauge
rso_data_items_total{type="posts"} {{.PostsTotal}}
rso_data_items_total{type="weekly_updates"} {{.WeeklyUpdatesTotal}}
rso_data_items_total{type="videos"} {{.VideosTotal}}
`))

func writePrometheusStats(client *DataClient) error {
	f, err := os.Create("static/metrics.tmp")
	if err != nil {
		return fmt.Errorf("creating static/metrics.tmp failed: %w", err)
	}
	err = prometheusTemplate.Execute(f, map[string]interface{}{
		"LastRunTS":          time.Now().Unix(),
		"PostsTotal":         len(client.Posts),
		"WeeklyUpdatesTotal": len(client.WeeklyUpdates),
		"VideosTotal":        len(client.Videos),
	})
	f.Sync()
	f.Close()
	if err != nil {
		return err
	}
	return os.Rename("static/metrics.tmp", "static/metrics.txt")
}

// https://www.reddit.com/r/TheRedditSymphony/search.json?restrict_sr=1&sort=new&q=flair:%22Approved%20Project%22&limit=100
// next page with &count=100&after=xyz

var cachedFlag = flag.Bool("cached", false, "use cached data")

func main() {
	flag.Parse()

	client := NewDataClient()

	var err error

	if *cachedFlag {
		if err = client.LoadFromCache(); err != nil {
			fmt.Printf("couldn't load from cache: %s\n", err)
			return
		}
	} else {
		if err = client.Init(); err != nil {
			fmt.Printf("couldn't initialize data client: %s\n", err)
			return
		}
		if err = client.FetchPosts(); err != nil {
			fmt.Printf("couldn't fetch posts: %s\n", err)
			return
		}
		if err = client.FetchWeeklyUpdates(); err != nil {
			fmt.Printf("couldn't fetch weekly updates: %s\n", err)
			return
		}
		if err = client.FetchVideos(); err != nil {
			fmt.Printf("couldn't fetch videos: %s\n", err)
			return
		}
	}

	//printProjects(&search)
	//printGanttChartData(&search)
	createHTMLPage(client)

	if err = writePrometheusStats(client); err != nil {
		fmt.Printf("failed writing stats: %s\n", err)
	}
}
