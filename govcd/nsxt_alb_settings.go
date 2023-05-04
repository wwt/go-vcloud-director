package govcd

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// GetAlbSettings retrieves NSX-T ALB settings for a particular Edge Gateway
func (egw *NsxtEdgeGateway) GetAlbSettings(ctx context.Context) (*types.NsxtAlbConfig, error) {
	client := egw.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbEdgeGateway
	apiVersion, err := client.getOpenApiHighestElevatedVersion(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, egw.EdgeGateway.ID))
	if err != nil {
		return nil, err
	}

	typeResponse := &types.NsxtAlbConfig{}
	err = client.OpenApiGetItem(ctx, apiVersion, urlRef, nil, &typeResponse, nil)
	if err != nil {
		return nil, err
	}

	return typeResponse, nil
}

// UpdateAlbSettings updates NSX-T ALB settings for a particular Edge Gateway
func (egw *NsxtEdgeGateway) UpdateAlbSettings(ctx context.Context, config *types.NsxtAlbConfig) (*types.NsxtAlbConfig, error) {
	client := egw.client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbEdgeGateway
	apiVersion, err := client.getOpenApiHighestElevatedVersion(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, egw.EdgeGateway.ID))
	if err != nil {
		return nil, err
	}

	typeResponse := &types.NsxtAlbConfig{}
	err = client.OpenApiPutItem(ctx, apiVersion, urlRef, nil, config, typeResponse, nil)
	if err != nil {
		return nil, err
	}

	return typeResponse, nil
}

// DisableAlb is a shortcut wrapping UpdateAlbSettings which disables ALB configuration
func (egw *NsxtEdgeGateway) DisableAlb(ctx context.Context) error {
	config := &types.NsxtAlbConfig{
		Enabled: false,
	}
	_, err := egw.UpdateAlbSettings(ctx, config)
	if err != nil {
		return fmt.Errorf("error disabling NSX-T ALB: %s", err)
	}

	return nil
}
