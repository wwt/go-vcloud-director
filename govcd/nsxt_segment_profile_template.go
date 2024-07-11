/*
 * Copyright 2023 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

const labelNsxtSegmentProfileTemplate = "NSX-T Segment Profile Template"

// NsxtSegmentProfileTemplate contains a structure for configuring Segment Profile Templates
type NsxtSegmentProfileTemplate struct {
	NsxtSegmentProfileTemplate *types.NsxtSegmentProfileTemplate
	VCDClient                  *VCDClient
}

// wrap is a hidden helper that facilitates the usage of a generic CRUD function
//
//lint:ignore U1000 this method is used in generic functions, but annoys staticcheck
func (n NsxtSegmentProfileTemplate) wrap(inner *types.NsxtSegmentProfileTemplate) *NsxtSegmentProfileTemplate {
	n.NsxtSegmentProfileTemplate = inner
	return &n
}

// CreateSegmentProfileTemplate creates a Segment Profile Template that can later be assigned to
// global VCD configuration, Org VDC or Org VDC Network
func (vcdClient *VCDClient) CreateSegmentProfileTemplate(ctx context.Context, segmentProfileConfig *types.NsxtSegmentProfileTemplate) (*NsxtSegmentProfileTemplate, error) {
	c := crudConfig{
		endpoint:    types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentProfileTemplates,
		entityLabel: labelNsxtSegmentProfileTemplate,
	}
	outerType := NsxtSegmentProfileTemplate{VCDClient: vcdClient}
	return createOuterEntity(ctx, &vcdClient.Client, outerType, c, segmentProfileConfig)
}

// GetAllSegmentProfileTemplates retrieves all Segment Profile Templates
func (vcdClient *VCDClient) GetAllSegmentProfileTemplates(ctx context.Context, queryFilter url.Values) ([]*NsxtSegmentProfileTemplate, error) {
	c := crudConfig{
		endpoint:        types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentProfileTemplates,
		entityLabel:     labelNsxtSegmentProfileTemplate,
		queryParameters: queryFilter,
	}

	outerType := NsxtSegmentProfileTemplate{VCDClient: vcdClient}
	return getAllOuterEntities[NsxtSegmentProfileTemplate, types.NsxtSegmentProfileTemplate](ctx, &vcdClient.Client, outerType, c)
}

// GetSegmentProfileTemplateById retrieves Segment Profile Template by ID
func (vcdClient *VCDClient) GetSegmentProfileTemplateById(ctx context.Context, id string) (*NsxtSegmentProfileTemplate, error) {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentProfileTemplates,
		endpointParams: []string{id},
		entityLabel:    labelNsxtSegmentProfileTemplate,
	}

	outerType := NsxtSegmentProfileTemplate{VCDClient: vcdClient}
	return getOuterEntity[NsxtSegmentProfileTemplate, types.NsxtSegmentProfileTemplate](ctx, &vcdClient.Client, outerType, c)
}

// GetSegmentProfileTemplateByName retrieves Segment Profile Template by ID
func (vcdClient *VCDClient) GetSegmentProfileTemplateByName(ctx context.Context, name string) (*NsxtSegmentProfileTemplate, error) {
	filterByName := copyOrNewUrlValues(nil)
	filterByName = queryParameterFilterAnd(fmt.Sprintf("name==%s", name), filterByName)

	allSegmentProfileTemplates, err := vcdClient.GetAllSegmentProfileTemplates(ctx, filterByName)
	if err != nil {
		return nil, err
	}

	singleSegmentProfileTemplate, err := oneOrError("name", name, allSegmentProfileTemplates)
	if err != nil {
		return nil, err
	}

	return singleSegmentProfileTemplate, nil
}

// Update Segment Profile Template
func (spt *NsxtSegmentProfileTemplate) Update(ctx context.Context, nsxtSegmentProfileTemplateConfig *types.NsxtSegmentProfileTemplate) (*NsxtSegmentProfileTemplate, error) {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentProfileTemplates,
		endpointParams: []string{nsxtSegmentProfileTemplateConfig.ID},
		entityLabel:    labelNsxtSegmentProfileTemplate,
	}
	outerType := NsxtSegmentProfileTemplate{VCDClient: spt.VCDClient}
	return updateOuterEntity(ctx, &spt.VCDClient.Client, outerType, c, nsxtSegmentProfileTemplateConfig)
}

// Delete allows deleting NSX-T Segment Profile Template
func (spt *NsxtSegmentProfileTemplate) Delete(ctx context.Context) error {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentProfileTemplates,
		endpointParams: []string{spt.NsxtSegmentProfileTemplate.ID},
		entityLabel:    labelNsxtSegmentProfileTemplate,
	}
	return deleteEntityById(ctx, &spt.VCDClient.Client, c)
}
