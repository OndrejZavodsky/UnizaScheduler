package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Class struct {
	Start   int
	Name    string
	Room    string
	Teacher string
	Id      string
	Url     string
}

type Block struct {
	class         []Class
	possibilities []Block
}

type Calendar struct {
	classes []Block
}

func FetchScheduleBlocks(skupina string) (string, error) {
	baseURL := "https://vzdelavanie.uniza.sk/vzdelavanie/rozvrh2.php"
	reqURL := fmt.Sprintf("%s?sq=2&id=%s", baseURL, url.QueryEscape(skupina))

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML document: %w", err)
	}

	targetSelection := doc.Find("div#miestnost_bloky")
	if targetSelection.Length() == 0 {
		return "", fmt.Errorf("target div#miestnost_bloky element not found in response HTML")
	}

	htmlContent, err := targetSelection.Html()
	if err != nil {
		return "", fmt.Errorf("failed to extract inner HTML: %w", err)
	}

	return strings.TrimSpace(htmlContent), nil
}

func ParseClasses(htmlSnippet string) ([]Class, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSnippet))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML snippet: %w", err)
	}

	var classes []Class

	doc.Find("div.rozvrh_tyzden").Each(func(_ int, daySelection *goquery.Selection) {
		blockIndex := 0
		daySelection.Children().Each(func(_ int, s *goquery.Selection) {
			if s.HasClass("rozvrh_nazov") {
				return
			}
			blockIndex++
			groupLink := s.Find("a[href*='sq=4']")
			if groupLink.Length() > 0 {
				groupLink.Each(func(_ int, link *goquery.Selection) {
					href, exists := link.Attr("href")
					if !exists {
						return
					}
					classID := ""
					if parts := strings.Split(href, "id="); len(parts) > 1 {
						classID = parts[1]
					}
					trimmedID := classID
					if percentIndex := strings.Index(classID, "%"); percentIndex != -1 {
						trimmedID = classID[:percentIndex]
					}
					teacher := strings.TrimSpace(s.Find("a[href*='sq=1']").Text())
					room := strings.TrimSpace(s.Find("a[href*='sq=3']").Text())
					name := strings.TrimSpace(link.Text())

					classes = append(classes, Class{
						Start:   blockIndex,
						Name:    name,
						Room:    room,
						Teacher: teacher,
						Id:      trimmedID,
						Url:     href,
					})
				})
			}
		})
	})

	return classes, nil
}
