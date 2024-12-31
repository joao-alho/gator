# Gator CLI

A simple CLI RSS feed.

## Requirements
- go v1.23.x
- A running postgres instance
- a `.gatorconfig.json` file on your `$HOME`

## Installation

```bash
go install github.com/joao-alho/gator
```

## .gatorconfig.json

The configuration for Gator should live in `$HOME/.gatorconfig.json`.
This is a JSON file with 2 options:
  - db_url: the URL for your postgres instance
  - current_user

## Usage:

Register a user:
```bash
gator register <username>
```

Add a feed to be tracked for the current user:
```bash
gator addfeed <feed name> <feed url>
```

Start collecting posts from tracked feeds, you must pass an 
interval in the format `30s, 60m, 2h`:
```bash
gator agg <interval>
```

The above command is a long living process that will fetch 
new posts from the tracked feeds every `interval`.


To browse collected posts use the `browse` command: 
```bash
gator browse <n>
```

where `n` is how many posts to show in the terminal.
This will look like:

```
- <Post Title>
  * URL: <URL>
  * Published at: <published timestamp>
```
