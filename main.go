package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"golang.org/x/term"
)

var Token TokenList

func loadJson() {
	tokenFile, err := os.Open("tokens.json")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// fmt.Println("Successfully Opened users.json")
	defer tokenFile.Close()

	byteValue, _ := io.ReadAll(tokenFile)
	// fmt.Println("error", err)
	json.Unmarshal(byteValue, &Token)
	// spew.Dump(byteValue)
}

func saveToken() {

	file, _ := json.MarshalIndent(Token, "", "    ")
	// fmt.Println(Token)
	// fmt.Println("\n\n\n\n\n")
	// err := os.WriteFile("tokens.json", file, 0644)
	_ = os.WriteFile("tokens.json", file, 0644)
	// fmt.Println(err)
	// fmt.Println("\n\n\n\n\n")
	// fmt.Println("\n\n\n\n\n")
}

func main() {
	// godotenv.Load()
	loadJson()
	token1 := fetchUserToken("user-read-currently-playing")
	// fmt.Println(token1)
	Token.UserReadCurrentlyPlaying = token1
	saveToken()

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
	playStatus := name.Newline()
	playProgress := name.Newline()

	name.Start()

	failCount := 0

	for {

		currentSong, err := loadSong(Token.UserReadCurrentlyPlaying)
		if err != nil {
			failCount++
		} else {
			failCount = 0

			fmt.Fprintf(name, "Name: %s\n", magenta(currentSong.Name))
			fmt.Fprintf(album, "Album: %s\n", magenta(currentSong.AlbumName))
			fmt.Fprintf(artist, "Artist: %s\n", magenta(currentSong.ArtistName))

			fmt.Fprintf(emptyLine, "\n")
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
			fmt.Fprintf(playProgress, "%s%s\n", green(plus), yellow(minus))

			// for i := 0; i < leftHeight; i++ {
			// 	fmt.Println("")
			// }
			time.Sleep(500 * time.Millisecond)
		}

		if failCount > 3 {
			fmt.Println(err)
			return
		}
	}
	// name.Stop()
}

func loadSong(token1 string) (currentTrack, error) {
	url1 := "https://api.spotify.com/v1/me/player/currently-playing"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url1, nil)
	req.Header.Set("Authorization", "Bearer "+Token.UserReadCurrentlyPlaying)
	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)
	// var d any

	var d currentTrackAPI
	json.Unmarshal([]byte(body), &d)
	// fmt.Printf("%+v\n", d)

	// return currentTrack{}, errors.New("a")

	if d.Error.Status == 401 {
		fmt.Println("get new token")
		fmt.Println(d.Error.Message)
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(Token)
		panic("panic")
		token1 := fetchUserToken("user-read-currently-playing")
		Token.UserReadCurrentlyPlaying = token1
		fmt.Println(Token)
		saveToken()
		loadJson()
		// fmt.Println(Token)
		return loadSong(Token.UserReadCurrentlyPlaying)
	}
	if (d.Error != errorMsg{}) {
		return currentTrack{}, errors.New(d.Error.Message)
	}
	if d.Item.Name == "" {
		return currentTrack{}, errors.New("no retrun")
	}

	data := currentTrack{ProgressMs: d.ProgressMs, IsPlaying: d.IsPlaying, AlbumName: d.Item.Album.Name, ArtistName: d.Item.Artists[0].Name, DurationMs: d.Item.DurationMs, Href: d.Item.Href, Name: d.Item.Name}
	return data, nil
}

type TokenList struct {
	UserReadCurrentlyPlaying string `json:"user-read-currently-playing"`
	UserReadPlaybackState    string `json:"user-read-playback-state"`
	UserModifyPlaybackState  string `json:"user-modify-playback-state"`
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
