# Plex Authentication (in Go!)

[![GoDoc](https://pkg.go.dev/badge/github.com/BrenekH/go-plexauth)](https://pkg.go.dev/github.com/BrenekH/go-plexauth)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/BrenekH/go-plexauth?label=version)
[![License](https://img.shields.io/github/license/BrenekH/go-plexauth)](https://github.com/BrenekH/go-plexauth/tree/master/LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/BrenekH/go-plexauth)

This tool is based on the instructions available in [this Plex forum post](https://forums.plex.tv/t/authenticating-with-plex/609370), which you should look over to understand the authentication flow.
To reiterate the high level steps, a PIN must be generated and formed into a URL that the user can visit which allows them to claim the PIN.
After the PIN has been claimed, the Plex API will return an authentication token which can be used to access the Plex API on the user's behalf.

Go PlexAuth automates this process making it easy to obtain a permanent auth token for API calls.

## Installation

As a Go library: `go get github.com/BrenekH/go-plexauth`

As a CLI: `go install github.com/BrenekH/go-plexauth/cli@latest`

## Usage

### CLI

Run the executable and enter a device name.
You will then need to visit the URL it prints out to the console to complete the authentication.
Once authentication is complete, your authentication token will be printed to the console.

### Go Library

Basic example:

```go
package main

import (
    "fmt"

    "github.com/BrenekH/go-plexauth"
)

func main() {
    clientID := "my-random-client-id"

    // Request that Plex generate a new PIN for a user to claim.
    pinID, pinCode, err := plexauth.GetPlexPIN(appName, clientID)
    if err != nil {
        panic(err)
    }

    // Create the URL that the user needs to visit to claim the PIN.
    authUrl, err := plexauth.GenerateAuthURL(appName, clientID, pinCode, plexauth.ExtraAuthURLOptions{})
    if err != nil {
        panic(err)
    }

    // Inform the user of the URL they need to visit. Typically this would
    // automatically open the browser, but here it is printed to the console for simplicity's sake.
    fmt.Printf("Please visit %s to authenticate.\n", authUrl)

    // Repeatedly ask Plex if the PIN has been claimed yet. Times out after 30 minutes.
    authToken, err := plexauth.PollForAuthToken(context.Background(), pinID, pinCode, clientID)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Your authentication token is \"%s\"\n", authToken)
}
```

For the full documentation, visit [pkg.go.dev/github.com/BrenekH/go-plexauth](https://pkg.go.dev/github.com/BrenekH/go-plexauth).

## License

This project is licensed under the Apache 2.0 License, a copy of which can be found in [LICENSE](https://github.com/BrenekH/go-plexauth/tree/master/LICENSE).
