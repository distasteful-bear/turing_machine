# Turing Machine

A web version of the Turing Machine logic puzzle game, built in Go.

Play it here: https://turing-machine.james-metz.com

## Highlights

- Generates solvable random puzzles by combining verifier rules until the solution is uniquely identifiable.
- Provides both a Gin-based web application and a local CLI game mode.
- Uses Auth0 login, cookie sessions, and server-side puzzle storage to keep active games isolated from the client.
- Records completed games in Firestore and computes leaderboard rankings from wins and guess efficiency.
- Includes focused Go tests around API and puzzle setup behavior.

## Tech Stack

Go, Gin, Auth0, Firestore, Tailwind CSS, Docker.

## Running Locally

```sh
./scripts/run.sh
```

For CLI mode:

```sh
go run . --local
```

Run tests:

```sh
go test ./...
```
