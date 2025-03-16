# Gator

## Dependancies

- Postgres
- Go

## Installation

- from root directory run `go install`
- call gator with command `gator`

## Configuration

- gator relies on a `.gatorconfig.json` file

```json
{
    "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name":"kahya"
}
```

- postgres connection string will need sslmode disabled for local
- current_user_name is controlled via gator cli

## Running Gator

- gator is not secure, users are identified by unique names not password protected
- users can add RSSFeeds and follow/unfollow these feeds
- users can `browse` their feeds to see the latest posts
- commands are called with the syntax `gator {command} {args...}`

## Commands

- `register {user}` registers a new user with name {user}
- `login {user}` sets active user to {user}
- `reset` truncates all application data in postgres
- `users` displays registered users, emphasizing currently logged in user
- `agg {duration}` starts post aggregation for a set time
  - valid values for duration are like 1h, 1m, 1s, 4s, 3d...
- `addfeed {url}` adds a feed to the database by feed URL
- `feeds` prints out currently available feeds
- `follow {url}` follows a feed by url
- `following` prints out what feeds the current user is following
- `unfollow {url}` unfollows a feed by url
- `browse {optional: limit}` displays up to limit # of posts.
  - If limit is ommitted defaults to 2
