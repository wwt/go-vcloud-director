package govcd

import (
	"context"
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/http"
	"net/url"
)

// ProviderVdc is the basic Provider VDC structure, contains the minimum set of attributes.
type ProviderVdc struct {
	ProviderVdc *types.ProviderVdc
	client      *Client
}

// ProviderVdcExtended is the extended Provider VDC structure, contains same attributes as ProviderVdc plus some more.
type ProviderVdcExtended struct {
	VMWProviderVdc *types.VMWProviderVdc
	client         *Client
}

func newProviderVdc(cli *Client) *ProviderVdc {
	return &ProviderVdc{
		ProviderVdc: new(types.ProviderVdc),
		client:      cli,
	}
}

func newProviderVdcExtended(cli *Client) *ProviderVdcExtended {
	return &ProviderVdcExtended{
		VMWProviderVdc: new(types.VMWProviderVdc),
		client:         cli,
	}
}

// GetProviderVdcByHref finds a Provider VDC by its HREF.
// On success, returns a pointer to the ProviderVdc structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcByHref(ctx context.Context, providerVdcHref string) (*ProviderVdc, error) {
	providerVdc := newProviderVdc(&vcdClient.Client)

	_, err := vcdClient.Client.ExecuteRequest(ctx, providerVdcHref, http.MethodGet,
		"", "error retrieving Provider VDC: %s", nil, providerVdc.ProviderVdc)
	if err != nil {
		return nil, err
	}

	return providerVdc, nil
}

// GetProviderVdcExtendedByHref finds a Provider VDC with extended attributes by its HREF.
// On success, returns a pointer to the ProviderVdcExtended structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcExtendedByHref(ctx context.Context, providerVdcHref string) (*ProviderVdcExtended, error) {
	providerVdc := newProviderVdcExtended(&vcdClient.Client)

	_, err := vcdClient.Client.ExecuteRequest(ctx, getAdminExtensionURL(providerVdcHref), http.MethodGet,
		"", "error retrieving extended Provider VDC: %s", nil, providerVdc.VMWProviderVdc)
	if err != nil {
		return nil, err
	}

	return providerVdc, nil
}

// GetProviderVdcById finds a Provider VDC by URN.
// On success, returns a pointer to the ProviderVdc structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcById(ctx context.Context, providerVdcId string) (*ProviderVdc, error) {
	providerVdcHref := vcdClient.Client.VCDHREF
	providerVdcHref.Path += "/admin/providervdc/" + extractUuid(providerVdcId)

	return vcdClient.GetProviderVdcByHref(ctx, providerVdcHref.String())
}

// GetProviderVdcExtendedById finds a Provider VDC with extended attributes by URN.
// On success, returns a pointer to the ProviderVdcExtended structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcExtendedById(ctx context.Context, providerVdcId string) (*ProviderVdcExtended, error) {
	providerVdcHref := vcdClient.Client.VCDHREF
	providerVdcHref.Path += "/admin/extension/providervdc/" + extractUuid(providerVdcId)

	return vcdClient.GetProviderVdcExtendedByHref(ctx, providerVdcHref.String())
}

// GetProviderVdcByName finds a Provider VDC by name.
// On success, returns a pointer to the ProviderVdc structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcByName(ctx context.Context, providerVdcName string) (*ProviderVdc, error) {
	providerVdc, err := getProviderVdcByName(ctx, vcdClient, providerVdcName, false)
	if err != nil {
		return nil, err
	}
	return providerVdc.(*ProviderVdc), err
}

// GetProviderVdcExtendedByName finds a Provider VDC with extended attributes by name.
// On success, returns a pointer to the ProviderVdcExtended structure and a nil error
// On failure, returns a nil pointer and an error
func (vcdClient *VCDClient) GetProviderVdcExtendedByName(ctx context.Context, providerVdcName string) (*ProviderVdcExtended, error) {
	providerVdcExtended, err := getProviderVdcByName(ctx, vcdClient, providerVdcName, true)
	if err != nil {
		return nil, err
	}
	return providerVdcExtended.(*ProviderVdcExtended), err
}

// Refresh updates the contents of the Provider VDC associated to the receiver object.
func (providerVdc *ProviderVdc) Refresh(ctx context.Context) error {
	if providerVdc.ProviderVdc.HREF == "" {
		return fmt.Errorf("cannot refresh, receiver Provider VDC is empty")
	}

	unmarshalledVdc := &types.ProviderVdc{}

	_, err := providerVdc.client.ExecuteRequest(ctx, providerVdc.ProviderVdc.HREF, http.MethodGet,
		"", "error refreshing Provider VDC: %s", nil, unmarshalledVdc)
	if err != nil {
		return err
	}

	providerVdc.ProviderVdc = unmarshalledVdc

	return nil
}

// Refresh updates the contents of the extended Provider VDC associated to the receiver object.
func (providerVdcExtended *ProviderVdcExtended) Refresh(ctx context.Context) error {
	if providerVdcExtended.VMWProviderVdc.HREF == "" {
		return fmt.Errorf("cannot refresh, receiver extended Provider VDC is empty")
	}

	unmarshalledVdc := &types.VMWProviderVdc{}

	_, err := providerVdcExtended.client.ExecuteRequest(ctx, providerVdcExtended.VMWProviderVdc.HREF, http.MethodGet,
		"", "error refreshing extended Provider VDC: %s", nil, unmarshalledVdc)
	if err != nil {
		return err
	}

	providerVdcExtended.VMWProviderVdc = unmarshalledVdc

	return nil
}

// ToProviderVdc converts the receiver ProviderVdcExtended into the subset ProviderVdc
func (providerVdcExtended *ProviderVdcExtended) ToProviderVdc(ctx context.Context) (*ProviderVdc, error) {
	providerVdcHref := providerVdcExtended.client.VCDHREF
	providerVdcHref.Path += "/admin/providervdc/" + extractUuid(providerVdcExtended.VMWProviderVdc.ID)

	providerVdc := newProviderVdc(providerVdcExtended.client)

	_, err := providerVdcExtended.client.ExecuteRequest(ctx, providerVdcHref.String(), http.MethodGet,
		"", "error retrieving Provider VDC: %s", nil, providerVdc.ProviderVdc)
	if err != nil {
		return nil, err
	}

	return providerVdc, nil
}

// getProviderVdcByName finds a Provider VDC with extension (extended=true) or without extension (extended=false) by name
// On success, returns a pointer to the ProviderVdc (extended=false) or ProviderVdcExtended (extended=true) structure and a nil error
// On failure, returns a nil pointer and an error
func getProviderVdcByName(ctx context.Context, vcdClient *VCDClient, providerVdcName string, extended bool) (interface{}, error) {
	foundProviderVdcs, err := vcdClient.QueryWithNotEncodedParams(ctx, nil,
		map[string]string{
			"type":          "providerVdc",
			"filter":        fmt.Sprintf("name==%s", url.QueryEscape(providerVdcName)),
			"filterEncoded": "true",
		})
	if err != nil {
		return nil, err
	}
	if len(foundProviderVdcs.Results.VMWProviderVdcRecord) == 0 {
		return nil, ErrorEntityNotFound
	}
	if len(foundProviderVdcs.Results.VMWProviderVdcRecord) > 1 {
		return nil, fmt.Errorf("more than one Provider VDC found with name '%s'", providerVdcName)
	}
	if extended {
		return vcdClient.GetProviderVdcExtendedByHref(ctx, foundProviderVdcs.Results.VMWProviderVdcRecord[0].HREF)
	}
	return vcdClient.GetProviderVdcByHref(ctx, foundProviderVdcs.Results.VMWProviderVdcRecord[0].HREF)
}
