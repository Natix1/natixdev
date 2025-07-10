package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	profileURL    string
	tmpl          = template.Must(template.ParseFiles("templates/index.html"))
	bindAddr      = "0.0.0.0:8080"
	client        http.Client
	apiProfileURL = "https://thumbnails.roblox.com/v1/users/avatar-headshot?userIds=1862042823&size=420x420&format=Png&isCircular=false&thumbnailType=HeadShot"

	robloxProfileURL = "https://www.roblox.com/users/1862042823/profile"
	robloxGroupURL   = "https://www.roblox.com/communities/33883891/Ronalds#!/about"
)

type RobloxAPIData struct {
	ImageURL string `json:"imageUrl"`
}

type RobloxAPIResp struct {
	Data []RobloxAPIData `json:"data"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		RobloxImageURL   string
		RobloxProfileURL string
		RobloxGroupURL   string
	}{
		RobloxImageURL:   profileURL,
		RobloxProfileURL: robloxProfileURL,
		RobloxGroupURL:   robloxGroupURL,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}
}

func fetchProfileURL() {
	resp, err := client.Get(apiProfileURL)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("Non-200 for profile url periodic fetch: ", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var respData RobloxAPIResp
	err = json.Unmarshal(body, &respData)

	if err != nil {
		log.Println(err)
		return
	}

	if len(respData.Data) < 1 {
		log.Println("No data received from roblox API")
		return
	}

	profileURL = respData.Data[0].ImageURL
	log.Println("Refresh of profile URL done")
}

func main() {
	fetchProfileURL()
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			fetchProfileURL()
		}
	}()

	http.HandleFunc("/", rootHandler)
	log.Println("Server running at", bindAddr)
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
