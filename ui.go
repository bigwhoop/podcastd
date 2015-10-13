package main

import (
	"fmt"
	"github.com/gizak/termui"
	"github.com/nareix/curl"
	"os"
	"sort"
	"time"
)

var (
	ui Ui
)

func NewUi() Ui {
	headerWidget := termui.NewPar(fmt.Sprintf("podcastd v%s\nCopyright 2015 Philippe Gerber\nhttps://github.com/bigwhoop/podcastd", VERSION))
	headerWidget.Height = 5
	headerWidget.HasBorder = false
	headerWidget.PaddingTop = 1
	headerWidget.PaddingBottom = 1
	headerWidget.PaddingLeft = 1

	infoWidget := termui.NewPar("")
	infoWidget.HasBorder = false
	infoWidget.Text = fmt.Sprintf("Press 'q' to quit")

	feedsWidget := termui.NewList()
	feedsWidget.Border.Label = "Feeds"

	return Ui{
		termui.TermWidth(),
		headerWidget,
		infoWidget,
		feedsWidget,
		make(map[string]bool, 0),
		make([]*termui.Gauge, 0),
		make(map[string]curl.ProgressStatus, 0),
	}
}

type Ui struct {
	gridWidth             int
	headerWidget          *termui.Par
	infoWidget            *termui.Par
	feedsWidget           *termui.List
	feeds                 map[string]bool
	activeDownloadWidgets []*termui.Gauge
	activeDownloads       map[string]curl.ProgressStatus
}

func (u Ui) refresh() {
	grid := termui.NewGrid(
		termui.NewRow(
			termui.NewCol(9, 0, u.headerWidget),
			termui.NewCol(3, 0, u.infoWidget),
		),
		termui.NewRow(
			termui.NewCol(12, 0, u.feedsWidget),
		),
	)

	for _, widget := range u.activeDownloadWidgets {
		grid.AddRows(
			termui.NewRow(
				termui.NewCol(12, 0, widget),
			),
		)
	}

	grid.Width = u.gridWidth
	grid.Align()
	termui.Render(grid)
}

func (u Ui) SetFeedStatus(uri string, loading bool) {
	u.feeds[uri] = loading
	u.refreshFeedsWidget()
}

func (u *Ui) SetActiveDownload(uri string, status curl.ProgressStatus) {
	u.activeDownloads[uri] = status
	u.refreshActiveDownloadWidgets()
}

func (u *Ui) RemoveActiveDownload(uri string) {
	delete(u.activeDownloads, uri)
	u.refreshActiveDownloadWidgets()
}

func (u *Ui) refreshFeedsWidget() {
	uris := make([]string, 0, len(u.feeds))
	for uri, isLoading := range u.feeds {
		if isLoading {
			uri += " (loading)"
		}
		uris = append(uris, uri)
	}
	sort.Strings(uris)

	u.feedsWidget.Items = uris
	u.feedsWidget.Height = len(uris) + 2
}

func (u *Ui) refreshActiveDownloadWidgets() {
	u.activeDownloadWidgets = make([]*termui.Gauge, 0, len(u.activeDownloads))

	uris := make([]string, 0, len(u.activeDownloads))
	for uri, _ := range u.activeDownloads {
		uris = append(uris, uri)
	}
	sort.Strings(uris)

	for _, uri := range uris {
		progress := u.activeDownloads[uri]

		widget := termui.NewGauge()
		widget.Height = 3
		widget.Percent = int(progress.Percent * 100)
		widget.Border.Label = "Downloading: " + uri
		widget.Label = fmt.Sprintf("{{percent}}%% (%s)", curl.PrettySpeedString(progress.Speed))

		u.activeDownloadWidgets = append(u.activeDownloadWidgets, widget)
	}
}

func init() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	termui.UseTheme("default")

	ui = NewUi()

	refreshTicker := time.NewTicker(time.Millisecond * 50)
	evt := termui.EventCh()

	ui.refresh()

	go func() {
		for {
			select {
			case e := <-evt:
				if e.Type == termui.EventKey && e.Ch == 'q' {
					os.Exit(1)
				}
				if e.Type == termui.EventResize {
					ui.gridWidth = termui.TermWidth()
					ui.refresh()
				}
			case <-refreshTicker.C:
				ui.refresh()
			}
		}
	}()
}
