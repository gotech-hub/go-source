package middlewares

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"go-source/pkg/client"
	"go-source/pkg/utils"
)

var AddRegionFromCtxToHeader = []client.RequestMiddlewareFunc{addRegionFromCtxToHeader}

func addRegionFromCtxToHeader(client *resty.Client, request *resty.Request) error {
	ctx := request.Context()
	region, ok := ctx.Value(utils.KeyRegion).(string)
	if !ok {
		return fmt.Errorf("region not found in context")
	}

	request.SetHeader(utils.KeyRegion, region)
	return nil
}
