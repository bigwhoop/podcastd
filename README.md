# podcastd

A super simple daemon that periodically checks RSS feeds for new audio content.


## Installation

	go get -u github.com/bigwhoop/podcastd

## Configuration

Create a `podcastd.yml` file along these lines:

	folder: "/home/phil/podcasts"
	feed_format: "%feed.title%"
	item_format: "%item.title%"
	interval: 30m
	feeds:
	  - url: http://www.nerdtalk.de/feed/podcast/
		placeholders:
		  feed.title: "Nerdtalk"
	  - url: http://www.heise.de/developer/podcast/itunes/heise-developer-podcast-softwarearchitektour.rss
	  - url: "http://www.npr.org/rss/podcast.php?id=510289"

The following placeholders are available:

- `%feed.title%` - The title of the feed
- `%item.title%` - The title of the feed
- `%item.author%` - The author of the item
- `%item.category%` - The category of the item
- `%item.date%` - The publication time of the item in format YYYY-MM-DD
- `%item.time%` - The publication time of the item in format HHMMSS

By default the values come from the actually RSS feeds and are applied to the `feed_format` and `item_format` directives.
However you can define static values using the `placeholders` map.

## Usage

	podcastd /path/to/podcastd.yml

**I haven't really tested this thing ...**
	
## License

MIT