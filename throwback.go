package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
	"google.golang.org/api/youtube/v3"
)

// postThrowback posts a random video as "Thursday Throwback".
func postThrowback(client *DataClient) error {
	bot, err := reddit.NewBotFromAgentFile("agentfile", 1*time.Second)
	if err != nil {
		return fmt.Errorf("creating reddit bot failed: %w", err)
	}

	posts, err := fetchPreviousThrowbackPosts(bot)
	if err != nil {
		return err
	}

	prevPostTime := time.Unix(int64(posts[0].CreatedUTC), 0)
	if prevPostTime.AddDate(0, 0, 6).After(time.Now()) {
		return fmt.Errorf("previous throwback post too young (%s)", prevPostTime.Format(time.RFC3339))
	}

	video, err := chooseThrowbackVideo(client.Videos, posts)
	if err != nil {
		return fmt.Errorf("couldn't choose video: %w", err)
	}

	title := fmt.Sprintf("Throwback %s! %s", time.Now().Format("Monday"), video.Snippet.Title)
	url := fmt.Sprintf("https://youtu.be/%s", video.ContentDetails.VideoId)
	fmt.Printf("%s <%s>\n", title, url)

	if err = bot.PostLink("TheRedditSymphony", title, url); err != nil {
		return fmt.Errorf("couldn't post throwback link: %w", err)
	}

	return nil
}

// fetchPreviousThrowbackPosts fetches previous posts flaired "Throwback Thursday".
func fetchPreviousThrowbackPosts(bot reddit.Scanner) ([]*reddit.Post, error) {
	results, err := bot.ListingWithParams("/r/TheRedditSymphony/search", map[string]string{
		"restrict_sr": "1",
		"sort":        "new",
		"limit":       "10",
		"q":           "flair:\"Throwback Thursday\"",
	})
	if err != nil {
		return nil, fmt.Errorf("fetching throwback posts failed: %w", err)
	}
	return results.Posts, nil
}

// chooseThrowbackVideo (intelligently) chooses a random throwback video.
func chooseThrowbackVideo(videos []youtube.PlaylistItem, posts []*reddit.Post) (*youtube.PlaylistItem, error) {
	var v youtube.PlaylistItem
CHOOSE:
	for i := 0; i < len(videos)/2; i++ {
		v = videos[rand.Intn(len(videos))]

		// No new videos published in the previous month
		publishedAt, err := time.Parse(time.RFC3339, v.ContentDetails.VideoPublishedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid published at for video %s: %w", v.Snippet.Title, err)
		}
		if publishedAt.After(time.Now().AddDate(0, -1, 0)) {
			continue
		}

		// No videos from the last 10 throwback posts.
		for _, post := range posts {
			if strings.Contains(post.URL, v.ContentDetails.VideoId) {
				continue CHOOSE
			}
		}
	}
	return &v, nil
}
