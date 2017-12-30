package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const redirectURI = "http://localhost:8080/callback"

var callFrquency = flag.Int("callFrequency", 60, "Time in seconds between calls to the SPotify API")

var (
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
	)
	ch           = make(chan *spotify.Client)
	tk           = make(chan *oauth2.Token)
	state        = "abc123"
	clientId     = os.Getenv("SPOTIFY_ID")
	clientSecret = os.Getenv("SPOTIFY_SECRET")
	slackToken   = os.Getenv("SLACK_TOKEN")
)

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

func setSlackStatus(text string) {
	api := slack.New(slackToken)
	api.SetUserCustomStatus(text, ":headphones:")
}

func getRefreshResponse(refreshToken string) *RefreshResponse {
	b64 := fmt.Sprintf("%s:%s", clientId, clientSecret)
	authValue := base64.StdEncoding.EncodeToString([]byte(b64))
	url := "https://accounts.spotify.com/api/token"
	payload := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", refreshToken))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authorization", fmt.Sprintf("Basic %s", authValue))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error making request: %s", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var ref *RefreshResponse
	json.Unmarshal(body, &ref)
	return ref
}

func getNowPlaying(accessToken string) (*spotify.CurrentlyPlaying, error) {
	url := "https://api.spotify.com/v1/me/player/currently-playing"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error making request: %s", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ref *spotify.CurrentlyPlaying
	json.Unmarshal(body, &ref)
	return ref, nil
}

func main() {
	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Printf("Login URL: %s\n", url)
	tok := <-tk
	accessToken := tok.AccessToken

	for {
		current, err := getNowPlaying(accessToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			break
		}
		text := fmt.Sprintf("%s - %s", current.Item.Artists[0].Name, current.Item.Name)
		fmt.Println(text)
		setSlackStatus(text)
		time.Sleep(60 * time.Second)
		refresh := getRefreshResponse(tok.RefreshToken)
		accessToken = refresh.AccessToken
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	w.Header().Set("Content-Type", "text/html")
	tk <- tok
}
