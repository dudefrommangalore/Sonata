package main

import (
	"Presto.Sonata/lib/opensearch"
	"context"
	"flag"
	"fmt"
	"log"
)

func printWorks(works []*opensearch.Work) {
	for _, work := range works {
		fmt.Printf("Work id %s title %s\n", work.Id, work.Title)
	}
}

func printArtist(artists []*opensearch.Artist) {
	for _, artist := range artists {
		fmt.Printf("Artist id %s title %s\n", artist.Id, artist.FullName)
	}
}

func printPlaylist(playlists []*opensearch.Playlist) {
	for _, playlist := range playlists {
		fmt.Printf("Playlist id %s title %s\n", playlist.AppleId, playlist.Title)
	}
}

func printAlbum(albums []*opensearch.Album) {
	for _, album := range albums {
		fmt.Printf("Album id %s title %s\n", album.Id, album.Title)
	}
}

func main() {

	// index := flag.String("index", "work-20230202-134024", "Index to query")
	// id := flag.String("id", "corciolli-1968-pp50", "id to fetch from the index")

	flag.Parse()

	// Basic information for the Amazon OpenSearch Service domain
	domain := "https://vpc-presto-contentdb-fd-qa-u7syqouyfyal5j4ajrd4bbkoim.us-west-2.es.amazonaws.com" // e.g. https://my-domain.region.es.amazonaws.com
	client, err := opensearch.NewClient(domain, "442f80")

	if err != nil {
		log.Fatal("failed to create opensearch client")
	}

	bachWorks, err := client.SearchWorks(context.Background(), "bach")
	if err != nil {
		log.Fatal("failed to search works %v", err)
	} else {
		printWorks(bachWorks.Works)
	}

	fmt.Printf("executing artist query\n")
	beethArtists, err := client.SearchArtist(context.Background(), "beethoven")
	if err != nil {
		log.Fatal("failed to search artist")
	} else {
		printArtist(beethArtists.Artists)
	}

	fmt.Printf("executing playlist query\n")
	playlistResp, err := client.SearchPlaylist(context.Background(), "beethoven")
	if err != nil {
		log.Fatal("failed to search playlist")
	} else {
		printPlaylist(playlistResp.Playlists)
	}

	fmt.Printf("executing album query\n")
	albumResp, err := client.SearchAlbum(context.Background(), "tchaikovsky")
	if err != nil {
		log.Fatal("failed to search album")
	} else {
		printAlbum(albumResp.Albums)
	}

	tokenizeResponse, err := client.TokenizeQuery("live-442f80-work", "beethoven's is great")
	if err != nil {
		log.Fatal("failed to tokenize query %v", err)
	}

	fmt.Printf("tokenized output %s", tokenizeResponse)

}
