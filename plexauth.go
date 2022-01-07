package plexauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// IsTokenValid checks if a Plex token is valid using the Plex API.
func IsTokenValid(appName, clientID, token string) (bool, error) {
	return IsTokenValidContext(context.Background(), appName, clientID, token)
}

// IsTokenValidContext checks if a Plex token is valid using the Plex API and a custom request context.
func IsTokenValidContext(ctx context.Context, appName, clientID, token string) (bool, error) {
	data := strings.NewReader(fmt.Sprintf(`X-Plex-Product=%s&X-Plex-Client-Identifier=%s&X-Plex-Token=%s`, appName, clientID, token))

	req, err := http.NewRequestWithContext(ctx, "GET", "https://plex.tv/api/v2/user", data)
	if err != nil {
		return false, fmt.Errorf("IsTokenValid: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("IsTokenValid: %w", err)
	}
	resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return true, nil
	case 401:
		return false, nil
	default:
		return false, fmt.Errorf("IsTokenValid: unable to determine token validity, received status code %v", resp.StatusCode)
	}
}

// GetPlexPIN requests a new claimable pin from the Plex API which can be used to authenticate a user.
func GetPlexPIN(appName, clientID string) (pinID int, pinCode string, err error) {
	return GetPlexPINContext(context.Background(), appName, clientID)
}

// GetPlexPINContext requests a new claimable pin from the Plex API which can be used to authenticate a user.
func GetPlexPINContext(ctx context.Context, appName, clientID string) (pinID int, pinCode string, err error) {
	data := strings.NewReader(fmt.Sprintf(`strong=true&X-Plex-Product=%s&X-Plex-Client-Identifier=%s`, appName, clientID))

	req, err := http.NewRequestWithContext(ctx, "POST", "https://plex.tv/api/v2/pins", data)
	if err != nil {
		return 0, "", fmt.Errorf("GetPlexPIN: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("GetPlexPIN: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("GetPlexPIN: %w", err)
	}

	pinResp := plexPINResponse{}
	if err = json.Unmarshal(b, &pinResp); err != nil {
		return 0, "", fmt.Errorf("GetPlexPIN: %w", err)
	}

	return pinResp.ID, pinResp.Code, nil
}

// GenerateAuthURL creates an authorization link that a user can visit to authorize an application.
// The application name, unique client id, and pin code are all required for the link to work.
//
// To send extra info that can be displayed on the Authorized Devices dashboard, the ExtraAuthURLOptions struct can be used.
func GenerateAuthURL(appName, clientID, pinCode string, extraOpts ExtraAuthURLOptions) (string, error) {
	// Copy all parameters into a single encodeable struct
	queryOpts := authURLQueryOpts{
		AppName:         appName,
		ClientID:        clientID,
		PinCode:         pinCode,
		AppVersion:      extraOpts.AppVersion,
		DeviceName:      extraOpts.DeviceName,
		Device:          extraOpts.Device,
		Platform:        extraOpts.Platform,
		PlatformVersion: extraOpts.PlatformVersion,
	}

	v, err := query.Values(queryOpts)
	if err != nil {
		return "", fmt.Errorf("GenerateAuthURL: %w", err)
	}

	optsStr := v.Encode()

	return "https://app.plex.tv/auth#?" + optsStr, nil
}

// PollForAuthToken gets a new auth token by waiting for the user to authenticate the pin in a web browser.
//
// A custom context is required, but will be restricted to a 30 minute timeout. This is because Plex pins are
// only valid for 30 minutes.
func PollForAuthToken(inCtx context.Context, pinID int, pinCode, clientID string) (string, error) {
	// Set a maximum timeout to 30 minutes, since that's how long a pin is good for.
	ctx, cancel := context.WithTimeout(inCtx, 30*time.Minute)
	defer cancel()

	for {
		select {
		case <-time.After(time.Second):
			//? Should errors be cause for immediate termination, or should a simple warning be printed to the console?

			data := strings.NewReader(fmt.Sprintf(`code=%s&X-Plex-Client-Identifier=%s`, pinCode, clientID))
			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://plex.tv/api/v2/pins/%v", pinID), data)
			if err != nil {
				return "", fmt.Errorf("PollForAuthToken: %w", err)
			}

			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return "", fmt.Errorf("PollForAuthToken: %w", err)
			}

			b, err := io.ReadAll(resp.Body)
			resp.Body.Close() // To avoid accumulating 1800 response bodies using defer, just close the body as soon as it's done being used.
			if err != nil {
				return "", fmt.Errorf("PollForAuthToken: %w", err)
			}

			respStruct := authTokenResp{}
			if err = json.Unmarshal(b, &respStruct); err != nil {
				return "", fmt.Errorf("PollForAuthToken: %w", err)
			}

			if respStruct.AuthToken != nil {
				return *(respStruct.AuthToken), nil
			}

		case <-ctx.Done():
			return "", fmt.Errorf("PollForAuthToken: could not retrieve auth token, exceeded context")
		}
	}
}

// ExtraAuthURLOptions provides a way to provide extra metadata about the device being authorized.
// Any zero-value will be omitted from the final authentication url.
type ExtraAuthURLOptions struct {
	AppVersion      string
	DeviceName      string
	Device          string // Small descriptor of the device
	Platform        string // Determines what icon is used in the authorized devices dashboard.
	PlatformVersion string
}

type authURLQueryOpts struct {
	ClientID        string `url:"clientID"`
	PinCode         string `url:"code"`
	AppName         string `url:"context[device][product]"`
	AppVersion      string `url:"context[device][version],omitempty"`
	DeviceName      string `url:"context[device][deviceName],omitempty"`
	Device          string `url:"context[device][device],omitempty"`
	Platform        string `url:"context[device][platform],omitempty"`
	PlatformVersion string `url:"context[device][platformVersion],omitempty"`
}

type plexPINResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

type authTokenResp struct {
	AuthToken *string `json:"authToken"`
}
