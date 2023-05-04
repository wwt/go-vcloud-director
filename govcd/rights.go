/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */
package govcd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// getAllRights retrieves all rights. Query parameters can be supplied to perform additional
// filtering
func getAllRights(ctx context.Context, client *Client, queryParameters url.Values, additionalHeader map[string]string) ([]*types.Right, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointRights
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.Right{{}}
	err = client.OpenApiGetAllItems(ctx, minimumApiVersion, urlRef, queryParameters, &typeResponses, additionalHeader)
	if err != nil {
		return nil, err
	}

	return typeResponses, nil
}

// GetAllRights retrieves all available rights.
// Query parameters can be supplied to perform additional filtering
func (client *Client) GetAllRights(ctx context.Context, queryParameters url.Values) ([]*types.Right, error) {
	return getAllRights(ctx, client, queryParameters, nil)
}

// GetAllRights retrieves all available rights. Query parameters can be supplied to perform additional
// filtering
func (adminOrg *AdminOrg) GetAllRights(ctx context.Context, queryParameters url.Values) ([]*types.Right, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getAllRights(ctx, adminOrg.client, queryParameters, getTenantContextHeader(tenantContext))
}

// getRights retrieves rights belonging to a given Role or similar container (global role, rights bundle).
// Query parameters can be supplied to perform additional filtering
func getRights(ctx context.Context, client *Client, roleId, endpoint string, queryParameters url.Values, additionalHeader map[string]string) ([]*types.Right, error) {
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint + roleId + "/rights")
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.Right{{}}
	err = client.OpenApiGetAllItems(ctx, minimumApiVersion, urlRef, queryParameters, &typeResponses, additionalHeader)
	if err != nil {
		return nil, err
	}

	return typeResponses, nil
}

// GetRights retrieves all rights belonging to a given Role. Query parameters can be supplied to perform additional
// filtering
func (role *Role) GetRights(ctx context.Context, queryParameters url.Values) ([]*types.Right, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointRoles
	return getRights(ctx, role.client, role.Role.ID, endpoint, queryParameters, getTenantContextHeader(role.TenantContext))
}

// getRightByName retrieves a right by given name
func getRightByName(ctx context.Context, client *Client, name string, additionalHeader map[string]string) (*types.Right, error) {
	var params = url.Values{}

	slowSearch := false

	// When the right name contains commas or semicolons, the encoding is rejected by the API.
	// For this reason, when one or more commas or semicolons are present (6 occurrences in more than 300 right names)
	// we run the search brute force, by fetching all the rights, and comparing the names.
	// This problem should be fixed in 10.3
	// TODO: revisit this function after 10.3 is released
	if strings.Contains(name, ",") || strings.Contains(name, ";") {
		slowSearch = true
	} else {
		params.Set("filter", "name=="+name)
	}
	rights, err := getAllRights(ctx, client, params, additionalHeader)
	if err != nil {
		return nil, err
	}
	if len(rights) == 0 {
		return nil, ErrorEntityNotFound
	}

	if slowSearch {
		for _, right := range rights {
			if right.Name == name {
				return right, nil
			}
		}
		return nil, ErrorEntityNotFound
	}

	if len(rights) > 1 {
		return nil, fmt.Errorf("more than one right found with name '%s'", name)
	}
	return rights[0], nil
}

// GetRightByName retrieves right by given name
func (client *Client) GetRightByName(ctx context.Context, name string) (*types.Right, error) {
	return getRightByName(ctx, client, name, nil)
}

// GetRightByName retrieves right by given name
func (adminOrg *AdminOrg) GetRightByName(ctx context.Context, name string) (*types.Right, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getRightByName(ctx, adminOrg.client, name, getTenantContextHeader(tenantContext))
}

func getRightById(ctx context.Context, client *Client, id string, additionalHeader map[string]string) (*types.Right, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointRights
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty role id")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	right := &types.Right{}

	err = client.OpenApiGetItem(ctx, minimumApiVersion, urlRef, nil, right, additionalHeader)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func (client *Client) GetRightById(ctx context.Context, id string) (*types.Right, error) {
	return getRightById(ctx, client, id, nil)
}

func (adminOrg *AdminOrg) GetRightById(ctx context.Context, id string) (*types.Right, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getRightById(ctx, adminOrg.client, id, getTenantContextHeader(tenantContext))
}

// getAllRightsCategories retrieves all rights categories. Query parameters can be supplied to perform additional
// filtering
func getAllRightsCategories(ctx context.Context, client *Client, queryParameters url.Values, additionalHeader map[string]string) ([]*types.RightsCategory, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointRightsCategories
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponses := []*types.RightsCategory{{}}
	err = client.OpenApiGetAllItems(ctx, minimumApiVersion, urlRef, queryParameters, &typeResponses, additionalHeader)
	if err != nil {
		return nil, err
	}

	return typeResponses, nil
}

// GetAllRightsCategories retrieves all rights categories. Query parameters can be supplied to perform additional
// filtering
func (client *Client) GetAllRightsCategories(ctx context.Context, queryParameters url.Values) ([]*types.RightsCategory, error) {
	return getAllRightsCategories(ctx, client, queryParameters, nil)
}

// GetAllRightsCategories retrieves all rights categories. Query parameters can be supplied to perform additional
// filtering
func (adminOrg *AdminOrg) GetAllRightsCategories(ctx context.Context, queryParameters url.Values) ([]*types.RightsCategory, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getAllRightsCategories(ctx, adminOrg.client, queryParameters, getTenantContextHeader(tenantContext))
}

func getRightCategoryById(ctx context.Context, client *Client, id string, additionalHeader map[string]string) (*types.RightsCategory, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointRightsCategories
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty category id")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	rightsCategory := &types.RightsCategory{}

	err = client.OpenApiGetItem(ctx, minimumApiVersion, urlRef, nil, rightsCategory, additionalHeader)
	if err != nil {
		return nil, err
	}

	return rightsCategory, nil
}

// GetRightsCategoryById retrieves a rights category from its ID
func (adminOrg *AdminOrg) GetRightsCategoryById(ctx context.Context, id string) (*types.RightsCategory, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getRightCategoryById(ctx, adminOrg.client, id, getTenantContextHeader(tenantContext))
}

// GetRightsCategoryById retrieves a rights category from its ID
func (client *Client) GetRightsCategoryById(ctx context.Context, id string) (*types.RightsCategory, error) {
	return getRightCategoryById(ctx, client, id, nil)
}
