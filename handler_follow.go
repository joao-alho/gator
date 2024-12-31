package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joao-alho/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Insufficient number of arguments")
	}
	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to retrieve feed info: %w", err)
	}
	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Printf("%s follows feed %s\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Insufficient number of arguments")
	}
	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	}); err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Wrong number of arguments")
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve following info: %w", err)
	}
	if len(follows) > 0 {
		fmt.Println("Following:")
		for _, follow := range follows {
			fmt.Printf("  - %s\n", follow.FeedName)
		}
	} else {
		fmt.Println("Following nothing.")
	}
	return nil
}
