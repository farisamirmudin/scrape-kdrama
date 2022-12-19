package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func Choose() uint {
	var idx uint
	c := color.New(color.FgGreen)
	fmt.Print(c.Sprint("Choose: "))
	fmt.Scanln(&idx)
	return idx
}

func GetLink(link string, keyword *string) (string, string) {
	if keyword != nil {
		color.Magenta("Searching Episode...")
	} else {
		color.Magenta("Searching Drama...")
	}
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// Create a goquery document from the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		panic(err)
	}

	var options []string
	a := "li.video-block > a"
	if keyword != nil {
		a = fmt.Sprintf("ul.listing.items.lists > li.video-block > a[href*='%s']", *keyword)
	}
	doc.Find(a).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			options = append(options, href)
			fmt.Printf("[%d] %s\n", i+1, strings.ReplaceAll(href[8:], `-`, ` `))
		}
	})

	idx := Choose()
	title, _ := doc.Find(fmt.Sprintf("a[href='%s']", options[idx-1])).Find("div.name").Html()
	return options[idx-1], strings.TrimSpace(title)
}

func GetEmbeddedLink(_searchLink string) string {
	color.Magenta("Getting embedded link...")
	resp, err := http.Get(_searchLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`<iframe src="(.+?)".*`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		panic("Link not found")
	}

	return "https:" + matches[1]
}

func GetM3u8Link(link string) string {

	ajaxURL := os.Getenv("AJAX_URL")

	// Extract the ID from the URL
	re := regexp.MustCompile(`id=(.+?)&`)
	match := re.FindStringSubmatch(link)[1]
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
	form := url.Values{}
	form.Set("id", encodedID)
	req, err := http.NewRequest("POST", ajaxURL, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	type Data struct {
		Data string `json:"data"`
	}
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	// Decode and Decrypt the response
	decodedData := Decrypt(data.Data)

	// Get the m3u8 link
	re = regexp.MustCompile(`(https.+?.m3u8)`)
	match1 := re.FindAllString(string(decodedData), -1)

	// Getting the longest m3u8 link
	longest := match1[0]
	for _, m := range match1 {
		if len(m) > len(longest) {
			longest = m
		}
	}
	m3u8 := strings.ReplaceAll(longest, `\`, ``)
	return m3u8
}

func GetMp4Link(id string) string {

	link := os.Getenv("EMBED_URL") + id
	req, err := http.NewRequest("POST", link, strings.NewReader(""))
	if err != nil {
		panic(err)
	}
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// Print the resolution
	re := regexp.MustCompile("({\"file\":.+?})")
	match := re.FindAllString(string(responseBody), -1)

	type File struct {
		File  string `json:"file"`
		Label string `json:"label"`
		Type  string `json:"type"`
	}

	resolution := make([]File, len(match))
	for i, m := range match {
		err = json.Unmarshal([]byte(m), &resolution[i])
		if err != nil {
			panic(err)
		}
		if resolution[i].Label == "720p" {
			fmt.Printf("[%d] %s (Recommended)\n", i+1, resolution[i].Label)
		} else {
			fmt.Printf("[%d] %s\n", i+1, resolution[i].Label)
		}
	}

	idx := Choose()
	return resolution[idx-1].File

}

func GetVideoLink(link string) string {
	color.Magenta("Getting video link...")
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("v/(.+?)#")
	match := re.FindSubmatch(body)
	if len(match) < 2 {
		return GetM3u8Link(link)
	} else {
		return GetMp4Link(string(match[1]))
	}
}
