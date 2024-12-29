package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("error decoding xml response: %w", err)
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for n, i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[n].Title = html.UnescapeString(i.Title)
		rssFeed.Channel.Item[n].Description = html.UnescapeString(i.Description)
	}

	return &rssFeed, nil
}

func printRssFeed(f *RSSFeed) {
	fmt.Printf("Title: %s\n", f.Channel.Title)
	fmt.Printf("Link: %s\n", f.Channel.Link)
	fmt.Printf("Description: %s\n", f.Channel.Description)
	fmt.Println("Items: ")
	for n, item := range f.Channel.Item {
		fmt.Printf("  %d. Title: %s\n", n+1, item.Title)
		fmt.Printf("  %d. Link: %s\n", n+1, item.Link)
		fmt.Printf("  %d. Description: %s\n", n+1, item.Description)
		fmt.Printf("  %d. PubDate: %s\n", n+1, item.PubDate)
	}
}
