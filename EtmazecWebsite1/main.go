package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rapito/go-spotify/spotify"
)

type TopHitsJSONParsed struct {
	Tracks []struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		DiscNumber  int  `json:"disc_number"`
		DurationMs  int  `json:"duration_ms"`
		Explicit    bool `json:"explicit"`
		ExternalIds struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string      `json:"href"`
		ID          string      `json:"id"`
		IsLocal     bool        `json:"is_local"`
		IsPlayable  bool        `json:"is_playable"`
		Name        string      `json:"name"`
		Popularity  int         `json:"popularity"`
		PreviewURL  interface{} `json:"preview_url"`
		TrackNumber int         `json:"track_number"`
		Type        string      `json:"type"`
		URI         string      `json:"uri"`
	} `json:"tracks"`
}

type SingerJSONParsed struct {
	Artists struct {
		Href  string `json:"href"`
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  interface{} `json:"href"`
				Total int         `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
			URI        string `json:"uri"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"artists"`
}

type AlbumJSONParsed struct {
	Href  string `json:"href"`
	Items []struct {
		AlbumGroup string `json:"album_group"`
		AlbumType  string `json:"album_type"`
		Artists    []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		TotalTracks          int    `json:"total_tracks"`
		Type                 string `json:"type"`
		URI                  string `json:"uri"`
	} `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

var tpl *template.Template

type pageData struct {
	Title     string
	Firstname string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func databaseAlbums(str1 string, str2 string, flag bool) {

	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS albumshistorytable (id INTEGER PRIMARY KEY, albumsname TEXT, albumartist TEXT)")
	statement.Exec()

	//Inserting the parsed JSON of artist's name and their albums' name
	statement, _ = database.Prepare("INSERT INTO albumshistorytable (albumsname, albumartist) VALUES (?, ?)")
	statement.Exec(str1, str2)
	rows, _ := database.Query("SELECT id, albumsname, albumartist FROM albumshistorytable")
	//Printing the db of album
	if flag == true {
		fmt.Println("DATABASE ALBUMS SEARCH HISTORY : ")
		var id int
		var albumsnamedb string
		var albumartistdb string
		for rows.Next() {
			rows.Scan(&id, &albumsnamedb, &albumartistdb)
			fmt.Println(strconv.Itoa(id), " ", albumsnamedb, albumartistdb)
		}
	}
}

func databaseTopHits(str1 string, str2 string, flag bool) {
	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS tophitshistorytable (id INTEGER PRIMARY KEY, trackname TEXT, trackpopularity TEXT)")
	statement.Exec()

	//Inserting the parsed JSON of artist's name and their albums' name
	statement, _ = database.Prepare("INSERT INTO tophitshistorytable (trackname, trackpopularity) VALUES (?, ?)")
	statement.Exec(str1, str2)
	rows, _ := database.Query("SELECT id, trackname, trackpopularity FROM tophitshistorytable")
	//Printing the db of album
	if flag == true {
		fmt.Println("DATABASE TOP HITS SEARCH HISTORY : ")
		var id int
		var track string
		var popularity string
		for rows.Next() {
			rows.Scan(&id, &track, &popularity)
			fmt.Println(strconv.Itoa(id), " ", track, popularity)
		}
	}

}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/home", home)
	http.HandleFunc("/albums", albums)
	http.HandleFunc("/tophits", tophits)
	http.ListenAndServe(":5051", nil)

}

func home(w http.ResponseWriter, req *http.Request) {

	pd := pageData{
		Title: "Etmazec Home Page",
	}

	err := tpl.ExecuteTemplate(w, "homepage.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func albums(w http.ResponseWriter, req *http.Request) {

	pd := pageData{
		Title: "Albums Page",
	}

	var first string
	if req.Method == http.MethodPost {
		first = req.FormValue("fname")
		pd.Firstname = first
		artistalbumsget(first)
	}
	err := tpl.ExecuteTemplate(w, "albums.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

	}
}
func tophits(w http.ResponseWriter, req *http.Request) {
	pd := pageData{
		Title: "Top Hits Page",
	}

	var first string
	if req.Method == http.MethodPost {
		first = req.FormValue("fname")
		pd.Firstname = first
		artistshitsget(first)
	}
	err := tpl.ExecuteTemplate(w, "tophits.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

	}
}

func artistalbumsget(word string) {

	spot := spotify.New("756100c43d724ae6b791f7804c82b219", "841a197013f847bca78c01b3b69fc72d")
	authorized, _ := spot.Authorize()

	if authorized {

		artistIDget := join("search?query=", word, "&type=artist&market=EG&offset=0&limit=1")
		response, _ := spot.Get(artistIDget, nil)
		var artistIdstring SingerJSONParsed
		if err := json.Unmarshal([]byte(response), &artistIdstring); err != nil {
			fmt.Println("ugh: ", err)
		}
		idofartist := artistIdstring.Artists.Items[0].ID

		albumsretrieved := join("artists/", idofartist, "/albums?market=ES&limit=10")
		response, _ = spot.Get(albumsretrieved, nil)
		var i AlbumJSONParsed
		if err := json.Unmarshal([]byte(response), &i); err != nil {
			fmt.Println("ugh: ", err)
		}
		flag := false
		for j := 0; j < 10; j++ {
			if j == 9 {
				flag = true
			}
			databaseAlbums(join("Album Name : ", i.Items[j].Name, "\n"), join("Artist Name: ", i.Items[j].Artists[0].Name, "\n"), flag)
		}
		fmt.Println("CURRENT ARTIST : ", word, "'s ALBUMS RETRIEVAL : ")

		for k := 0; k < 10; k++ {
			fmt.Println()
			fmt.Println("Album Name : ", i.Items[k].Name)
			fmt.Println()
			fmt.Println("Artist Name: ", i.Items[k].Artists[0].Name)
			fmt.Println()
		}

	}

	// Parse response to a JSON Object and
}
func artistshitsget(word string) string {
	retstring := ""
	spot := spotify.New("756100c43d724ae6b791f7804c82b219", "841a197013f847bca78c01b3b69fc72d")
	authorized, _ := spot.Authorize()

	if authorized {

		artistIDget := join("search?query=", word, "&type=artist&market=EG&offset=0&limit=1")
		response, _ := spot.Get(artistIDget, nil)
		var artistIdstring SingerJSONParsed
		if err := json.Unmarshal([]byte(response), &artistIdstring); err != nil {
			fmt.Println("ugh: ", err)
		}
		idofartist := artistIdstring.Artists.Items[0].ID

		tophitsretrieved := join("artists/", idofartist, "/top-tracks?country=ES")
		response, _ = spot.Get(tophitsretrieved, nil)
		var i TopHitsJSONParsed
		if err := json.Unmarshal([]byte(response), &i); err != nil {
			fmt.Println("ugh: ", err)
		}

		flag := false

		for k := 0; k < 5; k++ {
			if k == 4 {
				flag = true
			}
			databaseTopHits(join("Track Title : ", i.Tracks[k].Name, "\n"), join("Popularity : ", strconv.Itoa(i.Tracks[k].Popularity), "\n"), flag)
		}
		fmt.Println("CURRENT ARTIST : ", word, "'s TOP HITS RETRIEVAL : ")

		for j := 0; j < 5; j++ {
			retstring = join(retstring, "\n")
			fmt.Println()
			retstring = join(retstring, "Track Title : ", i.Tracks[j].Name)
			fmt.Println("Track Title : ", i.Tracks[j].Name)
			retstring = join(retstring, "\n")
			fmt.Println()
			retstring = join(retstring, "Popularity : ", strconv.Itoa(i.Tracks[j].Popularity))
			fmt.Println("Popularity of track : ", i.Tracks[j].Popularity)
			retstring = join(retstring, "\n")
			fmt.Println()
		}

	}
	return retstring
}

func join(strs ...string) string {
	var ret string
	for _, str := range strs {
		ret += str
	}
	return ret
}
