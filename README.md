# hearts-go

The card game [Hearts](https://bicyclecards.com/how-to-play/hearts/) implemented in [Go](https://go.dev/).

## Project Status

Nothing working yet. Building core game logic.

## Setup

To use the command-line interface (cli), no particular setup is required. It doesn't work yet but the idea is that eventually you can play a game against three bots using just the cli binary.

To host games as a web service, you'll need a database and a front-end (see parappayo/hearts-deno). I'd like to support MongoDB as an option but for now the work in progress is a Postgres table.

### Postgres Setup

You must define the environment variable `DB_CONN` to a Postgres connection string in order for the server to start succesfully. The server will run `db/schema.sql` on start in an attempt to create the necessary database schema.

It's a good idea to create a new user and database specifically for the hearts server, in order to keep it separate from anything else that might be using your database. In a typical Linux setup, first connect to Postgres as an admin user:

```
sudo -i -u postgres
psql
```

Then run the following queries:

```
CREATE USER hearts_user WITH PASSWORD 'secure_pass';
CREATE DATABASE hearts WITH OWNER hearts_user;
```

Now disconnect from the db (ctrl-d), log out from the postgres user (ctrl-d), and you use an evn var for local development:

```
export DB_CONN=postgresql://hearts_user:secure_pass@localhost:5432/hearts
```

## Usage

Run the test suite: `make test`

Run the command-line interface: `make run`

Start the web API: `make serve`

## Goals

The goal is to strengthen my knowledge of Go. I'd also like to end up with a suite of unit tests that could be used as a starting point for implementing Hearts in other programming languages, and a sophisticated enough version of the game that I can play it for fun.

### Progress

Some core game logic has been implemented, including scoring rounds. The test suite needs to be expanded.

A simple http server has been added to serve game state. This is to facilitate progress on a frontend since I also want to learn Deno.

Work is needed to persist game state. Ideally I'd like to do some event sourcing and have the option for the database to be either a Postgres table or a MongoDB instance.

### Bot AIs

Developing a sophisticated AI for playing Hearts is not a priority, but it would be a shame not to put some effort into having bots that are fun to play against (at least at a novice level.)

In increasing levels of sophistication, I may attempt the following:

- Chaos Bot - plays a random, legal move each turn
- Always Low Bot - plays the lowest card available, plays points cards when they have a void
- Naive Greedy Bot - tries to surrender control while discarding the highest cards they can, uses table history to recognize which suits are safer to lead
- Greedy Bot - work in some probability methods to evaluate plays, recognize when a shoot the moon is likely and be able to play in that situation
- [A* Search](https://en.wikipedia.org/wiki/A*_search_algorithm) Bot - (very unlikely to get this far), classical AI methods: use probability to score the best move on each turn with look-ahead search for a few rounds (branching factor too high?), with heuristics to shortcut evaluation of common scenarios and ideally an opening book
- Neural Net AI Bot - (very unlikely to get this far), modern AI methods: train a neural network on a large history of Hearts games and be able to identify strong moves
