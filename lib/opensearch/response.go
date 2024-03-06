package opensearch

// The response returned by OpenSearch (generic response plus the entity specific response)

type SearchHit struct {
	Id     string      `json:"_id,omitempty"`
	Score  float32     `json:"_score,omitempty"`
	Source interface{} `json:"_source,omitempty"` // source can be any of the entity type. Need to parse this with right struct type once the search result is parsed
}

type SearchHits struct {
	Hits     []SearchHit `json:"hits,omitempty"`
	MaxScore float32     `json:"max_score,omitempty"`
	Total    struct {
		Value float32 `json:"value,omitempty"`
	} `json:"total,omitempty"`
}

type SearchResponse struct {
	Hits     SearchHits `json:"hits,omitempty"`
	TimedOut bool       `json:"timed_out,omitempty"`
	Took     int32      `json:"took,omitempty"`
	Request  string     `json:"request,omitempty"`
}

type WorkResponse struct {
	Debug string  `json:"debug,omitempty"`
	Works []*Work `json:"works,omitempty"`
}

type Work struct {
	Id               string `json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	AlternativeTitle string `json:"alternativeTitle,omitempty"`
	Key              string `json:"key,omitempty"`
	ComposerHeader   struct {
		Id        string `json:"id,omitempty"`
		AdamId    string `json:"adamId,omitempty"`
		ShortName string `json:"shortName,omitempty"`
		FullName  string `json:"fullName,omitempty"`
		Image     struct {
			Type     string `json:"type,omitempty"`
			ImageUrl string `json:"imageUrl,omitempty"`
		} `json:"image,omitempty"`
		IsComposer bool   `json:"isComposer,omitempty"`
		StartYear  string `json:"startYear,omitempty"`
		EndYear    string `json:"endYear,omitempty"`
	} `json:"composerHeader,omitempty"`
	CatalogNumbers     []string         `json:"catalogNumbers,omitempty"`
	NumberOfRecordings map[string]int32 `json:"numberOfRecordings,omitempty"`
	Popularity         struct {
		Overall float32 `json:"overall,omitempty"`
	} `json:"popularity,omitempty"`
	Description []string `json:"description,omitempty"`
}

type ArtistResponse struct {
	Debug   string    `json:"debug,omitempty"`
	Artists []*Artist `json:"artists,omitempty"`
}

type Artist struct {
	Id              string   `json:"id,omitempty"`
	FullName        string   `json:"fullName,omitempty"`
	ShortName       string   `json:"shortName,omitempty"`
	RepertoireRoles []string `json:"repertoireRoles,omitempty"`
	Popularity      struct {
		Overall    float32 `json:"overall,omitempty"`
		AsComposer struct {
			Overall float32 `json:"overall,omitempty"`
		} `json:"asComposer,omitempty"`
	} `json:"popularity,omitempty"`
}

type Playlist struct {
	AppleId     string `json:"appleId,omitempty"`
	Title       string `json:"title,omitempty"`
	SubTitle    string `json:"subtitle,omitempty"`
	Description string `json:"description,omitempty"`
	CoverImage  struct {
		ImageUrl string `json:"imageUrl,omitempty"`
		Locales  struct {
			ImageUrl map[string]string `json:"imageUrl,omitempty"`
		} `json:"locales,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"coverImage,omitempty"`
	Children []struct {
		Type                   string   `json:"type,omitempty"`
		Id                     string   `json:"id,omitempty"`
		UnsupportedStorefronts []string `json:"unsupportedStorefronts,omitempty"`
	} `json:"children,omitempty"`
	Metadata struct {
		Id    string `json:"id,omitempty"`
		Type  string `json:"type,omitempty"`
		Space string `json:"space,omitempty"`
	} `json:"metadata,omitempty"`
}

type PlaylistResponse struct {
	Debug     string      `json:"debug,omitempty"`
	Playlists []*Playlist `json:"artists,omitempty"`
}

type ArtistHeader struct {
	Id           string   `json:"id,omitempty"`
	AdamId       string   `json:"adamId,omitempty"`
	FullName     string   `json:"fullName,omitempty"`
	ShortName    string   `json:"shortName,omitempty"`
	StartYear    string   `json:"startYear,omitempty"`
	EndYear      string   `json:"endYear,omitempty"`
	DisplayRoles []string `json:"displayRoles,omitempty"`
	Locales      struct {
		DisplayRoles map[string]string `json:"displayRoles,omitempty"`
	} `json:"locales,omitempty"`
}

type Album struct {
	Id                string `json:"id,omitempty"`
	Title             string `json:"title,omitempty"`
	DisplayArtistName string `json:"DisplayArtistName"`
	DisplayArtists    []struct {
		ContributorId string       `json:"contributorId,omitempty"`
		ArtistHeader  ArtistHeader `json:"artistHeader,omitempty"`
	} `json:"DisplayArtists,omitempty"`
	Composers []struct {
		ContributorId string       `json:"contributorId,omitempty"`
		ArtistHeader  ArtistHeader `json:"artistHeader,omitempty"`
	} `json:"composers,omitempty"`
	Artists []struct {
		ContributorId string       `json:"contributorId,omitempty"`
		ArtistHeader  ArtistHeader `json:"artistHeader,omitempty"`
	} `json:"artists,omitempty"`
	ReleaseTitle string `json:"releaseTitle,omitempty"`
	Children     []struct {
		Type string `json:"type,omitempty"`
		Id   string `json:"id,omitempty"`
	} `json:"children,omitempty"`
}

type AlbumResponse struct {
	Debug  string   `json:"debug,omitempty"`
	Albums []*Album `json:"albums,omitempty"`
}
