package opensearch

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	awsv4 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/opensearch-project/opensearch-go/signer/aws"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type OpenSearchClient struct {
	client    *opensearch.Client
	indexHash string
}

type tokenizerOutput struct {
	Tokens []struct {
		Token       string `json:"token"`
		StartOffset int32  `json:"start_offset"`
		EndOffset   int32  `json:"end_offset"`
		Position    int32  `json:"position"`
	} `json:"tokens"`
}

type field struct {
	name   string
	weight float32
}

var (
	profile = flag.String("profile", "presto/qa", "Profile to use")
	region  = flag.String("region", "us-west-2", "AWS region to use")
)

func NewClient(endpoint string, indexHash string) (*OpenSearchClient, error) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	signer, err := aws.NewSigner(session.Options{
		Profile:            *profile,
		Config:             awsv4.Config{Region: region},
		AssumeRoleDuration: time.Duration(24) * time.Hour,
	})
	if err != nil {
		return nil, err
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{endpoint},
		Signer:    signer,
	})

	if err != nil {
		fmt.Errorf("failed to create open search client %v", err)
		return nil, err
	}

	if info, err := client.Info(); err != nil {
		return nil, err
	} else {
		fmt.Errorf("got info {}", info)
	}
	if err != nil {
		return nil, err
	}

	return &OpenSearchClient{client: client, indexHash: indexHash}, nil
}

type tokenizeQuery struct {
	Text string `json:"text"`
}

func (c *OpenSearchClient) TokenizeQuery(ctx context.Context, index string, query string) ([]string, error) {

	q := tokenizeQuery{Text: query}
	strQ, err := json.Marshal(&q)
	if err != nil {
		return nil, err
	}
	analyzeRequest := opensearchapi.IndicesAnalyzeRequest{
		Index: index,
		Body:  strings.NewReader(string(strQ)),
	}
	resp, err := analyzeRequest.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	strRes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var out tokenizerOutput
	err = json.Unmarshal(strRes, &out)
	if err != nil {
		return nil, err
	}

	fmt.Printf("tokenizer output %s\n", string(strRes))

	tokens := make([]string, 0)
	for _, tok := range out.Tokens {
		tokens = append(tokens, tok.Token)
	}

	return tokens, nil

}

func (c *OpenSearchClient) SearchWorks(ctx context.Context, query string) (*WorkResponse, error) {

	fieldsToSearch := []field{
		{name: "srchTitle.default", weight: 2.0},
		{name: "composerHeader.srchFullName.default", weight: 2.0},
	}

	fieldValueFactors := []field{
		{name: "popularity.overall", weight: 2.0},
	}

	searchRequest := constructQuery(query, fieldsToSearch, fieldValueFactors)

	if searchResp, err := c.executeQuery(ctx, searchRequest, "work"); err == nil {

		workResp := make([]*Work, 0)
		for _, hit := range searchResp.Hits.Hits {

			// convert map to string and back to struct
			if sourceStr, err := json.Marshal(hit.Source); err != nil {
				return nil, err
			} else {
				var work Work
				if err := json.Unmarshal(sourceStr, &work); err != nil {
					return nil, err
				}
				workResp = append(workResp, &work)
			}
		}
		return &WorkResponse{
			Debug: searchResp.Request,
			Works: workResp,
		}, nil

	} else {
		return nil, err
	}
}

func (c *OpenSearchClient) SearchArtist(ctx context.Context, query string) (*ArtistResponse, error) {

	fieldsToSearch := []field{
		{name: "srchFullName.default", weight: 2.0},
	}

	scoringFields := []field{
		{name: "popularity.overall", weight: 2.0},
		{name: "popularity.asComposer.overall", weight: 2.0},
	}

	searchRequest := constructQuery(query, fieldsToSearch, scoringFields)

	if searchResp, err := c.executeQuery(ctx, searchRequest, "artist"); err == nil {

		artistResp := make([]*Artist, 0)
		for _, hit := range searchResp.Hits.Hits {

			// convert map to string and back to struct
			if sourceStr, err := json.Marshal(hit.Source); err != nil {
				return nil, err
			} else {
				var artist Artist
				if err := json.Unmarshal(sourceStr, &artist); err != nil {
					return nil, err
				}
				artistResp = append(artistResp, &artist)
			}
		}
		return &ArtistResponse{
			Debug:   searchResp.Request,
			Artists: artistResp,
		}, nil
	} else {
		return nil, err
	}
}

func (c *OpenSearchClient) SearchPlaylist(ctx context.Context, query string) (*PlaylistResponse, error) {

	fieldsToSearch := []field{
		{name: "srchTitle.default", weight: 2.0},
	}

	scoringFields := []field{}

	searchRequest := constructQuery(query, fieldsToSearch, scoringFields)

	if searchResp, err := c.executeQuery(ctx, searchRequest, "playlist"); err == nil {

		playlistResp := make([]*Playlist, 0)
		for _, hit := range searchResp.Hits.Hits {

			// convert map to string and back to struct
			if sourceStr, err := json.Marshal(hit.Source); err != nil {
				return nil, err
			} else {
				var playlist Playlist
				if err := json.Unmarshal(sourceStr, &playlist); err != nil {
					return nil, err
				}
				playlistResp = append(playlistResp, &playlist)
			}
		}
		return &PlaylistResponse{
			Debug:     searchResp.Request,
			Playlists: playlistResp,
		}, nil
	} else {
		return nil, err
	}
}

func (c *OpenSearchClient) SearchAlbum(ctx context.Context, query string) (*AlbumResponse, error) {

	fieldsToSearch := []field{
		{name: "srchTitle.default", weight: 2.0},
		{name: "srchPrimary.default", weight: 2.0},
	}

	scoringFields := []field{
		{name: "popularityByContributors", weight: 2.0},
	}

	if searchResp, err := executeAndParse[Album](ctx, c, "album", query, fieldsToSearch, scoringFields); err == nil {
		return &AlbumResponse{
			Albums: searchResp,
		}, nil
	} else {
		return nil, err
	}
}

func (c *OpenSearchClient) executeQuery(ctx context.Context, query interface{}, tableName string) (*SearchResponse, error) {
	requestBodyTemplateStr, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Executing %s query for %s\n", requestBodyTemplateStr, tableName)
	search := opensearchapi.SearchRequest{
		Index: []string{"live-" + c.indexHash + "-" + tableName},
		Body:  strings.NewReader(string(requestBodyTemplateStr)),
	}
	resp, err := search.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	strResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse string response as SearchResponse
	var searchResp SearchResponse
	if err := json.Unmarshal(strResp, &searchResp); err != nil {
		return nil, err
	}
	fmt.Printf("Query for %s executed in %d ms\n", tableName, searchResp.Took)

	searchResp.Request = string(requestBodyTemplateStr)
	return &searchResp, nil

}

func constructQuery(query string, fields []field, scoringFields []field) *SearchRequest {
	queryTokens := strings.Split(query, " ")
	mustMatchQueries := make([]interface{}, 0)
	for _, token := range queryTokens {

		tokenFieldQueries := make([]interface{}, 0)
		for _, field := range fields {
			matchQuery := CreateMatchQuery(field.name, token, field.weight)
			tokenFieldQueries = append(tokenFieldQueries, matchQuery)
		}
		mustMatchQueries = append(mustMatchQueries, CreateOrQuery(tokenFieldQueries))
	}

	fieldValueFactors := make([]FieldValueFactorQuery, 0)
	for _, field := range scoringFields {
		fieldValueFactors = append(fieldValueFactors, FieldValueFactorQuery{FieldValueFactor: FieldValueFactor{
			Field:    field.name,
			Factor:   field.weight,
			Modifier: "sqrt",
			Missing:  0,
		}})
	}

	return &SearchRequest{
		Query: Query{
			FunctionScore: &FunctionQuery{
				Query: Query{
					Bool: BoolQuery{
						Must: mustMatchQueries,
					},
				},
				Functions: fieldValueFactors,
				Boost:     0.5,
				BoostMode: "",
			},
		},
		Size: 50,
	}
}

func executeAndParse[T Album | Playlist | Work | Artist](ctx context.Context, c *OpenSearchClient, entityType string, query string, fieldsToSearch []field, scoringFields []field) ([]*T, error) {
	searchRequest := constructQuery(query, fieldsToSearch, scoringFields)
	if searchResp, err := c.executeQuery(ctx, searchRequest, entityType); err == nil {
		responses := make([]*T, 0)
		for _, hit := range searchResp.Hits.Hits {
			if sourceStr, err := json.Marshal(hit.Source); err != nil {
				return nil, err
			} else {
				var entityResp T
				if err := json.Unmarshal(sourceStr, &entityResp); err != nil {
					return nil, err
				} else {
					responses = append(responses, &entityResp)
				}
			}
		}

		return responses, nil
	} else {
		return nil, err
	}
}
