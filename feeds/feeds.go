package feeds

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Feed struct {
	Title       string
	Description string
	Items       []FeedItem
}

func (f *Feed) AddItem(i FeedItem) {
	f.Items = append(f.Items, i)
}

type FeedItem struct {
	Title       string
	Description string
	Author      string
	Category    string
	PubDate     time.Time
	FileUrl     string
	FileSize    int64
}

type FeedReader interface {
	ParseXml([]byte) (Feed, error)
}

func Download(url string) (Feed, error) {
	xmlData, err := downloadFile(url)
	if err != nil {
		return Feed{}, err
	}

	var feedReader FeedReader
	if isRss2(xmlData) {
		feedReader = Rss2FeedReader{}
	} else if isAtom(xmlData) {
		feedReader = AtomFeedReader{}
	} else {
		return Feed{}, errors.New("Failed detecting feed format.")
	}

	return feedReader.ParseXml(xmlData)
}

func downloadFile(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	xmlData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return xmlData, nil
}

func isAtom(xmlData []byte) bool {
	return strings.Index(string(xmlData), "xmlns:atom=\"http://www.w3.org/2005/Atom\"") > 0
}

func isRss2(xmlData []byte) bool {
	r := regexp.MustCompile("<rss.* version=\"2\\.0\"")
	return r.Match(xmlData)
}

func normalizeUrl(uri string) (string, error) {
	theUrl, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return theUrl.Scheme + "://" + theUrl.Host + theUrl.Path, nil
}
