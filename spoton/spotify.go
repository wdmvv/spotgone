package spoton

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ClientID     string
	ClientSecret string
	AuthToken    string `json:"access_token"`
}

type Album struct {
	Name   string `json:"name"`
	Total  int    `json:"total"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Release    string `json:"release_date"`
	Label      string `json:"label"`
	Popularity int    `json:"popularity"`
	Tracks     []struct {
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name   string `json:"name"`
		Number int    `json:"track_number"`
		Disc   int    `json:"disc_number"`
	} `json:"items"`
}

// TODO: https://developer.spotify.com/documentation/web-api/reference/get-playlists-tracks
type Track struct {
	Album struct {
	} `json:"album"`
}

type Playlist struct {
}

// i have no idea how to do this efficiently so almost same structs /shrug
// this one is for processing response while capital one is for storing tracks only
// ill see how this turns out to be and maybe will tweak stuff, who knows
type playlistResponse struct {
	Total int `json:"total"`
	Items []struct {
		Track struct {
			Album struct {
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name    string `json:"name"`
				Release string `json:"release_date"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Disc       int    `json:"disc_number"`
			DurationMs int    `json:"duration_ms"`
			Name       string `json:"name"`
			Number     int    `json:"track_number"`
			Popularity int    `json:"popularity"`
		} `json:"track"`
	} `json:"items"`
}

// this should be in some kind of config but eh
var (
	tokenEndpoint string = "https://accounts.spotify.com/api/token"
)

// Gets auth token and sets in User struct
func (u *User) SetAuth() error {
	auth := fmt.Sprintf("%s:%s", u.ClientID, u.ClientSecret)
	authb64 := b64.StdEncoding.EncodeToString([]byte(auth))

	data := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+authb64)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return &ErrBadRequest{"auth", resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)
	return d.Decode(&u)
}

// if something's missing then one should set market in request
func (u *User) GetAlbum(id string) (Album, error) {
	var a Album
	if u.AuthToken == "" {
		return Album{}, &ErrNoAuth{}
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Album{}, err
	}

	req.Header.Add("Authorization", "Bearer "+u.AuthToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Album{}, err
	} else if resp.StatusCode != http.StatusOK {
		return Album{}, &ErrBadRequest{"album", resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&a)
	if err != nil {
		return Album{}, err
	}

	url = fmt.Sprintf("https://api.spotify.com/v1/albums/%s", id)
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Album{}, err
	}

	req.Header.Add("Authorization", "Bearer "+u.AuthToken)
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return Album{}, err
	} else if resp.StatusCode != http.StatusOK {
		return Album{}, &ErrBadRequest{"album", resp.StatusCode}
	}

	d = json.NewDecoder(resp.Body)
	err = d.Decode(&a)
	return a, err
}

func (u *User) GetPlaylist(id string) (Playlist, error) {
	if u.AuthToken == "" {
		return Playlist{}, &ErrNoAuth{}
	}

	// yes i couldve used a body for this but why bother? this is not a complex request, just simple math
	// if this was bigger then i wouldve used struct -> json marshaling and proper body and not nil
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=0&limit=50", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Playlist{}, err
	}
	req.Header.Add("Authorization", "Bearer "+u.AuthToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Playlist{}, err
	} else if resp.StatusCode != http.StatusOK {
		return Playlist{}, &ErrBadRequest{"playlist", resp.StatusCode}
	}
	d := json.NewDecoder(resp.Body)
	var pl playlistResponse
	err = d.Decode(&pl)
	if err != nil {
		return Playlist{}, err
	}
	return Playlist{}, nil
}
