package govcd

/*
 * Copyright 2023 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// VgpuProfile defines a vGPU profile which is fetched from vCenter
type VgpuProfile struct {
	VgpuProfile *types.VgpuProfile
	client      *Client
}

// GetAllVgpuProfiles gets all vGPU profiles that are available to VCD
func (client *VCDClient) GetAllVgpuProfiles(ctx context.Context, queryParameters url.Values) ([]*VgpuProfile, error) {
	return getAllVgpuProfiles(ctx, queryParameters, &client.Client)
}

// GetVgpuProfilesByProviderVdc gets all vGPU profiles that are available to a specific provider VDC
func (client *VCDClient) GetVgpuProfilesByProviderVdc(ctx context.Context, providerVdcUrn string) ([]*VgpuProfile, error) {
	queryParameters := url.Values{}
	queryParameters = queryParameterFilterAnd(fmt.Sprintf("pvdcId==%s", providerVdcUrn), queryParameters)
	return client.GetAllVgpuProfiles(ctx, queryParameters)
}

// GetVgpuProfileById gets a vGPU profile by ID
func (client *VCDClient) GetVgpuProfileById(ctx context.Context, vgpuProfileId string) (*VgpuProfile, error) {
	return getVgpuProfileById(ctx, vgpuProfileId, &client.Client)
}

// GetVgpuProfileByName gets a vGPU profile by name
func (client *VCDClient) GetVgpuProfileByName(ctx context.Context, vgpuProfileName string) (*VgpuProfile, error) {
	return getVgpuProfileByFilter(ctx, "name", vgpuProfileName, &client.Client)
}

// GetVgpuProfileByTenantFacingName gets a vGPU profile by its tenant facing name
func (client *VCDClient) GetVgpuProfileByTenantFacingName(ctx context.Context, tenantFacingName string) (*VgpuProfile, error) {
	return getVgpuProfileByFilter(ctx, "tenantFacingName", tenantFacingName, &client.Client)
}

// Update updates a vGPU profile with new parameters
func (profile *VgpuProfile) Update(ctx context.Context, newProfile *types.VgpuProfile) error {
	client := profile.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVgpuProfile
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, "/", profile.VgpuProfile.Id)
	if err != nil {
		return err
	}

	err = client.OpenApiPutItemSync(ctx, minimumApiVersion, urlRef, nil, newProfile, nil, nil)
	if err != nil {
		return err
	}

	// We need to refresh here, as PUT returns the original struct instead of the updated one
	err = profile.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Refresh updates the current state of the vGPU profile
func (profile *VgpuProfile) Refresh(ctx context.Context) error {
	var err error
	newProfile, err := getVgpuProfileById(ctx, profile.VgpuProfile.Id, profile.client)
	if err != nil {
		return err
	}
	profile.VgpuProfile = newProfile.VgpuProfile

	return nil
}

func getVgpuProfileByFilter(ctx context.Context, filter, filterValue string, client *Client) (*VgpuProfile, error) {
	queryParameters := url.Values{}
	queryParameters = queryParameterFilterAnd(fmt.Sprintf("%s==%s", filter, filterValue), queryParameters)
	vgpuProfiles, err := getAllVgpuProfiles(ctx, queryParameters, client)
	if err != nil {
		return nil, err
	}

	vgpuProfile, err := oneOrError(filter, filterValue, vgpuProfiles)
	if err != nil {
		return nil, err
	}

	return vgpuProfile, nil
}

func getVgpuProfileById(ctx context.Context, vgpuProfileId string, client *Client) (*VgpuProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVgpuProfile
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, "/", vgpuProfileId)
	if err != nil {
		return nil, err
	}

	profile := &VgpuProfile{
		client: client,
	}
	err = client.OpenApiGetItem(ctx, minimumApiVersion, urlRef, nil, &profile.VgpuProfile, nil)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func getAllVgpuProfiles(ctx context.Context, queryParameters url.Values, client *Client) ([]*VgpuProfile, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointVgpuProfile
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	responses := []*types.VgpuProfile{{}}

	err = client.OpenApiGetAllItems(ctx, minimumApiVersion, urlRef, queryParameters, &responses, nil)
	if err != nil {
		return nil, err
	}

	wrappedVgpuProfiles := make([]*VgpuProfile, len(responses))
	for index, response := range responses {
		wrappedVgpuProfile := &VgpuProfile{
			client:      client,
			VgpuProfile: response,
		}
		wrappedVgpuProfiles[index] = wrappedVgpuProfile
	}

	return wrappedVgpuProfiles, nil
}
