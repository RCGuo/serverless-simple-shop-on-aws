package models

import "encoding/json"

type SearchResponse struct {
	Took     int                   `json:"_took"`
	TimedOut bool                  `json:"timed_out"`
	Shards   *SearchResponseShards `json:"shards,omitempty"`
	Hits     *SearchResponseHits   `json:"hits,omitempty"`
}

type SearchResponseShards struct {
	Total      int
	Successful int
	Skipped    int
	Failed     int
}

type SearchResponseHits struct {
	Total    SearchResponseHitsTotal
	MaxScore json.Number         `json:"max_score"`
	Hits     []SearchResponseHit `json:"hits,omitempty"`
}

type SearchResponseHitsTotal struct {
	Value    int
	Relation string
}

type SearchResponseHit struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	ID     string          `json:"_id"`
	Score  json.Number     `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

type SearchResponseError struct {
	Error  SearchError `json:"error"`
	Status int         `json:"status"`
}

type SearchError struct {
	RootCause []ErrorCause `json:"root_cause"`
	ErrorCause
}

type ErrorCause struct {
	Type         string `json:"type"`
	Reason       string `json:"reason"`
	Index        string `json:"index"`
	ResourceId   string `json:"resource.id"`
	ResourceType string `json:"resource.type"`
	IndexUUId    string `json:"index_uuid"`
}