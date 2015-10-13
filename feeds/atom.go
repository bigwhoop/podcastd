package feeds

import (
	"encoding/xml"
	"errors"
	"time"
)

type AtomFeedReader struct{}

type atomFeed struct {
	Title   string          `xml:"title"`
	Entries []atomFeedEntry `xml:"entry"`
}

type atomFeedEntry struct {
	Id       string      `xml:"id"`
	Title    string      `xml:"title"`
	Category string      `xml:"category,omitempty"`
	Updated  string      `xml:"updated"`
	Link     atomLink    `xml:"link"`
	Summary  atomSummary `xml:"summary"`
	Author   atomAuthor  `xml:"author"`
}

type atomLink struct {
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	Length int64  `xml:"length,attr,omitempty"`
}

type atomAuthor struct {
	Name string `xml:"name,omitempty"`
}

type atomSummary struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

func (r AtomFeedReader) ParseXml(data []byte) (Feed, error) {
	feed := Feed{}
	atomFeed := atomFeed{}

	if err := xml.Unmarshal(data, &atomFeed); err != nil {
		return feed, err
	}

	feed.Title = atomFeed.Title
	feed.Description = ""

	for _, atomFeedEntry := range atomFeed.Entries {
		if atomFeedEntry.Link.Href == "" {
			continue
		}

		pubDate, err := parseAtomDate(atomFeedEntry.Updated)
		if err != nil {
			return feed, err
		}

		url, err := normalizeUrl(atomFeedEntry.Link.Href)
		if err != nil {
			return feed, err
		}

		feedItem := FeedItem{}
		feedItem.Title = atomFeedEntry.Title
		feedItem.Description = atomFeedEntry.Summary.Content
		feedItem.Author = atomFeedEntry.Author.Name
		feedItem.Category = atomFeedEntry.Category
		feedItem.PubDate = pubDate
		feedItem.FileUrl = url
		feedItem.FileSize = atomFeedEntry.Link.Length

		feed.AddItem(feedItem)
	}

	return feed, nil
}

var possibleAtomDateFormats = []string{
	time.RFC3339,
}

func parseAtomDate(date string) (time.Time, error) {
	for _, format := range possibleAtomDateFormats {
		if pubDate, err := time.Parse(format, date); err == nil {
			return pubDate, nil
		}
	}

	return time.Now(), errors.New("Failed parsing item date: " + date)
}
