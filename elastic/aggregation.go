/**
* @program: eutils
*
* @create: 2025-04-21 15:44
**/

package elastic

import (
	"github.com/lemonyxk/kitty/json"
	"io"
)

type Aggregation struct {
	Reader io.ReadCloser
	Error  error
}

func (a *Aggregation) All(result interface{}) error {
	if a.Error != nil {
		return a.Error
	}

	defer func() { _ = a.Reader.Close() }()

	type Agg struct {
		Took     int  `json:"took"`
		TimedOut bool `json:"timed_out"`
		Shards   struct {
			Total      int `json:"total"`
			Successful int `json:"successful"`
			Skipped    int `json:"skipped"`
			Failed     int `json:"failed"`
		} `json:"_shards"`
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			MaxScore interface{}   `json:"max_score"`
			Hits     []interface{} `json:"hits"`
		} `json:"hits"`
		Aggregations json.RawMessage `json:"aggregations"`
	}

	var agg Agg
	var err = json.NewDecoder(a.Reader).Decode(&agg)
	if err != nil {
		return err
	}

	return json.Unmarshal(agg.Aggregations, result)
}
