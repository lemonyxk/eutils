package types

import (
	"fmt"
)

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type UpdateResponse struct {
	Index       string `json:"_index"`
	Id          string `json:"_id"`
	Version     int    `json:"_version"`
	Result      string `json:"result"`
	Shards      Shards `json:"_shards"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Status      int    `json:"status"`
	Error       *Error `json:"error,omitempty"`
}

type CausedBy struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type Error struct {
	Type     string   `json:"type"`
	Reason   string   `json:"reason"`
	CausedBy CausedBy `json:"caused_by"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s %s", e.Type, e.Reason, e.CausedBy.Reason)
}

type RootCause struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	IndexUuid string `json:"index_uuid"`
	Index     string `json:"index"`
}

type FailedShard struct {
	Shard  int    `json:"shard"`
	Index  string `json:"index"`
	Node   string `json:"node"`
	Reason struct {
		Type      string `json:"type"`
		Reason    string `json:"reason"`
		IndexUuid string `json:"index_uuid"`
		Index     string `json:"index"`
	} `json:"reason"`
}

type QueryError struct {
	RootCause    []RootCause   `json:"root_cause"`
	Type         string        `json:"type"`
	Reason       string        `json:"reason"`
	Phase        string        `json:"phase"`
	Grouped      bool          `json:"grouped"`
	FailedShards []FailedShard `json:"failed_shards"`
}

func (e *QueryError) Error() string {
	var s = fmt.Sprintf("%s: %s - ", e.Type, e.Reason)
	for i := 0; i < len(e.RootCause); i++ {
		s += fmt.Sprintf("%s: %s", e.RootCause[i].Type, e.RootCause[i].Reason)
		if i != len(e.RootCause)-1 {
			s += " | "
		}
	}
	return s
}

type Retries struct {
	Bulk   int `json:"bulk"`
	Search int `json:"search"`
}

type UpdateByQueryResponse struct {
	Took                 int     `json:"took"`
	TimedOut             bool    `json:"timed_out"`
	Total                int     `json:"total"`
	Updated              int     `json:"updated"`
	Deleted              int     `json:"deleted"`
	Batches              int     `json:"batches"`
	VersionConflicts     int     `json:"version_conflicts"`
	Noops                int     `json:"noops"`
	Retries              Retries `json:"retries"`
	ThrottledMillis      int     `json:"throttled_millis"`
	RequestsPerSecond    float64 `json:"requests_per_second"`
	ThrottledUntilMillis int     `json:"throttled_until_millis"`
	Failures             []any   `json:"failures"`
}

type DeleteByQueryResponse struct {
	Took                 int     `json:"took"`
	TimedOut             bool    `json:"timed_out"`
	Total                int     `json:"total"`
	Deleted              int     `json:"deleted"`
	Batches              int     `json:"batches"`
	VersionConflicts     int     `json:"version_conflicts"`
	Noops                int     `json:"noops"`
	Retries              Retries `json:"retries"`
	ThrottledMillis      int     `json:"throttled_millis"`
	RequestsPerSecond    float64 `json:"requests_per_second"`
	ThrottledUntilMillis int     `json:"throttled_until_millis"`
	Failures             []any   `json:"failures"`
}

type MultiIndexResponse struct {
	Took   int    `json:"took"`
	Errors bool   `json:"errors"`
	Items  []Item `json:"items"`
}

type Item struct {
	Index UpdateResponse `json:"index"`
}
