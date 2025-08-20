# superfight

An online game based on Pipework Studio's [SUPERFIGHT](https://store.steampowered.com/app/404770/SUPERFIGHT/), which was sadly discontinued.

## Prerequisites

- [Go](https://go.dev/doc/install)
- [Yarn](https://yarnpkg.com/getting-started/install)

## Building

Card data is not provided in this repository due to copyright concerns. In order to build the project, you will need to create your own deck. First, copy the example files:

```bash
cp cards/black.txt.example cards/black.txt
cp cards/white.txt.example cards/white.txt
```

Then populate the files with your own cards, one per line.

To build the client:

```bash
cp client
yarn
yarn build
```

To build the server:

```bash
go build
```

## Running

To run the game, first build the client per above, then run:

```bash
go run .
```
