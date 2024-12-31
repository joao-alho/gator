package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joao-alho/gator/internal/database"
)

const (
	bootDevLongForm = "Wed, 03 Jul 2019 00:00:00"
)

func handlerBrowse(s *state, cmd command) error {
	var limit int32
	if len(cmd.Args) == 1 {
		arg_int, err := strconv.ParseInt(cmd.Args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse limit argument: %w", err)
		}
		limit = int32(arg_int)
	} else {
		limit = 2
	}
	posts, err := s.db.GetPosts(context.Background(), limit)
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}
	for _, p := range posts {
		if err := printPost(&p); err != nil {
			return fmt.Errorf("failed to print post: %w", err)
		}
		if err := s.db.MarkPostFetched(context.Background(), database.MarkPostFetchedParams{
			ID:            p.ID,
			UpdatedAt:     time.Now(),
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to mark as fetched: %w", err)
		}
	}

	return nil
}

func printPost(p *database.Post) error {
	fmt.Printf("- %s\n", p.Title)
	fmt.Printf("  * URL: %s\n", p.Url)
	fmt.Printf("  * Published at: %s\n", p.PublishedAt.Time.String())
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Insufficient arguments")
	}
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", time_between_reqs.String())
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	fmt.Printf("Fetching %s feed.\n", next_feed.Name)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	if err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            next_feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
	}); err != nil {
		return fmt.Errorf("failed to mark feed: %w\n", err)
	}

	rss_feed, err := fetchFeed(context.Background(), next_feed.Url)

	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	for _, item := range rss_feed.Channel.Item {
		if item.Title != "" {
			var published_at sql.NullTime
			pub_date, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				log.Printf("Failed to parse published date: %s\n", item.PubDate)
			}
			if err := published_at.Scan(pub_date); err != nil {
				fmt.Printf("did not push: %s\n", err.Error())
			}

			_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description},
				PublishedAt: published_at,
				FeedID:      next_feed.ID,
			})
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("failed to create post: %w", err)
			}
		}
	}
	return nil
}
