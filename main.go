package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

func main() {
	baseURL := "https://www.rvp.co.th"
	pageURL := "https://www.rvp.co.th/rvpmedia2.php"

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1.0; +https://go-colly.org/)"),
	)

	c.OnHTML("div.home-spotlight", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			href := el.Attr("href")
			alt := el.ChildAttr("img", "alt")
			imgSrc := el.ChildAttr("img", "src")
			title := el.ChildText("h3")

			imgURL, err := url.Parse(imgSrc)
			if err != nil {
				log.Printf("Failed to parse image URL: %s, error: %v\n", imgSrc, err)
				return
			}

			if !imgURL.IsAbs() {
				imgURL = resolveURL(baseURL, imgSrc)
			}
			fmt.Printf("Link: %s\n", href)
			fmt.Printf("Image Alt: %s\n", alt)
			fmt.Printf("Image Src: %s\n", imgURL.String())
			fmt.Printf("Title: %s\n", title)
			fmt.Println("-------------")

			imagePath := filepath.Join("images", filepath.Base(imgURL.Path))

			err = downloadFile(imgURL.String(), imagePath)
			if err != nil {
				log.Printf("Failed to download image: %s, error: %v\n", imgURL.String(), err)
			}
		})
	})

	err := c.Visit(pageURL)
	if err != nil {
		log.Fatal(err)
	}
}

func resolveURL(base, rel string) *url.URL {
	baseURL, _ := url.Parse(base)
	relURL, _ := url.Parse(rel)
	return baseURL.ResolveReference(relURL)
}

func downloadFile(url, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
