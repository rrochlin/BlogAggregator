package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rrochlin/BlogAggregator/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Missing Arguments for %v\n", cmd.name)
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Login Failed: %v", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("User Lookup Failed: %v", err)
	}
	for _, user := range users {
		fmt.Printf("* %v", user.Name)
		if s.cfg.CurrentUserName == user.Name {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Missing Name for %v\n", cmd.name)
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, cmd.args[0])
	if err == nil {
		return fmt.Errorf("User Already Exists\n")
	}
	user, err = s.db.CreateUser(
		ctx,
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return err
	}
	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}
	fmt.Printf("New User Created\n")
	fmt.Printf("ID: %v, created\n", user.ID)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)
	fmt.Printf("Name: %v\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	return s.db.TruncateTable(context.Background())
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg takes 1 argument: time between requests")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	newFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	multipleUrlErr := "pq: duplicate key value violates unique constraint \"posts_url_key\""
	for _, item := range newFeed.Channel.Item {
		publishDate := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishDate = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: publishDate,
				FeedID:      feed.ID,
			})
		if err != nil && err.Error() != multipleUrlErr {
			return err
		}

	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error

	if len(cmd.args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("post.Title: %v\n", post.Title)
		fmt.Printf("post.Description: %v\n", post.Description)
		fmt.Printf("post.PublishedAt: %v\n", post.PublishedAt)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Add feed needs 2 arguments, Name of feed, and URL")
	}
	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Created Feed: %v\n", feed)

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	return nil

}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByUID(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		fmt.Printf("Created By: %v\n", user.Name)
	}
	return nil

}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow takes 1 argument: url\nProvided: %v", cmd.args)
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v is now following %v\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following takes 0 arguments")
	}
	ctx := context.Background()

	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%v follows:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf(" * %v\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unfollow takes 1 argument: Feed URL")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	_, err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
