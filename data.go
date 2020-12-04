package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/turnage/graw/reddit"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var rsoPlaylistID = "PLAl3fvW4KndjZDMFIs7w-f6Cm7Bp49gPA"

// DataClient fetches RSO posts and comments from reddit and videos from  YouTube.
type DataClient struct {
	bot     reddit.Bot
	youtube *youtube.Service

	Posts         []reddit.Post
	WeeklyUpdates []reddit.Comment
	Videos        []youtube.PlaylistItemSnippet
}

// NewDataClient creates a new, unitialized client.
func NewDataClient() *DataClient {
	return &DataClient{}
}

// Init reads auth data from "agentfile" to initialize a reddit
// client and an API key from the YOUTUBE_API_KEY environment variable to
// create a YouTube client.
func (c *DataClient) Init() error {
	var err error
	c.bot, err = reddit.NewBotFromAgentFile("agentfile", 1*time.Second)
	if err != nil {
		return fmt.Errorf("creating reddit bot failed: %w", err)
	}

	c.youtube, err = youtube.NewService(context.TODO(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		return fmt.Errorf("creating YouTube client failed: %w", err)
	}
	return nil
}

// LoadFromCache populates posts and comments from data/*.json
func (c *DataClient) LoadFromCache() error {
	if err := loadFromCache("posts.json", &c.Posts); err != nil {
		return err
	}
	if err := loadFromCache("weekly_updates.json", &c.WeeklyUpdates); err != nil {
		return err
	}
	if err := loadFromCache("videos.json", &c.Videos); err != nil {
		return err
	}
	return nil
}

// FetchPosts fetches the latest Approved Projects, Official Projects and
// Official posts from Reddit.
func (c *DataClient) FetchPosts() error {
	results, err := c.bot.ListingWithParams("/r/TheRedditSymphony/search", map[string]string{
		"restrict_sr": "1",
		"sort":        "new",
		"limit":       "100",
		"q":           "flair:\"Approved Project\" OR flair:\"Official Project\" OR flair:\"Official\"",
	})
	if err != nil {
		return fmt.Errorf("fetching posts failed: %w", err)
	}

	// Copy posts to avoid pointers.
	posts := make([]reddit.Post, len(results.Posts))
	for i, post := range results.Posts {
		posts[i] = *post
	}

	c.Posts = posts

	return writeToCache("posts.json", posts)
}

// FetchWeeklyUpdates fetches the comments on the last weekly project update threads.
func (c *DataClient) FetchWeeklyUpdates() error {
	result, err := c.bot.ListingWithParams("/r/TheRedditSymphony/search", map[string]string{
		"restrict_sr": "1",
		"sort":        "new",
		"limit":       "3",
		"q":           "Weekly Project Update Thread author:AutoModerator",
	})
	if err != nil {
		return fmt.Errorf("fetching weekly update posts failed: %w", err)
	}

	var comments []reddit.Comment
	for _, post := range result.Posts {
		// Fetch comments.
		fullpost, err := c.bot.Thread(post.Permalink)
		if err != nil {
			return fmt.Errorf("fetching comments for %s failed: %w", post.Title, err)
		}
		for _, comment := range fullpost.Replies {
			c := *comment
			// We are only interested in top-level comments.
			c.Replies = nil
			comments = append(comments, c)
		}
	}

	c.WeeklyUpdates = comments

	return writeToCache("weekly_updates.json", comments)
}

// FetchVideos fetches the latest videos from YouTube.
func (c *DataClient) FetchVideos() error {
	var videos []youtube.PlaylistItemSnippet
	call := c.youtube.PlaylistItems.List([]string{"snippet"}).PlaylistId(rsoPlaylistID).MaxResults(50)
	err := call.Pages(context.TODO(), func(res *youtube.PlaylistItemListResponse) error {
		for _, item := range res.Items {
			videos = append(videos, *item.Snippet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(videos, func(i, j int) bool { return videos[i].PublishedAt < videos[j].PublishedAt })

	c.Videos = videos

	return writeToCache("videos.json", videos)
}

func writeToCache(name string, data interface{}) error {
	if err := os.MkdirAll("data", 0777); err != nil {
		return fmt.Errorf("couldn't create data directory: %w", err)
	}
	f, err := os.Create("data/" + name)
	if err != nil {
		return fmt.Errorf("couldn't create data/%s: %w", name, err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(data); err != nil {
		return fmt.Errorf("couldn't encode %s: %w", name, err)
	}
	return nil
}

func loadFromCache(name string, data interface{}) error {
	f, err := os.Open("data/" + name)
	if err != nil {
		return fmt.Errorf("couldn't load data/%s: %w", name, err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(data); err != nil {
		return fmt.Errorf("couldn't decode data/%s: %w", name, err)
	}
	return nil
}
