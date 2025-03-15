package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rrochlin/BlogAggregator/gator/internal/database"
	"time"
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
	// if len(cmd.args) == 0 {
	// 	return fmt.Errorf("Missing Name for %v\n", cmd.name)
	// }
	// feed, err := fethFeed(context.Background(), cmd.args[0])
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fethFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("feed: %v\n", feed)
	return nil
}
