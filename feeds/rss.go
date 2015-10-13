package feeds

import (
	"encoding/xml"
	"errors"
	"time"
)

type Rss2FeedReader struct{}

type rssFeed struct {
	Version string         `xml:"version,attr"`
	Channel rssFeedChannel `xml:"channel"`
}

type rssFeedChannel struct {
	Title       string               `xml:"title"`
	Description string               `xml:"description"`
	Items       []rssFeedChannelItem `xml:"item"`
}

type rssFeedChannelItem struct {
	Guid        string       `xml:"guid,omitempty"`
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	Author      string       `xml:"author,omitempty"`
	Category    string       `xml:"category,omitempty"`
	PubDate     string       `xml:"pubDate,omitempty"`
	Enclosure   rssEnclosure `xml:"enclosure,omitempty"`
}

type rssEnclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func (r Rss2FeedReader) ParseXml(data []byte) (Feed, error) {
	feed := Feed{}
	rssFeed := rssFeed{}

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return feed, err
	}

	feed.Title = rssFeed.Channel.Title
	feed.Description = rssFeed.Channel.Description

	for _, rssFeedItem := range rssFeed.Channel.Items {
		if rssFeedItem.Enclosure.Url == "" {
			continue
		}

		pubDate, err := parseRssDate(rssFeedItem.PubDate)
		if err != nil {
			return feed, err
		}

		url, err := normalizeUrl(rssFeedItem.Enclosure.Url)
		if err != nil {
			return feed, err
		}

		feedItem := FeedItem{}
		feedItem.Title = rssFeedItem.Title
		feedItem.Description = rssFeedItem.Description
		feedItem.Author = rssFeedItem.Author
		feedItem.Category = rssFeedItem.Category
		feedItem.PubDate = pubDate
		feedItem.FileUrl = url
		feedItem.FileSize = rssFeedItem.Enclosure.Length

		feed.AddItem(feedItem)
	}

	return feed, nil
}

var possibleRssDateFormats = []string{
	"Mon, _2 Jan 2006",
	"Mon, _2 January 2006",
	"Mon, _2 Jan 2006 15:04:05 MST",
	"Mon, _2 Jan 2006 15:04:05 -0700",
}

func parseRssDate(date string) (time.Time, error) {
	for _, format := range possibleRssDateFormats {
		if pubDate, err := time.Parse(format, date); err == nil {
			return pubDate, nil
		}
	}

	return time.Now(), errors.New("Failed parsing item date: " + date)
}
