package lib

import (
	"bytes"
	"encoding/json"
	"html"

	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/farisamirmudin/gowatch/model"
)

func GetFilm(keyword string) []model.Film {
	if keyword == "" {
		return []model.Film{}
	}
	resp, err := http.Get(os.Getenv("BASE_URL") + "/search.html?keyword=" + keyword)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	// Create a goquery document from the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Print(err)
	}

	var results []model.Film

	doc.Find("li.video-block").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find(".name").Html()
		href, _ := s.Find("a").Attr("href")

		results = append(results, model.Film{Title: html.UnescapeString(title), Path: href})
	})

	return results
}

func GetFilmEpisodes(path string) []model.Film {
	resp, err := http.Get(os.Getenv("BASE_URL") + path)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	// Create a goquery document from the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Print(err)
	}

	var results []model.Film

	doc.Find("ul.listing.items.lists li").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find(".name").Html()
		href, _ := s.Find("a").Attr("href")

		results = append(results, model.Film{Title: html.UnescapeString(title), Path: href})
	})

	return results
}

func GetEmbeddedLink(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	matches := regexp.MustCompile(`<iframe src="(.+?)".*`).FindStringSubmatch(string(body))
	if len(matches) < 2 {
		log.Print("Link not found")
	}

	return "https:" + matches[1]
}

func GetFilmServers(path string) string {
	ajaxURL := os.Getenv("BASE_URL") + "/encrypt-ajax.php"
	searchLink := GetEmbeddedLink(os.Getenv("BASE_URL") + path)

	// Extract the ID from the URL
	match := regexp.MustCompile(`id=(.+?)&`).FindStringSubmatch(searchLink)[1]
	id := []byte(match)

	// Pad the id
	blockSize := 16
	padding := blockSize - len(id)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	id = append(id, padtext...)

	// Encrypt the ID
	encryptedID := Encrypt(id)

	// Send an HTTP request with the encrypted ID
	client := &http.Client{}
	req, err := http.NewRequest("POST", ajaxURL+"?id="+encryptedID, nil)
	if err != nil {
		log.Print(err)
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Print(err)
	}
	type Data struct {
		Data string `json:"data"`
	}
	var data Data

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	// Decrypt the response
	decryptedData := Decrypt(data.Data)

	// Get the m3u8 link
	links := regexp.MustCompile(`(https.+?.m3u8)`).FindAllString(string(decryptedData), -1)

	parsedLink, _ := url.Parse(links[0])

	return parsedLink.String()
}
