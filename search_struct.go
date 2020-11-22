package main

// RedditSearch is the schema of reddit's search endpoint.
type RedditSearch struct {
	Data struct {
		After    interface{} `json:"after"`
		Before   interface{} `json:"before"`
		Children []struct {
			Data struct {
				AllAwardings []struct {
					AwardSubType                     string      `json:"award_sub_type"`
					AwardType                        string      `json:"award_type"`
					AwardingsRequiredToGrantBenefits int64       `json:"awardings_required_to_grant_benefits"`
					CoinPrice                        int64       `json:"coin_price"`
					CoinReward                       int64       `json:"coin_reward"`
					Count                            int64       `json:"count"`
					DaysOfDripExtension              int64       `json:"days_of_drip_extension"`
					DaysOfPremium                    int64       `json:"days_of_premium"`
					Description                      string      `json:"description"`
					EndDate                          interface{} `json:"end_date"`
					GiverCoinReward                  int64       `json:"giver_coin_reward"`
					IconFormat                       string      `json:"icon_format"`
					IconHeight                       int64       `json:"icon_height"`
					IconURL                          string      `json:"icon_url"`
					IconWidth                        int64       `json:"icon_width"`
					ID                               string      `json:"id"`
					IsEnabled                        bool        `json:"is_enabled"`
					IsNew                            bool        `json:"is_new"`
					Name                             string      `json:"name"`
					PennyDonate                      int64       `json:"penny_donate"`
					PennyPrice                       int64       `json:"penny_price"`
					ResizedIcons                     []struct {
						Height int64  `json:"height"`
						URL    string `json:"url"`
						Width  int64  `json:"width"`
					} `json:"resized_icons"`
					ResizedStaticIcons []struct {
						Height int64  `json:"height"`
						URL    string `json:"url"`
						Width  int64  `json:"width"`
					} `json:"resized_static_icons"`
					StartDate                interface{} `json:"start_date"`
					StaticIconHeight         int64       `json:"static_icon_height"`
					StaticIconURL            string      `json:"static_icon_url"`
					StaticIconWidth          int64       `json:"static_icon_width"`
					SubredditCoinReward      int64       `json:"subreddit_coin_reward"`
					SubredditID              interface{} `json:"subreddit_id"`
					TiersByRequiredAwardings struct {
						Zero struct {
							AwardingsRequired int64 `json:"awardings_required"`
							Icon              struct {
								Format string `json:"format"`
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"icon"`
							ResizedIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_icons"`
							ResizedStaticIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_static_icons"`
							StaticIcon struct {
								Format interface{} `json:"format"`
								Height int64       `json:"height"`
								URL    string      `json:"url"`
								Width  int64       `json:"width"`
							} `json:"static_icon"`
						} `json:"0"`
						Three struct {
							AwardingsRequired int64 `json:"awardings_required"`
							Icon              struct {
								Format string `json:"format"`
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"icon"`
							ResizedIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_icons"`
							ResizedStaticIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_static_icons"`
							StaticIcon struct {
								Format interface{} `json:"format"`
								Height int64       `json:"height"`
								URL    string      `json:"url"`
								Width  int64       `json:"width"`
							} `json:"static_icon"`
						} `json:"3"`
						Six struct {
							AwardingsRequired int64 `json:"awardings_required"`
							Icon              struct {
								Format string `json:"format"`
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"icon"`
							ResizedIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_icons"`
							ResizedStaticIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_static_icons"`
							StaticIcon struct {
								Format interface{} `json:"format"`
								Height int64       `json:"height"`
								URL    string      `json:"url"`
								Width  int64       `json:"width"`
							} `json:"static_icon"`
						} `json:"6"`
						Nine struct {
							AwardingsRequired int64 `json:"awardings_required"`
							Icon              struct {
								Format string `json:"format"`
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"icon"`
							ResizedIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_icons"`
							ResizedStaticIcons []struct {
								Height int64  `json:"height"`
								URL    string `json:"url"`
								Width  int64  `json:"width"`
							} `json:"resized_static_icons"`
							StaticIcon struct {
								Format interface{} `json:"format"`
								Height int64       `json:"height"`
								URL    string      `json:"url"`
								Width  int64       `json:"width"`
							} `json:"static_icon"`
						} `json:"9"`
					} `json:"tiers_by_required_awardings"`
				} `json:"all_awardings"`
				AllowLiveComments          bool        `json:"allow_live_comments"`
				ApprovedAtUtc              interface{} `json:"approved_at_utc"`
				ApprovedBy                 interface{} `json:"approved_by"`
				Archived                   bool        `json:"archived"`
				Author                     string      `json:"author"`
				AuthorFlairBackgroundColor string      `json:"author_flair_background_color"`
				AuthorFlairCSSClass        string      `json:"author_flair_css_class"`
				AuthorFlairRichtext        []struct {
					A string `json:"a"`
					E string `json:"e"`
					T string `json:"t"`
					U string `json:"u"`
				} `json:"author_flair_richtext"`
				AuthorFlairTemplateID string        `json:"author_flair_template_id"`
				AuthorFlairText       string        `json:"author_flair_text"`
				AuthorFlairTextColor  string        `json:"author_flair_text_color"`
				AuthorFlairType       string        `json:"author_flair_type"`
				AuthorFullname        string        `json:"author_fullname"`
				AuthorPatreonFlair    bool          `json:"author_patreon_flair"`
				AuthorPremium         bool          `json:"author_premium"`
				Awarders              []interface{} `json:"awarders"`
				BannedAtUtc           interface{}   `json:"banned_at_utc"`
				BannedBy              interface{}   `json:"banned_by"`
				CanGild               bool          `json:"can_gild"`
				CanModPost            bool          `json:"can_mod_post"`
				Category              interface{}   `json:"category"`
				Clicked               bool          `json:"clicked"`
				ContentCategories     interface{}   `json:"content_categories"`
				ContestMode           bool          `json:"contest_mode"`
				Created               float64       `json:"created"`
				CreatedUtc            float64       `json:"created_utc"`
				DiscussionType        interface{}   `json:"discussion_type"`
				Distinguished         interface{}   `json:"distinguished"`
				Domain                string        `json:"domain"`
				Downs                 int64         `json:"downs"`
				Edited                interface{}   `json:"edited"`
				Gilded                int64         `json:"gilded"`
				Gildings              struct {
					Gid1 int64 `json:"gid_1"`
				} `json:"gildings"`
				Hidden                   bool   `json:"hidden"`
				HideScore                bool   `json:"hide_score"`
				ID                       string `json:"id"`
				IsCrosspostable          bool   `json:"is_crosspostable"`
				IsMeta                   bool   `json:"is_meta"`
				IsOriginalContent        bool   `json:"is_original_content"`
				IsRedditMediaDomain      bool   `json:"is_reddit_media_domain"`
				IsRobotIndexable         bool   `json:"is_robot_indexable"`
				IsSelf                   bool   `json:"is_self"`
				IsVideo                  bool   `json:"is_video"`
				Likes                    bool   `json:"likes"`
				LinkFlairBackgroundColor string `json:"link_flair_background_color"`
				LinkFlairCSSClass        string `json:"link_flair_css_class"`
				LinkFlairRichtext        []struct {
					E string `json:"e"`
					T string `json:"t"`
				} `json:"link_flair_richtext"`
				LinkFlairTemplateID   string        `json:"link_flair_template_id"`
				LinkFlairText         string        `json:"link_flair_text"`
				LinkFlairTextColor    string        `json:"link_flair_text_color"`
				LinkFlairType         string        `json:"link_flair_type"`
				Locked                bool          `json:"locked"`
				Media                 interface{}   `json:"media"`
				MediaEmbed            struct{}      `json:"media_embed"`
				MediaOnly             bool          `json:"media_only"`
				ModNote               interface{}   `json:"mod_note"`
				ModReasonBy           interface{}   `json:"mod_reason_by"`
				ModReasonTitle        interface{}   `json:"mod_reason_title"`
				ModReports            []interface{} `json:"mod_reports"`
				Name                  string        `json:"name"`
				NoFollow              bool          `json:"no_follow"`
				NumComments           int64         `json:"num_comments"`
				NumCrossposts         int64         `json:"num_crossposts"`
				NumReports            interface{}   `json:"num_reports"`
				Over18                bool          `json:"over_18"`
				ParentWhitelistStatus string        `json:"parent_whitelist_status"`
				Permalink             string        `json:"permalink"`
				Pinned                bool          `json:"pinned"`
				PostHint              string        `json:"post_hint"`
				Preview               struct {
					Enabled bool `json:"enabled"`
					Images  []struct {
						ID          string `json:"id"`
						Resolutions []struct {
							Height int64  `json:"height"`
							URL    string `json:"url"`
							Width  int64  `json:"width"`
						} `json:"resolutions"`
						Source struct {
							Height int64  `json:"height"`
							URL    string `json:"url"`
							Width  int64  `json:"width"`
						} `json:"source"`
						Variants struct{} `json:"variants"`
					} `json:"images"`
				} `json:"preview"`
				Pwls                  int64         `json:"pwls"`
				Quarantine            bool          `json:"quarantine"`
				RemovalReason         interface{}   `json:"removal_reason"`
				RemovedBy             interface{}   `json:"removed_by"`
				RemovedByCategory     interface{}   `json:"removed_by_category"`
				ReportReasons         interface{}   `json:"report_reasons"`
				Saved                 bool          `json:"saved"`
				Score                 int64         `json:"score"`
				SecureMedia           interface{}   `json:"secure_media"`
				SecureMediaEmbed      struct{}      `json:"secure_media_embed"`
				Selftext              string        `json:"selftext"`
				SelftextHTML          string        `json:"selftext_html"`
				SendReplies           bool          `json:"send_replies"`
				Spoiler               bool          `json:"spoiler"`
				Stickied              bool          `json:"stickied"`
				Subreddit             string        `json:"subreddit"`
				SubredditID           string        `json:"subreddit_id"`
				SubredditNamePrefixed string        `json:"subreddit_name_prefixed"`
				SubredditSubscribers  int64         `json:"subreddit_subscribers"`
				SubredditType         string        `json:"subreddit_type"`
				SuggestedSort         string        `json:"suggested_sort"`
				Thumbnail             string        `json:"thumbnail"`
				ThumbnailHeight       interface{}   `json:"thumbnail_height"`
				ThumbnailWidth        interface{}   `json:"thumbnail_width"`
				Title                 string        `json:"title"`
				TopAwardedType        interface{}   `json:"top_awarded_type"`
				TotalAwardsReceived   int64         `json:"total_awards_received"`
				TreatmentTags         []interface{} `json:"treatment_tags"`
				Ups                   int64         `json:"ups"`
				UpvoteRatio           float64       `json:"upvote_ratio"`
				URL                   string        `json:"url"`
				UserReports           []interface{} `json:"user_reports"`
				ViewCount             interface{}   `json:"view_count"`
				Visited               bool          `json:"visited"`
				WhitelistStatus       string        `json:"whitelist_status"`
				Wls                   int64         `json:"wls"`
			} `json:"data"`
			Kind string `json:"kind"`
		} `json:"children"`
		Dist    int64    `json:"dist"`
		Facets  struct{} `json:"facets"`
		Modhash string   `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"`
}
