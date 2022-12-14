package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

func fetchUserToken(scope string) string {
	godotenv.Load()
	const (
		redirectURL     = "http://localhost:4321"
		spotifyLoginURL = "https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s"
	)

	var (
		clientID     = os.Getenv("CLIENT_ID")
		clientSecret = os.Getenv("CLIENT_SECRET")
		authHeader   = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	)

	if clientID == "" && clientSecret == "" {
		panic(fmt.Errorf("spotify client ID and secret missing"))
	}

	// authorization code - received in callback
	code := ""
	// local state parameter for cross-site request forgery prevention
	state := fmt.Sprint(rand.Int())
	// scope of the access: we want to modify user's playlists
	// scope := "user-read-currently-playing"

	// scope := "playlist-modify-private"
	// loginURL
	path := fmt.Sprintf(spotifyLoginURL, clientID, redirectURL, scope, state)

	// channel for signaling that server shutdown can be done
	messages := make(chan bool)

	// callback handler, redirect from authentication is handled here
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// check that the state parameter matches
		if s, ok := r.URL.Query()["state"]; ok && s[0] == state {
			// code is received as query parameter
			if codes, ok := r.URL.Query()["code"]; ok && len(codes) == 1 {
				// save code and signal shutdown
				code = codes[0]
				messages <- true
			}
		}
		// redirect user's browser to spotify home page
		http.Redirect(w, r, "https://www.spotify.com/", http.StatusSeeOther)
	})

	// open user's browser to login page
	if err := browser.OpenURL(path); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}

	server := &http.Server{Addr: ":4321"}
	// go routine for shutting down the server
	go func() {
		okToClose := <-messages
		if okToClose {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Println("Failed to shutdown server", err)
			}
		}
	}()
	// log.Println(server.ListenAndServe())
	server.ListenAndServe()

	data1 := url.Values{}
	data1.Set("grant_type", "authorization_code")
	data1.Set("code", code)
	data1.Set("redirect_uri", redirectURL)
	encodedData := data1.Encode()

	url1 := "https://accounts.spotify.com/api/token"
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url1, strings.NewReader(encodedData))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	data, _ := io.ReadAll(res.Body)
	if err == nil {
		response := AuthResponse{}
		if err = json.Unmarshal(data, &response); err == nil {
			// happy end: token parsed successfully
			return response.AccessToken
		}
	}
	fmt.Println(err)
	panic(fmt.Errorf("unable to acquire Spotify user token2"))
}
