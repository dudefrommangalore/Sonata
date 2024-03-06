package opensearch

import (
	"fmt"
)

type BoolQuery struct {
	Must   interface{} `json:"must,omitempty"`
	Should interface{} `json:"should,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}
type DisMax struct {
	Queries []interface{} `json:"queries,omitempty"`
}
type OrQuery struct {
	DisMax DisMax `json:"dis_max,omitempty"`
}

type matchQueryInternal struct {
	Query     string  `json:"query,omitempty"`
	Fuzziness string  `json:"fuzziness,omitempty"`
	Boost     float32 `json:"boost,omitempty"`
}

type MatchQuery struct {
	Match map[string]matchQueryInternal `json:"match,omitempty"`
}

type MatchPrefixQuery struct {
	Query     string  `json:"query,omitempty"`
	Fuzziness string  `json:"fuzziness,omitempty"`
	Boost     float32 `json:"boost,omitempty"`
	Slop      int32   `json:"query,omitempty"`
}

type FieldValueFactor struct {
	Field    string  `json:"field,omitempty"`
	Factor   float32 `json:"factor,omitempty"`
	Modifier string  `json:"modifier,omitempty"`
	Missing  float32 `json:"missing"`
}

type FieldValueFactorQuery struct {
	FieldValueFactor FieldValueFactor `json:"field_value_factor,omitempty"`
}

type FunctionQuery struct {
	Query     interface{} `json:"query,omitempty"`
	Functions interface{} `json:"functions,omitempty"`
	Boost     float32     `json:"boost,omitempty"`
	BoostMode string      `json:"boost_mode,omitempty"`
}

type Query struct {
	Bool          interface{}    `json:"bool,omitempty"`
	FunctionScore *FunctionQuery `json:"function_score,omitempty"`
}

type SearchRequest struct {
	Query          Query       `json:"query,omitempty"`
	Sort           interface{} `json:"sort,omitempty"`
	Size           int32       `json:"size,omitempty"`
	TrackScores    bool        `json:"track_scores,omitempty"`
	TrackTotalHits bool        `json:"track_total_hits,omitempty"`
}

func CreateMatchQuery(field string, query string, boost float32) *MatchQuery {
	return &MatchQuery{
		Match: map[string]matchQueryInternal{
			field: {
				Query: query,
				Boost: 0,
			},
		},
	}
}

func CreateMatchQueryWithFuzziness(field string, query string, boost float32, fuzziness int8) *MatchQuery {
	match := matchQueryInternal{
		Query: query,
		Boost: 0,
	}

	if fuzziness == -1 {
		match.Fuzziness = "AUTO"
	} else {
		match.Fuzziness = fmt.Sprintf("%d", fuzziness)
	}

	return &MatchQuery{
		Match: map[string]matchQueryInternal{
			field: match,
		},
	}
}

func CreateOrQuery(queries []interface{}) *OrQuery {
	return &OrQuery{DisMax: DisMax{Queries: queries}}
}
