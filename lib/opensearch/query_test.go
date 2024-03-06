package opensearch

import (
	"reflect"
	"testing"
)

func TestMatchQuery(t *testing.T) {
	query := CreateMatchQuery("field1", "foo", 0)
	expected := MatchQuery{
		Match: map[string]matchQueryInternal{
			"field1": matchQueryInternal{
				Query: "foo",
				Boost: 0,
			},
		},
	}
	if !reflect.DeepEqual(*query, expected) {
		t.Errorf("CreateMatchQuery failed")
	}

	queryWithFuzziness := CreateMatchQueryWithFuzziness("field1", "foo", 0, 0)
	expectedQueryWithFuzziness := MatchQuery{
		Match: map[string]matchQueryInternal{
			"field1": matchQueryInternal{
				Query:     "foo",
				Fuzziness: "0",
				Boost:     0,
			},
		},
	}

	if !reflect.DeepEqual(*queryWithFuzziness, expectedQueryWithFuzziness) {
		t.Errorf("CreateMatchQueryWithFuzziness with specific fuzziness failed expected %v got %v",
			expectedQueryWithFuzziness, *queryWithFuzziness)
	}

	queryWithAutoFuzziness := CreateMatchQueryWithFuzziness("field1", "foo", 0, -1)
	expectedQueryWithAutoFuzziness := MatchQuery{
		Match: map[string]matchQueryInternal{
			"field1": matchQueryInternal{
				Query:     "foo",
				Fuzziness: "AUTO",
				Boost:     0,
			},
		},
	}

	if !reflect.DeepEqual(*queryWithAutoFuzziness, expectedQueryWithAutoFuzziness) {
		t.Errorf("CreateMatchQueryWithFuzziness for AUTO fuzziness failed")
	}
}
