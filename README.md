# Plex Authentication (in Go!)

[![GoDoc](https://pkg.go.dev/badge/github.com/BrenekH/go-plexauth)](https://pkg.go.dev/github.com/BrenekH/go-plexauth)

This tool is based on the instructions available in [this Plex forum post](https://forums.plex.tv/t/authenticating-with-plex/609370), which you should look over to understand the authentication flow.
To reiterate the high level steps, a PIN must be generated and formed into a URL that the user can visit which allows them to claim the PIN.
After the PIN has been claimed, the Plex API will return an authentication token which can be used to access the Plex API on the user's behalf.

Go PlexAuth automates this process making it easy to obtain a permanent auth token for API calls.

## Installation

As a Go library: `go get github.com/BrenekH/go-plexauth`

As a CLI: `go install github.com/BrenekH/go-plexauth/cli@latest`

## Usage

### Go Library

Visit the documentation on [pkg.go.dev](https://pkg.go.dev/github.com/BrenekH/go-plexauth).

### CLI

Run the executable and enter a device name.
You will then need to visit the URL it prints out to the console to complete the authentication.
Once authentication is complete, your authentication token will be printed to the console.
