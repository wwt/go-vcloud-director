/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

const labelVdcNetworkProfile = "VDC Network Profile"

// VDC Network profiles have 1:1 mapping with VDC - each VDC has an option to configure VDC Network
// Profiles. types.VdcNetworkProfile holds more information about possible configurations

// GetVdcNetworkProfile retrieves VDC Network Profile configuration
// vdc.Vdc.ID must be set and valid present
func (vdc *Vdc) GetVdcNetworkProfile(ctx context.Context) (*types.VdcNetworkProfile, error) {
	if vdc == nil || vdc.Vdc == nil || vdc.Vdc.ID == "" {
		return nil, fmt.Errorf("cannot lookup VDC Network Profile configuration without VDC ID")
	}

	return getVdcNetworkProfile(ctx, vdc.client, vdc.Vdc.ID)
}

// GetVdcNetworkProfile retrieves VDC Network Profile configuration
// vdc.Vdc.ID must be set and valid present
func (adminVdc *AdminVdc) GetVdcNetworkProfile(ctx context.Context) (*types.VdcNetworkProfile, error) {
	if adminVdc == nil || adminVdc.AdminVdc == nil || adminVdc.AdminVdc.ID == "" {
		return nil, fmt.Errorf("cannot lookup VDC Network Profile configuration without VDC ID")
	}

	return getVdcNetworkProfile(ctx, adminVdc.client, adminVdc.AdminVdc.ID)
}

// UpdateVdcNetworkProfile updates the VDC Network Profile configuration
//
// Note. Whenever updating VDC Network Profile it is required to send all fields (not only the
// changed ones) as VCD will remove other configuration. Best practice is to fetch current
// configuration of VDC Network Profile using GetVdcNetworkProfile, alter it with new values and
// submit it to UpdateVdcNetworkProfile.
func (vdc *Vdc) UpdateVdcNetworkProfile(ctx context.Context, vdcNetworkProfileConfig *types.VdcNetworkProfile) (*types.VdcNetworkProfile, error) {
	if vdc == nil || vdc.Vdc == nil || vdc.Vdc.ID == "" {
		return nil, fmt.Errorf("cannot update VDC Network Profile configuration without ID")
	}

	return updateVdcNetworkProfile(ctx, vdc.client, vdc.Vdc.ID, vdcNetworkProfileConfig)
}

// UpdateVdcNetworkProfile updates the VDC Network Profile configuration
func (adminVdc *AdminVdc) UpdateVdcNetworkProfile(ctx context.Context, vdcNetworkProfileConfig *types.VdcNetworkProfile) (*types.VdcNetworkProfile, error) {
	if adminVdc == nil || adminVdc.AdminVdc == nil || adminVdc.AdminVdc.ID == "" {
		return nil, fmt.Errorf("cannot update VDC Network Profile configuration without ID")
	}

	return updateVdcNetworkProfile(ctx, adminVdc.client, adminVdc.AdminVdc.ID, vdcNetworkProfileConfig)
}

// DeleteVdcNetworkProfile deletes VDC Network Profile Configuration
func (vdc *Vdc) DeleteVdcNetworkProfile(ctx context.Context) error {
	if vdc == nil || vdc.Vdc == nil || vdc.Vdc.ID == "" {
		return fmt.Errorf("cannot lookup VDC Network Profile without VDC ID")
	}

	return deleteVdcNetworkProfile(ctx, vdc.client, vdc.Vdc.ID)
}

// DeleteVdcNetworkProfile deletes VDC Network Profile Configuration
func (adminVdc *AdminVdc) DeleteVdcNetworkProfile(ctx context.Context) error {
	if adminVdc == nil || adminVdc.AdminVdc == nil || adminVdc.AdminVdc.ID == "" {
		return fmt.Errorf("cannot lookup VDC Network Profile without VDC ID")
	}

	return deleteVdcNetworkProfile(ctx, adminVdc.client, adminVdc.AdminVdc.ID)
}

func getVdcNetworkProfile(ctx context.Context, client *Client, vdcId string) (*types.VdcNetworkProfile, error) {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVdcNetworkProfile,
		endpointParams: []string{vdcId},
		entityLabel:    labelVdcNetworkProfile,
	}
	return getInnerEntity[types.VdcNetworkProfile](ctx, client, c)
}

func updateVdcNetworkProfile(ctx context.Context, client *Client, vdcId string, vdcNetworkProfileConfig *types.VdcNetworkProfile) (*types.VdcNetworkProfile, error) {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVdcNetworkProfile,
		endpointParams: []string{vdcId},
		entityLabel:    labelVdcNetworkProfile,
	}
	return updateInnerEntity(ctx, client, c, vdcNetworkProfileConfig)
}

func deleteVdcNetworkProfile(ctx context.Context, client *Client, vdcId string) error {
	c := crudConfig{
		endpoint:       types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVdcNetworkProfile,
		endpointParams: []string{vdcId},
		entityLabel:    labelVdcNetworkProfile,
	}
	return deleteEntityById(ctx, client, c)
}
