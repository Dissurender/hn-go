[![](https://img.shields.io/github/issues/dissurender/hn-go)](https://github.com/Dissurender/hn-go/issues) [![](https://img.shields.io/github/license/dissurender/hn-go)](https://github.com/Dissurender/hn-go/blob/main/LICENSE) ![](https://img.shields.io/github/languages/top/dissurender/hn-go)

# HN-Go

## Description

Hn-go is an Proxy API to ingest and clean data from [HackerNews](https://news.ycombinator.com) created by [ycombinator](https://www.ycombinator.com)

- I built this project to have a more accessible API for HN
- This project ingests and creates a local cache of http response data from HN's firebase API.
- Using this code structure, I focused on cleaning up the responses to be friendlier for clients to parse with minimal fetching.
- Learned concepts: Chunked API responses, local caching structures

## Table of Contents

- [HN-Go](#hn-go)
  - [Description](#description)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Usage](#usage)

## Features

- Built-in cache for quick access to clean data
- Concurrent api calls to greatly increase speed vs Node

## Usage

To run locally use `go run cmd/main.go` in your terminal

