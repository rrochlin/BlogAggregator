package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rrochlin/BlogAggregator/gator/internal/config"
	"github.com/rrochlin/BlogAggregator/gator/internal/database"
	"os"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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
