# Bounce 'n Junk

This is a crappy clone of the glorius [Jump 'n Bump](https://en.wikipedia.org/wiki/Jump_%27n_Bump) DOS game

**The goal** of the project is to create a Jump 'n Bump like online game, which can be played in the browser
- The server is written in Go using the Echo web-framework and Gorilla's WebSocket implementation
- The command line interface is using spf13's Cobra + Viper framework
- The database is a tiny wrapper around BoltDB
- The client is a vanilla HTML CSS JavaScript application

## Project State
POC

## Setup
The whole project can be compiled into a single binary with Go and carryed around. The server maintains persistent data in a database file (the path of the database file can be configured or completly omitted). The behaviour of the server can be configured with the **Configuration File**, with environment variables and with command line arguments.

Once booted up, the server listens on a tcp port, serves the game site, upgrades the incoming WebSocket requests and connects them to the game. The game is constantly going on, players can join and disconnect anytime. The game consist of rounds. Before every round starts the previous one's
The project contains a **Dockerfile** which can build a small Docker image on top of Alipne Linux

## Commands
commands...
