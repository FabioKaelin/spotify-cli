package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/joho/godotenv"
	"golang.org/x/term"
)

var Token string

func main() {
	godotenv.Load()
	Token = os.Getenv("TOKEN_user_read_recently_played")

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	if !term.IsTerminal(0) {
		fmt.Println(red("not in a term"))
		return
	}

	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}

	for i := 0; i < height; i++ {

	}

	name := uilive.New()
	album := name.Newline()
	artist := name.Newline()
	emptyLine := name.Newline()
	fmt.Fprint(emptyLine, "")
	playStatus := name.Newline()
	playProgress := name.Newline()

	name.Start()

	// for i := 0; i <= 100; i++ {
	// 	fmt.Fprintf(name, "Downloading File 1.. %d %%\n", i)
	// 	fmt.Fprintf(album, "Downloading File 2.. %d %%\n", i)
	// 	fmt.Fprintf(artist, "Downloading File 3.. %d %%\n", i)
	// 	time.Sleep(time.Millisecond * 5)
	// }

	// return

	for {

		currentSong, err := loadSong()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(name, "Name: %s", magenta(currentSong.Name))
		fmt.Fprintf(album, "Album: %s", magenta(currentSong.AlbumName))
		fmt.Fprintf(artist, "Artist: %s", magenta(currentSong.ArtistName))

		if currentSong.IsPlaying {

			fmt.Fprintf(playStatus, "%s -- %d.%02d/%d.%02d\n", green("|> (Playing)"), (currentSong.ProgressMs/1000)/60, (currentSong.ProgressMs/1000)%60, (currentSong.DurationMs/1000)/60, (currentSong.DurationMs/1000)%60)
		} else {
			fmt.Fprintf(playStatus, "%s -- %d.%02d/%d.%02d\n", red("|| (paused)"), (currentSong.ProgressMs/1000)/60, (currentSong.ProgressMs/1000)%60, (currentSong.DurationMs/1000)/60, (currentSong.DurationMs/1000)%60)
		}

		listenedChars := math.Round(float64(width) / 100.0 * float64(100.0/float64(currentSong.DurationMs)) * float64(currentSong.ProgressMs))
		plus := ""
		for i := 0; i < int(listenedChars); i++ {
			plus = plus + "+"
		}
		minus := ""
		for i := 0; i < width-int(listenedChars); i++ {
			minus = minus + "-"
		}
		fmt.Fprintf(playProgress, "%s%s", green(plus), yellow(minus))

		// for i := 0; i < leftHeight; i++ {
		// 	fmt.Println("")
		// }
		time.Sleep(500 * time.Millisecond)
	}
	// name.Stop()
}

func getNewToken() error {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	encodedData := data.Encode()

	url1 := "https://accounts.spotify.com/api/token"
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url1, strings.NewReader(encodedData))

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	str := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)))

	req.Header.Set("Authorization", "Basic "+str)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)
	var d tokenResponseAPI
	json.Unmarshal([]byte(body), &d)
	if d.Error != "" {
		return errors.New(d.Error)
	}

	Token = d.AccessToken

	return nil
}

// func loadSong() currentTrack {
func loadSong() (currentTrack, error) {
	// url1 := "https://api.spotify.com/v1/me/following?type=artist"
	url1 := "https://api.spotify.com/v1/me"
	// url1 := "https://api.spotify.com/v1/me/player/currently-playing"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url1, nil)
	req.Header.Set("Authorization", "Bearer "+Token)
	// req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_user_read_recently_played"))
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)
	// var d any

	var d currentTrackAPI
	json.Unmarshal([]byte(body), &d)
	// fmt.Printf("%+v\n", d)
	if d.Error.Status == 401 {
		err := getNewToken()
		if err != nil {
			return currentTrack{}, err
		}
		fmt.Println(d.Error.Message)
		return loadSong()
	}
	if (d.Error != errorMsg{}) {
		// fmt.Println(d.Error.Status)
		return currentTrack{}, errors.New(d.Error.Message)
	}
	if d.Item.Name == "" {
		return currentTrack{}, errors.New("no retrun")
	}

	data := currentTrack{ProgressMs: d.ProgressMs, IsPlaying: d.IsPlaying, AlbumName: d.Item.Album.Name, ArtistName: d.Item.Artists[0].Name, DurationMs: d.Item.DurationMs, Href: d.Item.Href, Name: d.Item.Name}
	return data, nil
	// return currentTrack{}
}

type tokenResponseAPI struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

type currentTrackAPI struct {
	ProgressMs int  `json:"progress_ms"`
	IsPlaying  bool `json:"is_playing"`
	Item       struct {
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
		Artists    []artistsAPI `json:"artists"`
		DurationMs int          `json:"duration_ms"`
		Href       string       `json:"href"`
		Name       string       `json:"name"`
	} `json:"item"`
	Error errorMsg `json:"error"`
}

type errorMsg struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type artistsAPI struct {
	Name string `json:"name"`
}

type currentTrack struct {
	ProgressMs int    `json:"progress_ms"`
	IsPlaying  bool   `json:"is_playing"`
	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`
	DurationMs int    `json:"duration_ms"`
	Href       string `json:"href"`
	Name       string `json:"name"`
}
