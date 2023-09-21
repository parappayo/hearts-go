# hearts-go

The card game [Hearts](https://bicyclecards.com/how-to-play/hearts/) implemented in [Go](https://go.dev/).

## Project Status

Nothing working yet. Building core game logic.

## Goals

Generally speaking, the goal is to strengthen my knowledge of Go and have some practical example code to refer to. Ideally I'd like to end up with a suite of unit tests that could be used as a starting point for implementing Hearts in other programming languages, and a sophisticated enough version of the game that I can play it for fun.

### Project phases

1. Core game logic
2. Command-line interface with bots that play random moves
3. Event sourcing with Postgres
4. HTTP service backend
5. HTML+JS frontend

### Stretch Goals

Many more will occur to me as development progresses.

- Frontend tech improvements: TypeScript, React, Redux
- Stand-alone desktop version (Gtk?)
- Mobile client (React Native?)

### Bot AIs

Developing a sophisticated AI for playing Hearts is not a priority, but it would be a shame not to put some effort into having bots that are fun to play against (at least at a novice level.)

In increasing levels of sophistication, I may attempt the following:

- Chaos Bot - plays a random, legal move each turn
- Always Low Bot - plays the lowest card available, plays points cards when they have a void
- Naive Greedy Bot - tries to surrender control while discarding the highest cards they can, uses table history to recognize which suits are safer to lead
- Greey Bot - work in some probability methods to evaluate plays, recognize when a shoot the moon is likely and be able to play in that situation
- [A* Search](https://en.wikipedia.org/wiki/A*_search_algorithm) Bot - (very unlikely to get this far), classical AI methods: use probability to score the best move on each turn with look-ahead search for a few rounds (branching factor too high?), with heuristics to shortcut evaluation of common scenarios and ideally an opening book
- Neural Net AI Bot - (very unlikely to get this far), modern AI methods: train a neural network on a large history of Hearts games and be able to identify strong moves
