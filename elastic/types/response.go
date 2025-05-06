package types

import "fmt"

type Shards struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Failed     int `json:"failed,omitempty"`
}

type UpdateResponse struct {
	Index       string  `json:"_index,omitempty"`
	Id          string  `json:"_id,omitempty"`
	Version     int     `json:"_version,omitempty"`
	Result      string  `json:"result,omitempty"`
	Shards      *Shards `json:"_shards,omitempty"`
	SeqNo       int     `json:"_seq_no,omitempty"`
	PrimaryTerm int     `json:"_primary_term,omitempty"`
	Status      int     `json:"status,omitempty"`
	Error       *Error  `json:"error,omitempty"`
}

type CausedBy struct {
	Type   string `json:"type,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type Error struct {
	Type     string    `json:"type,omitempty"`
	Reason   string    `json:"reason,omitempty"`
	CausedBy *CausedBy `json:"caused_by,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s %s", e.Type, e.Reason, e.CausedBy.Reason)
}

type RootCause struct {
	Type      string `json:"type,omitempty"`
	Reason    string `json:"reason,omitempty"`
	IndexUuid string `json:"index_uuid,omitempty"`
	Index     string `json:"index,omitempty"`
}

type FailedShard struct {
	Shard  int                `json:"shard,omitempty"`
	Index  string             `json:"index,omitempty"`
	Node   string             `json:"node,omitempty"`
	Reason *FailedShardReason `json:"reason,omitempty"`
}

type FailedShardReason struct {
	Type      string `json:"type,omitempty"`
	Reason    string `json:"reason,omitempty"`
	IndexUuid string `json:"index_uuid,omitempty"`
	Index     string `json:"index,omitempty"`
}

type QueryError struct {
	RootCause    []RootCause   `json:"root_cause,omitempty"`
	Type         string        `json:"type,omitempty"`
	Reason       string        `json:"reason,omitempty"`
	Phase        string        `json:"phase,omitempty"`
	Grouped      bool          `json:"grouped,omitempty"`
	FailedShards []FailedShard `json:"failed_shards,omitempty"`
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
	Bulk   int `json:"bulk,omitempty"`
	Search int `json:"search,omitempty"`
}

type UpdateByQueryResponse struct {
	Took                 int      `json:"took,omitempty"`
	TimedOut             bool     `json:"timed_out,omitempty"`
	Total                int      `json:"total,omitempty"`
	Updated              int      `json:"updated,omitempty"`
	Deleted              int      `json:"deleted,omitempty"`
	Batches              int      `json:"batches,omitempty"`
	VersionConflicts     int      `json:"version_conflicts,omitempty"`
	Noops                int      `json:"noops,omitempty"`
	Retries              *Retries `json:"retries,omitempty"`
	ThrottledMillis      int      `json:"throttled_millis,omitempty"`
	RequestsPerSecond    float64  `json:"requests_per_second,omitempty"`
	ThrottledUntilMillis int      `json:"throttled_until_millis,omitempty"`
	Failures             []any    `json:"failures,omitempty"`
}

type DeleteByQueryResponse struct {
	Took                 int      `json:"took,omitempty"`
	TimedOut             bool     `json:"timed_out,omitempty"`
	Total                int      `json:"total,omitempty"`
	Deleted              int      `json:"deleted,omitempty"`
	Batches              int      `json:"batches,omitempty"`
	VersionConflicts     int      `json:"version_conflicts,omitempty"`
	Noops                int      `json:"noops,omitempty"`
	Retries              *Retries `json:"retries,omitempty"`
	ThrottledMillis      int      `json:"throttled_millis,omitempty"`
	RequestsPerSecond    float64  `json:"requests_per_second,omitempty"`
	ThrottledUntilMillis int      `json:"throttled_until_millis,omitempty"`
	Failures             []any    `json:"failures,omitempty"`
}

type MultiIndexResponse struct {
	Took   int    `json:"took,omitempty"`
	Errors bool   `json:"errors,omitempty"`
	Items  []Item `json:"items,omitempty"`
}

type Item struct {
	Index  *UpdateResponse `json:"index,omitempty"`
	Create *UpdateResponse `json:"create,omitempty"`
	Update *UpdateResponse `json:"update,omitempty"`
	Delete *UpdateResponse `json:"delete,omitempty"`
}
