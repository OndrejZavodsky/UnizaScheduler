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
	Day     string
	Name    string
	Room    string
	Teacher string
	ID      string
	URL     string
}

type Block struct {
	Classes       []Class
	Possibilities []Block
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
		currentDay := strings.TrimSpace(daySelection.Find("div.rozvrh_nazov").Text())

		blockIndex := 0

		children := daySelection.Children()

		children.Each(func(i int, s *goquery.Selection) {
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

					rawID := ""
					if parts := strings.Split(href, "id="); len(parts) > 1 {
						rawID = parts[1]
					}

					trimmedID := rawID
					if percentIndex := strings.Index(rawID, "%"); percentIndex != -1 {
						trimmedID = rawID[:percentIndex]
					}

					teacher := strings.TrimSpace(s.Find("a[href*='sq=1']").Text())
					room := strings.TrimSpace(s.Find("a[href*='sq=3']").Text())
					name := strings.TrimSpace(link.Text())

					classObj := Class{
						Start:   blockIndex,
						Day:     currentDay,
						Name:    name,
						Room:    room,
						Teacher: teacher,
						ID:      trimmedID,
						URL:     href,
					}

					classes = append(classes, classObj)

					if i+1 < children.Length() {
						nextSibling := children.Eq(i + 1)
						classAttr, _ := nextSibling.Attr("class")

						if strings.Contains(classAttr, "-c") {
							classes = append(classes, classObj)
						}
					}
				})
			}
		})
	})

	return classes, nil
}

func TransformClassesIntoBlocks(classes []Class) []Block {
	if len(classes) == 0 {
		return nil
	}

	var blocks []Block
	var currentGroup []Class

	for i := 0; i < len(classes); i++ {
		if len(currentGroup) == 0 {
			currentGroup = append(currentGroup, classes[i])
			continue
		}

		if classes[i].Start == currentGroup[0].Start {
			currentGroup = append(currentGroup, classes[i])
		} else {
			blocks = append(blocks, Block{
				Classes: currentGroup,
			})
			currentGroup = []Class{classes[i]}
		}
	}

	if len(currentGroup) > 0 {
		blocks = append(blocks, Block{
			Classes: currentGroup,
		})
	}

	return blocks
}
