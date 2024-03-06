package screen

import (
	"Presto.Sonata/lib/opensearch"
	"Presto.Sonata/search/protos/serving/api/commonpb"
	"strconv"
	"strings"
)

func WorkResponseToUI(resp opensearch.WorkResponse) []*commonpb.SearchScreenSectionComponentItem {

	workItems := make([]*commonpb.SearchScreenSectionComponentItem, 0)
	for _, work := range resp.Works {
		image := work.ComposerHeader.Image.ImageUrl
		title, opusAndTitle := toTitle(work, 1)
		workItem := commonpb.SearchScreenSectionComponentItem{
			Type:            "work",
			Title:           title,
			ComposerName:    work.ComposerHeader.FullName,
			OpusAndAltTitle: opusAndTitle,
			Image: &commonpb.ApiImage{
				Url: image,
			},
			Style: "no-image",
		}
		if n, ok := work.NumberOfRecordings["us"]; ok {
			workItem.RecordingsNumber = strconv.FormatInt(int64(n), 10)
		}

		var composerOrArtist string
		if work.ComposerHeader.IsComposer {
			composerOrArtist = "composer"
		} else {
			composerOrArtist = "artist"
		}
		// generate action
		action := &commonpb.ApiAction{
			Type:       "componentScreen",
			ScreenType: "work",
			Url:        "",
			Title:      work.Title,
			Header: &commonpb.ActionHeader{
				Type:         "work",
				Title:        workItem.Title,
				Subtitle:     workItem.OpusAndAltTitle,
				ComposerName: work.ComposerHeader.FullName,
				ComposerAction: &commonpb.ApiAction{
					Type:       "componentScreen",
					ScreenType: "artist",
					Url:        "",
					Title:      work.ComposerHeader.FullName,
					Header: &commonpb.ActionHeader{
						Type:     "artist",
						Title:    work.ComposerHeader.FullName,
						Subtitle: composerOrArtist,
						Image: &commonpb.ApiImage{
							Url: work.ComposerHeader.Image.ImageUrl,
						},
						StartYear: work.ComposerHeader.StartYear,
						EndYear:   work.ComposerHeader.EndYear,
					},
					Offline: false,
				},
			},
			Offline: false,
		}
		workItem.Action = action
		workItems = append(workItems, &workItem)
	}

	return workItems
}

// return title and opusAndAltTitle
func toTitle(work *opensearch.Work, maxCatalogNumbers int) (string, string) {
	desc := make([]string, 0)
	total := len(work.CatalogNumbers)
	var catalog string
	if total > 0 {
		if total > maxCatalogNumbers {
			total = maxCatalogNumbers
		}
		catalog = strings.Join(work.CatalogNumbers[0:total], ",")
		desc = append(desc, catalog)
	}

	if work.AlternativeTitle != "" {
		desc = append(desc, work.AlternativeTitle)
	}

	var opusAndAltTitle string
	if len(desc) > 0 {
		opusAndAltTitle = strings.Join(desc, "\u00b7")
	} else {
		opusAndAltTitle = ""
	}

	var title string
	if work.Key != "" {
		title = work.ComposerHeader.ShortName + ":" + work.Title + " in " + work.Key
	} else {
		if catalog == "" {
			title = work.ComposerHeader.ShortName + work.Title
		} else {
			title = work.ComposerHeader.ShortName + ":" + work.Title + catalog
		}
	}

	return title, opusAndAltTitle
}
