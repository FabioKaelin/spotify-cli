package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/term"
)

func main() {
	godotenv.Load()

	if !term.IsTerminal(0) {
		println("not in a term")
	}

	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}

	// usedLines := 0

	currentSong := loadSong()

	fmt.Println("Name:", currentSong.Name)
	fmt.Println("Album:", currentSong.AlbumName)
	fmt.Println("Artist:", currentSong.ArtistName)
	fmt.Println("")
	height -= 5

	if currentSong.IsPlaying {
		fmt.Println("|> (Playing) --", fmt.Sprint(currentSong.ProgressMs/1000)+"/"+fmt.Sprint(currentSong.DurationMs/1000))
	} else {
		fmt.Println("|| (paused)")
	}
	height--

	listenedChars := math.Round(float64(width) / 100.0 * float64(100.0/float64(currentSong.DurationMs)) * float64(currentSong.ProgressMs))
	for i := 0; i < int(listenedChars); i++ {
		fmt.Print("+")
	}
	for i := 0; i < width-int(listenedChars); i++ {
		fmt.Print("-")
	}
	height--
	fmt.Println("")
	// for true {
	// 	fmt.Println("hallo")
	// }
}

func loadSong() currentTrack {
	url := "https://api.spotify.com/v1/me/player/currently-playing"
	// header = { "Authorization" : "Bearer " + config.get("TOKEN_user_read_recently_played") }
	// r = requests.get(url, headers=header)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_user_read_recently_played"))
	res, _ := client.Do(req)
	// var d any
	var d currentTrackAPI
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &d)
	// fmt.Printf("%+v\n", d)
	// fmt.Printf("%+v\n", d.Item.Artists)
	// data := currentTrack{ProgressMs: d.ProgressMs, IsPlaying: d.IsPlaying, AlbumName: d.Item.Album.Name, DurationMs: d.Item.DurationMs, Href: d.Item.Href, Name: d.Item.Name}
	data := currentTrack{ProgressMs: d.ProgressMs, IsPlaying: d.IsPlaying, AlbumName: d.Item.Album.Name, ArtistName: d.Item.Artists[0].Name, DurationMs: d.Item.DurationMs, Href: d.Item.Href, Name: d.Item.Name}
	// fmt.Printf("%+v\n", data)
	return data
	// return currentTrack{}
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
