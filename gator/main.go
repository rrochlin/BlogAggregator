package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rrochlin/BlogAggregator/gator/internal/config"
	"github.com/rrochlin/BlogAggregator/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config")
		fmt.Println(err)
		return
	}
	db, err := sql.Open("postgres", cfg.DBUrl)
	s := state{
		cfg: &cfg,
		db:  database.New(db),
	}
	cmds := commands{m: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	if len(os.Args) < 2 {
		fmt.Printf("Not enough arguments: %v/2\n", len(os.Args))
		os.Exit(1)
	}
	cmdName := os.Args[1]
	args := os.Args[2:]
	fmt.Printf("Running %v, with args %v\n", cmdName, args)
	err = cmds.run(&s, command{name: cmdName, args: args})
	if err != nil {
		fmt.Printf("Encountered error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)

}

type commands struct {
	m map[string]func(*state, command) error
}

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

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

func (c *commands) register(name string, f func(*state, command) error) {
	c.m[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.m[cmd.name]
	if !ok {
		return fmt.Errorf("Command %v not found\n", cmd.name)
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}
