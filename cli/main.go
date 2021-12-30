package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/BrenekH/go-plexauth"
)

const appName string = "Go PlexAuth CLI"
const version string = "0.0.2"

func main() {
	deviceName := input("Please enter a device name: ")

	clientID := fmt.Sprintf("go-plexauth-%s", randSeq(10))

	pinID, pinCode, err := plexauth.GetPlexPIN(appName, clientID)
	if err != nil {
		panic(err)
	}

	authUrl, err := plexauth.GenerateAuthURL(appName, clientID, pinCode, plexauth.ExtraAuthURLOptions{DeviceName: deviceName, AppVersion: version})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Please visit %s to authenticate.\n", authUrl)

	authToken, err := plexauth.PollForAuthToken(context.Background(), pinID, pinCode, clientID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Your authentication token is \"%s\"\n", authToken)
}

// input mimics Python's input function, which outputs a prompt and
// takes bytes from stdin until a newline and returns a string.
func input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); ok {
		return scanner.Text()
	}
	return ""
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// randSeq creates a random, alphanumeric string of n characters.
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
