/**
* @program: engine
*
* @create: 2025-04-19 14:11
**/

package elastic

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"io"
)

func GetClusterState(client *elasticsearch.Client) ([]byte, error) {

	// GET /_cluster/state/metadata?filter_path=metadata.stored_scripts

	var req = esapi.ClusterStateRequest{
		//FilterPath: []string{"metadata.stored_scripts"},
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bts, nil
}
