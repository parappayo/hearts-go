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

## Golang Setup

You can use `go install` to manage the version of go.

```
go install golang.org/dl/go1.25.3@latest
~/go/bin/go1.25.3 download
~/go/bin/go1.25.3 run server/hearts.go
```

Using the above method won't work with the project's `Makefile`, but to fix that you can also update your environment to use a specific Go version:

```
export GOROOT='/home/your_user/sdk/go1.25.3'
export PATH=$GOROOT/bin:$PATH
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

### Event Sourcing

This project is being used as an example to put [event sourcing](https://martinfowler.com/eaaDev/EventSourcing.html) into practice. A key benefit of using event sourcing here is to be able to easily implement match replays, such that users can revisit completed Hearts matches and see the entire sequence of play.

Each event has an aggregate ID, which in this case is also the match ID used to identify instances of games of Hearts. Under this model, each match can be thought of as a domain object with an ID (the match ID) and a revision number that counts up incrementally. The current state of the match at any given revision number can always be reconstructed by replaying the events up to that point.

Events are not exposed directly through the service API. Endpoints are implemented in such a way as to expose matches as domain objects that are acted on through GET and POST operations. State changes are implemented as events being issued in the backend but this is hidden from the API client.

The following events are planned:

* `match-created` - In order for a match to exist, it must first be created. The match ID is assigned by this event.
* `player-joined` - Players may only join a match after it is created and before it has started. A maximum of four players are allowed to have joined a match at one time.
* `player-left` - Players may leave a match after it is created and before it has started.
* `match-started` - A match that has exactly four players joined-up may start. After a match has started, players may not join or leave, and may only play cards (on their turn, of course.) This event deals out the initial hands for the players.
* `card-played` - The player whose turn it currently is may play a card as long as it is a valid Hearts play.
* `round-finished` - When a card is played and all of the players now have empty hands, a round has concluded. If none of the players has passed the 100 point threshold, then the backend automatically issues this event to deal out a new hand and to snapshot the current state of scoring.
* `match-finished` - When a round has concluded and one of the players has passed the 100 point threshold, the game has concluded and no further events are accepted. The final scores are snapshotted in this event.

The following events may be implemented later:

* `player-profile-updated` - For players to assign themselves names, etc.

It's a bit of premature optimization, but this event model allows for conveniently querying the set of recently completed games by looking for `match-finished` events in a given time span. It's also easy to query for games in progress by looking for match IDs which have a `match-created` event but no `match-finished` event.

### Ports & Adapters

The Ports & Adapters system design is also known as [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)). I haven't thought ahead much about this, but an important thing to keep in mind is that core domain objects (Hearts players, cards, matches / tables), API objects (requests and responses), and persistence records (database rows and documents) should each live in their separate layers and interact only in specific modules of code -- the adapters, which create domain objects out of interface data and vice-versa.

When done well, Ports & Adapters architecture de-couples core domain objects from implementation concerns such that an http API could be substituted for a gRPC API, or a Postgres database could be substituted for a MongoDB database, without changes rippling through the entire system. It also facilitates writing cleaner, more meaningful tests, such as properly isolated unit tests and efficient, reliable integration tests.

### Bot AIs

Developing a sophisticated AI for playing Hearts is not a priority, but it would be a shame not to put some effort into having bots that are fun to play against (at least at a novice level.)

In increasing levels of sophistication, I may attempt the following:

- Chaos Bot - plays a random, legal move each turn
- Always Low Bot - plays the lowest card available, plays points cards when they have a void
- Naive Greedy Bot - tries to surrender control while discarding the highest cards they can, uses table history to recognize which suits are safer to lead
- Greedy Bot - work in some probability methods to evaluate plays, recognize when a shoot the moon is likely and be able to play in that situation
- [A* Search](https://en.wikipedia.org/wiki/A*_search_algorithm) Bot - (very unlikely to get this far), classical AI methods: use probability to score the best move on each turn with look-ahead search for a few rounds (branching factor too high?), with heuristics to shortcut evaluation of common scenarios and ideally an opening book
- Neural Net AI Bot - (very unlikely to get this far), modern AI methods: train a neural network on a large history of Hearts games and be able to identify strong moves
