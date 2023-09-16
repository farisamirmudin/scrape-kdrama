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
)

type Film struct {
	Title string
	Path  string
}

func ExtractTitleAndEpisode(input string) []string {
	matches := regexp.MustCompile(`(.+) Episode (.+)`).FindStringSubmatch(input)
	if len(matches) < 2 {
		return []string{}
	}

	return []string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2])}
}

func GetFilm(keyword string) []Film {
	if keyword == "" {
		return []Film{}
	}
	resp, err := http.Get(os.Getenv("BASE_URL") + "/search.html?keyword=" + keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Create a goquery document from the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	var results []Film

	doc.Find("li.video-block").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find(".name").Html()
		href, _ := s.Find("a").Attr("href")

		results = append(results, Film{Title: html.UnescapeString(title), Path: href})
	})

	return results
}

func GetFilmEpisodes(path string) []Film {
	resp, err := http.Get(os.Getenv("BASE_URL") + path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Create a goquery document from the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	var results []Film

	doc.Find("ul.listing.items.lists li").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find(".name").Html()
		href, _ := s.Find("a").Attr("href")

		results = append(results, Film{Title: html.UnescapeString(title), Path: href})
	})

	return results
}

func GetEmbeddedLink(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	matches := regexp.MustCompile(`<iframe src="(.+?)".*`).FindStringSubmatch(string(body))
	if len(matches) < 2 {
		log.Fatal("Link not found")
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

	// Encrypt and Encode the ID
	encodedID := Encrypt(id)

	// Send an HTTP request with the encoded ID
	client := &http.Client{}
	req, err := http.NewRequest("POST", ajaxURL+"?id="+encodedID, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type Data struct {
		Data string `json:"data"`
	}
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Decode and Decrypt the response
	decodedData := Decrypt(data.Data)

	// Get the m3u8 link
	links := regexp.MustCompile(`(https.+?.m3u8)`).FindAllString(string(decodedData), -1)

	parsedLink, _ := url.Parse(links[0])

	return parsedLink.String()
}
