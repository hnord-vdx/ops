//go:build digitalocean || do || !onlyprovider

package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/nanovms/ops/lepton"
)

// GetVPC returns a vpc by vpc name and zone
func (do *DigitalOcean) GetVPC(ctx *lepton.Context, vpcName string) (*godo.VPC, error) {

	if vpcName == "" {
		return nil, nil
	}

	page := 1
	var vpc *godo.VPC

	for {
		opts := &godo.ListOptions{
			Page:    page,
			PerPage: 200, // max allowed by DO
		}

		vpcs, _, err := do.Client.VPCs.List(context.TODO(), opts)
		if err != nil {
			return nil, err
		}
		if len(vpcs) == 0 {
			break
		}

		for _, v := range vpcs {
			if v.Name == vpcName {
				vpc = v
				break
			}
		}
		page++
	}

	if vpc == nil {
		ctx.Logger().Debugf("no vpcs with name %s found", vpcName)
		return nil, fmt.Errorf("vpc %s not found", vpcName)
	}

	return vpc, nil
}
