package search

import (
	"Presto.Sonata/lib/opensearch"
	"Presto.Sonata/lib/screen"
	"Presto.Sonata/search/protos/serving/api/commonpb"
	"context"
	"flag"
	"fmt"
	"sync"

	searchpb "Presto.Sonata/search/protos"
)

var (
	opensearchUrl = flag.String("opensearch-url", "https://vpc-presto-contentdb-fd-qa-u7syqouyfyal5j4ajrd4bbkoim.us-west-2.es.amazonaws.com", "URL for open search (example, https://localhost:8090")
	indexHash     = flag.String("index-hash", "442f80", "Index to use")
)

type EntitType struct{}

type ServiceSearch struct {
	client *opensearch.OpenSearchClient
}

func NewService() (*ServiceSearch, error) {
	if *opensearchUrl == "" || *indexHash == "" {
		return nil, fmt.Errorf("please specify opensearch-url and index-hash")
	}

	if openSearchClient, err := opensearch.NewClient(*opensearchUrl, *indexHash); err != nil {
		return nil, err
	} else {
		return &ServiceSearch{
			client: openSearchClient,
		}, nil
	}
}

func (s *ServiceSearch) Search(ctx context.Context, request *searchpb.SearchRequest) (*commonpb.ApiSearchComponentScreen, error) {

	query := request.Query

	// Search across different entities
	wait := sync.WaitGroup{}
	var artistResponse *opensearch.ArtistResponse
	var playlistResponse *opensearch.PlaylistResponse
	var workResponse *opensearch.WorkResponse

	wait.Add(1)
	go func() {
		defer wait.Done()
		if resp, err := s.client.SearchArtist(context.WithValue(ctx, EntitType{}, "artist"), query); err != nil {
			fmt.Printf("artist failed {}\n", err)
		} else {
			fmt.Printf("got artist response\n")
			artistResponse = resp
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		if resp, err := s.client.SearchPlaylist(context.WithValue(ctx, EntitType{}, "playlist"), query); err != nil {
			fmt.Printf("playlist failed {}", err)
		} else {
			fmt.Printf("got playlist response\n")
			playlistResponse = resp
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		if resp, err := s.client.SearchWorks(context.WithValue(ctx, EntitType{}, "work"), query); err != nil {
			fmt.Printf("work failed {}", err)
		} else {
			fmt.Printf("got work response\n")
			workResponse = resp
		}
	}()

	wait.Wait()
	if artistResponse == nil ||
		workResponse == nil ||
		playlistResponse == nil {
		return nil, fmt.Errorf("failed, please see logs")
	}

	workItems := screen.WorkResponseToUI(*workResponse)

	return &commonpb.ApiSearchComponentScreen{
		Type:       "componentScreen",
		ScreenType: "searchResult",
		Title:      "Search",
		Header: &commonpb.ApiSearchScreenHeader{
			Type:        "search",
			Title:       "Search",
			Placeholder: "Composers, Works, and More",
			Query:       "d959",
			Action: &commonpb.ApiSearchScreenAction{
				Type:  "searchScreen",
				Url:   "/query/view/us/search?q=",
				Title: "Search",
				Query: "",
			},
			Button: &commonpb.ApiSearchButton{
				Title: "Cancel",
				Action: &commonpb.ApiSearchScreenAction{
					Type:  "searchScreen",
					Url:   "/query/view/us/search",
					Title: "Popular Searches",
					Query: "",
				},
			},
		},
		UserComponentsAction: &commonpb.ApiUserComponentsAction{
			Type: "user-components",
			Url:  "/query/user-components/us?originalQuery=d959&screen=search&trackIds%5B%5D=1452199842&trackIds%5B%5D=1452199670&trackIds%5B%5D=1452199648&trackIds%5B%5D=1452199830&trackIds%5B%5D=1452615086",
		},
		Sections: []*commonpb.SearchScreenSection{
			&commonpb.SearchScreenSection{
				Type:     "",
				Priority: "",
				Heading:  &commonpb.SearchScreenSectionHeading{},
				Components: []*commonpb.SearchScreenSectionComponent{
					&commonpb.SearchScreenSectionComponent{
						Type:  "work",
						Items: workItems,
					},
				},
			},
		},
	}, nil
}

/*
func (s *SearchService) mustEmbedUnimplementedSearchServiceServer() {
	//TODO implement me
	return
}

*/
