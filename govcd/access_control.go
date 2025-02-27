/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// orgInfoCache is a cache to save org information, avoid repeated calls to compute the same result.
// The keys to this map are the requesting objects IDs.
var orgInfoCache = make(map[string]*TenantContext)

// GetAccessControl retrieves the access control information for the requested entity
func (client Client) GetAccessControl(ctx context.Context, href, entityType, entityName string, headerValues map[string]string) (*types.ControlAccessParams, error) {

	href += "/controlAccess"
	var controlAccess types.ControlAccessParams

	acUrl, err := url.ParseRequestURI(href)
	if err != nil {
		return nil, fmt.Errorf("[client.GetAccessControl] error parsing HREF %s: %s", href, err)
	}
	var additionalHeader = make(http.Header)

	if len(headerValues) > 0 {
		for k, v := range headerValues {
			additionalHeader.Add(k, v)
		}
	}
	req := client.newRequest(
		ctx,
		nil,               // params
		nil,               // notEncodedParams
		http.MethodGet,    // method
		*acUrl,            // reqUrl
		nil,               // body
		client.APIVersion, // apiVersion
		additionalHeader,  // additionalHeader
	)

	resp, err := checkResp(client.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("[client.GetAccessControl] error checking response to request %s: %s", href, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("[client.GetAccessControl] nil response received")
	}
	if err = decodeBody(types.BodyTypeXML, resp, &controlAccess); err != nil {
		return nil, fmt.Errorf("[client.GetAccessControl] error decoding response: %s", err)
	}

	return &controlAccess, nil
}

// SetAccessControl changes the access control information for this entity
// There are two ways of setting the access:
// with accessControl.IsSharedToEveryone = true we give access to everyone
// with accessControl.IsSharedToEveryone = false, accessControl.AccessSettings defines which subjects can access the vApp
// For each setting we must provide:
// * The subject (HREF and Type are mandatory)
// * The access level (one of ReadOnly, Change, FullControl)
func (client *Client) SetAccessControl(ctx context.Context, accessControl *types.ControlAccessParams, href, entityType, entityName string, headerValues map[string]string) error {
	return client.setAccessControlWithHttpMethod(ctx, http.MethodPost, accessControl, href, entityType, entityName, headerValues)
}

// setAccessControlWithMethod is the same as Client.SetAccessControl but allowing passing a different HTTP method.
// This method has been created since VDC accessControl endpoint works with PUT and SetAccessControl method worked
// exclusively with POST. This private method gives the flexibility to use both POST and PUT passing it as httpMethod parameter.
func (client *Client) setAccessControlWithHttpMethod(ctx context.Context, httpMethod string, accessControl *types.ControlAccessParams, href, entityType, entityName string, headerValues map[string]string) error {
	href += "/action/controlAccess"
	// Make sure that subjects in the setting list are used only once
	if accessControl.AccessSettings != nil && len(accessControl.AccessSettings.AccessSetting) > 0 {
		if accessControl.IsSharedToEveryone {
			return fmt.Errorf("[client.SetAccessControl] can't set IsSharedToEveryone and AccessSettings at the same time for %s %s (%s)", entityType, entityName, href)
		}
		var used = make(map[string]bool)
		for _, setting := range accessControl.AccessSettings.AccessSetting {
			_, seen := used[setting.Subject.HREF]
			if seen {
				return fmt.Errorf("[client.SetAccessControl] subject %s (%s) used more than once", setting.Subject.Name, setting.Subject.HREF)
			}
			used[setting.Subject.HREF] = true
			if setting.Subject.Type == "" {
				return fmt.Errorf("[client.SetAccessControl] subject %s (%s) has no type defined", setting.Subject.Name, setting.Subject.HREF)
			}
		}
	}

	accessControl.Xmlns = types.XMLNamespaceVCloud
	queryUrl, err := url.ParseRequestURI(href)
	if err != nil {
		return fmt.Errorf("[client.SetAccessControl] error parsing HREF %s: %s", href, err)
	}

	var header = make(http.Header)
	if len(headerValues) > 0 {
		for k, v := range headerValues {
			header.Add(k, v)
		}
	}

	marshaledXml, err := xml.MarshalIndent(accessControl, "  ", "    ")
	if err != nil {
		return fmt.Errorf("[client.SetAccessControl] error marshalling xml data: %s", err)
	}
	body := bytes.NewBufferString(xml.Header + string(marshaledXml))

	req := client.newRequest(
		ctx,
		nil,               // params
		nil,               // notEncodedParams
		httpMethod,        // method
		*queryUrl,         // reqUrl
		body,              // body
		client.APIVersion, // apiVersion
		header,            // additionalHeader
	)

	resp, err := checkResp(client.Http.Do(req))

	if err != nil {
		return fmt.Errorf("[client.SetAccessControl] error checking response to HREF %s: %s", href, err)
	}
	if resp == nil {
		return fmt.Errorf("[client.SetAccessControl] nil response received")
	}
	_, err = checkResp(resp, err)
	return err
}

// GetAccessControl retrieves the access control information for this vApp
func (vapp VApp) GetAccessControl(ctx context.Context, useTenantContext bool) (*types.ControlAccessParams, error) {

	if vapp.VApp.HREF == "" {
		return nil, fmt.Errorf("vApp HREF is empty")
	}
	// if useTenantContext is false, we use an empty header (= default behavior)
	// if it is true, we use a header populated with tenant context values
	accessControlHeader, err := vapp.getAccessControlHeader(ctx, useTenantContext)
	if err != nil {
		return nil, err
	}
	return vapp.client.GetAccessControl(ctx, vapp.VApp.HREF, "vApp", vapp.VApp.Name, accessControlHeader)
}

// SetAccessControl changes the access control information for this vApp
func (vapp VApp) SetAccessControl(ctx context.Context, accessControl *types.ControlAccessParams, useTenantContext bool) error {

	if vapp.VApp.HREF == "" {
		return fmt.Errorf("vApp HREF is empty")
	}

	// if useTenantContext is false, we use an empty header (= default behavior)
	// if it is true, we use a header populated with tenant context values
	accessControlHeader, err := vapp.getAccessControlHeader(ctx, useTenantContext)
	if err != nil {
		return err
	}
	return vapp.client.SetAccessControl(ctx, accessControl, vapp.VApp.HREF, "vApp", vapp.VApp.Name, accessControlHeader)

}

// RemoveAccessControl is a shortcut to SetAccessControl with all access disabled
func (vapp VApp) RemoveAccessControl(ctx context.Context, useTenantContext bool) error {
	return vapp.SetAccessControl(ctx, &types.ControlAccessParams{IsSharedToEveryone: false}, useTenantContext)
}

// IsShared shows whether a vApp is shared or not, regardless of the number of subjects sharing it
func (vapp VApp) IsShared(ctx context.Context, useTenantContext bool) (bool, error) {
	settings, err := vapp.GetAccessControl(ctx, useTenantContext)
	if err != nil {
		return false, err
	}
	if settings.IsSharedToEveryone {
		return true, nil
	}
	return settings.AccessSettings != nil, nil
}

// GetAccessControl retrieves the access control information for this catalog
func (adminCatalog AdminCatalog) GetAccessControl(ctx context.Context, useTenantContext bool) (*types.ControlAccessParams, error) {

	if adminCatalog.AdminCatalog.HREF == "" {
		return nil, fmt.Errorf("catalog HREF is empty")
	}
	href := strings.Replace(adminCatalog.AdminCatalog.HREF, "/admin/", "/", 1)

	// if useTenantContext is false, we use an empty header (= default behavior)
	// if it is true, we use a header populated with tenant context values
	accessControlHeader, err := adminCatalog.getAccessControlHeader(useTenantContext)
	if err != nil {
		return nil, err
	}
	return adminCatalog.client.GetAccessControl(ctx, href, "catalog", adminCatalog.AdminCatalog.Name, accessControlHeader)
}

// SetAccessControl changes the access control information for this catalog
func (adminCatalog AdminCatalog) SetAccessControl(ctx context.Context, accessControl *types.ControlAccessParams, useTenantContext bool) error {

	if adminCatalog.AdminCatalog.HREF == "" {
		return fmt.Errorf("catalog HREF is empty")
	}
	href := strings.Replace(adminCatalog.AdminCatalog.HREF, "/admin/", "/", 1)

	// if useTenantContext is false, we use an empty header (= default behavior)
	// if it is true, we use a header populated with tenant context values
	accessControlHeader, err := adminCatalog.getAccessControlHeader(useTenantContext)
	if err != nil {
		return err
	}
	return adminCatalog.client.SetAccessControl(ctx, accessControl, href, "catalog", adminCatalog.AdminCatalog.Name, accessControlHeader)
}

// RemoveAccessControl is a shortcut to SetAccessControl with all access disabled
func (adminCatalog AdminCatalog) RemoveAccessControl(ctx context.Context, useTenantContext bool) error {
	return adminCatalog.SetAccessControl(ctx, &types.ControlAccessParams{IsSharedToEveryone: false}, useTenantContext)
}

// IsShared shows whether a catalog is shared or not, regardless of the number of subjects sharing it
func (adminCatalog AdminCatalog) IsShared(ctx context.Context, useTenantContext bool) (bool, error) {
	settings, err := adminCatalog.GetAccessControl(ctx, useTenantContext)
	if err != nil {
		return false, err
	}
	if settings.IsSharedToEveryone {
		return true, nil
	}
	return settings.AccessSettings != nil, nil
}

// GetVappAccessControl is a convenience method to retrieve access control for a vApp
// from a VDC.
// The input variable vappIdentifier can be either the vApp name or its ID
func (vdc *Vdc) GetVappAccessControl(ctx context.Context, vappIdentifier string, useTenantContext bool) (*types.ControlAccessParams, error) {
	vapp, err := vdc.GetVAppByNameOrId(ctx, vappIdentifier, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vApp %s: %s", vappIdentifier, err)
	}
	return vapp.GetAccessControl(ctx, useTenantContext)
}

// GetCatalogAccessControl is a convenience method to retrieve access control for a catalog
// from an organization.
// The input variable catalogIdentifier can be either the catalog name or its ID
func (org *AdminOrg) GetCatalogAccessControl(ctx context.Context, catalogIdentifier string, useTenantContext bool) (*types.ControlAccessParams, error) {
	catalog, err := org.GetAdminCatalogByNameOrId(ctx, catalogIdentifier, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog %s: %s", catalogIdentifier, err)
	}
	return catalog.GetAccessControl(ctx, useTenantContext)
}

// GetCatalogAccessControl is a convenience method to retrieve access control for a catalog
// from an organization.
// The input variable catalogIdentifier can be either the catalog name or its ID
func (org *Org) GetCatalogAccessControl(ctx context.Context, catalogIdentifier string, useTenantContext bool) (*types.ControlAccessParams, error) {
	catalog, err := org.GetCatalogByNameOrId(ctx, catalogIdentifier, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog %s: %s", catalogIdentifier, err)
	}
	return catalog.GetAccessControl(ctx, useTenantContext)
}

// GetAccessControl retrieves the access control information for this catalog
func (catalog Catalog) GetAccessControl(ctx context.Context, useTenantContext bool) (*types.ControlAccessParams, error) {

	if catalog.Catalog.HREF == "" {
		return nil, fmt.Errorf("catalog HREF is empty")
	}
	href := strings.Replace(catalog.Catalog.HREF, "/admin/", "/", 1)
	accessControlHeader, err := catalog.getAccessControlHeader(useTenantContext)
	if err != nil {
		return nil, err
	}
	return catalog.client.GetAccessControl(ctx, href, "catalog", catalog.Catalog.Name, accessControlHeader)
}

// SetAccessControl changes the access control information for this catalog
func (catalog Catalog) SetAccessControl(ctx context.Context, accessControl *types.ControlAccessParams, useTenantContext bool) error {

	if catalog.Catalog.HREF == "" {
		return fmt.Errorf("catalog HREF is empty")
	}

	href := strings.Replace(catalog.Catalog.HREF, "/admin/", "/", 1)

	// if useTenantContext is false, we use an empty header (= default behavior)
	// if it is true, we use a header populated with tenant context values
	accessControlHeader, err := catalog.getAccessControlHeader(useTenantContext)
	if err != nil {
		return err
	}
	return catalog.client.SetAccessControl(ctx, accessControl, href, "catalog", catalog.Catalog.Name, accessControlHeader)
}

// RemoveAccessControl is a shortcut to SetAccessControl with all access disabled
func (catalog Catalog) RemoveAccessControl(ctx context.Context, useTenantContext bool) error {
	return catalog.SetAccessControl(ctx, &types.ControlAccessParams{IsSharedToEveryone: false}, useTenantContext)
}

// IsShared shows whether a catalog is shared or not, regardless of the number of subjects sharing it
func (catalog Catalog) IsShared(ctx context.Context, useTenantContext bool) (bool, error) {
	settings, err := catalog.GetAccessControl(ctx, useTenantContext)
	if err != nil {
		return false, err
	}
	if settings.IsSharedToEveryone {
		return true, nil
	}
	return settings.AccessSettings != nil, nil
}

// getAccessControlHeader builds the data needed to set the header when tenant context is required.
// If useTenantContext is false, it returns an empty map.
// Otherwise, it finds the Org ID and name (going up in the hierarchy through the VDC)
// and creates the header data
func (vapp *VApp) getAccessControlHeader(ctx context.Context, useTenantContext bool) (map[string]string, error) {
	if !useTenantContext {
		return map[string]string{}, nil
	}
	orgInfo, err := vapp.getOrgInfo(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]string{types.HeaderTenantContext: orgInfo.OrgId, types.HeaderAuthContext: orgInfo.OrgName}, nil
}

// getAccessControlHeader builds the data needed to set the header when tenant context is required.
// If useTenantContext is false, it returns an empty map.
// Otherwise, it finds the Org ID and name and creates the header data
func (catalog *Catalog) getAccessControlHeader(useTenantContext bool) (map[string]string, error) {
	if !useTenantContext {
		return map[string]string{}, nil
	}
	orgInfo, err := catalog.getOrgInfo()
	if err != nil {
		return nil, err
	}
	return map[string]string{types.HeaderTenantContext: orgInfo.OrgId, types.HeaderAuthContext: orgInfo.OrgName}, nil
}

// getAccessControlHeader builds the data needed to set the header when tenant context is required.
// If useTenantContext is false, it returns an empty map.
// Otherwise, it finds the Org ID and name and creates the header data
func (adminCatalog *AdminCatalog) getAccessControlHeader(useTenantContext bool) (map[string]string, error) {
	if !useTenantContext {
		return map[string]string{}, nil
	}
	orgInfo, err := adminCatalog.getOrgInfo()

	if err != nil {
		return nil, err
	}
	return map[string]string{types.HeaderTenantContext: orgInfo.OrgId, types.HeaderAuthContext: orgInfo.OrgName}, nil
}

// GetControlAccess read and returns the control access parameters from a VDC
func (vdc *Vdc) GetControlAccess(ctx context.Context, useTenantContext bool) (*types.ControlAccessParams, error) {
	err := checkSanityVdcControlAccess(vdc)
	if err != nil {
		return nil, err
	}

	var tenantContextHeaders map[string]string

	if useTenantContext {
		tenantContext, err := vdc.getTenantContext()
		if err != nil {
			return nil, fmt.Errorf("error getting the tenant context - %s", err)
		}

		tenantContextHeaders = getTenantContextHeader(tenantContext)
	}

	controlAccessParams, err := vdc.client.GetAccessControl(ctx, vdc.Vdc.HREF, "vdc", vdc.Vdc.Name, tenantContextHeaders)
	if err != nil {
		return nil, fmt.Errorf("there was an error when retrieving VDC control access params - %s", err)
	}

	return controlAccessParams, nil
}

// SetControlAccess sets VDC control access parameters for everybody or individual users/groups.
// This method either sets control for everybody, passing isSharedToEveryOne true, and everyoneAccessLevel (currently only ReadOnly is supported for VDC) and nil for accessSettings,
// or can set access control for specific users/groups, passing isSharedToEveryOne false, everyoneAccessLevel "" and accessSettings filled as desired.
// The method will fail if tries to configure access control for everybody and passes individual users/groups to configure.
// It returns the control access parameters that are read from the API (using Vdc.GetControlAccess).
func (vdc *Vdc) SetControlAccess(ctx context.Context, isSharedToEveryOne bool, everyoneAccessLevel string, accessSettings []*types.AccessSetting, useTenantContext bool) (*types.ControlAccessParams, error) {
	err := checkSanityVdcControlAccess(vdc)
	if err != nil {
		return nil, err
	}

	if (isSharedToEveryOne && accessSettings != nil) && len(accessSettings) > 0 {
		return nil, fmt.Errorf("either configure access for everybody or individual users, not both at the same time")
	}

	var tenantContextHeaders map[string]string
	var accessControl = &types.ControlAccessParams{
		Xmlns: types.XMLNamespaceVCloud,
	}

	if isSharedToEveryOne { // Do configuration for everyone
		if everyoneAccessLevel == "" {
			return nil, fmt.Errorf("everyoneAccessLevel needs to be set if isSharedToEveryOne is true")
		}

		accessControl.IsSharedToEveryone = true
		accessControl.EveryoneAccessLevel = takeStringPointer(everyoneAccessLevel)

	} else { // Do configuration for individual users/groups
		if len(accessSettings) > 0 {
			accessControl.AccessSettings = &types.AccessSettingList{
				AccessSetting: accessSettings,
			}
		}
	}

	if useTenantContext {
		tenantContext, err := vdc.getTenantContext()
		if err != nil {
			return nil, fmt.Errorf("error getting the tenant context - %s", err)
		}

		tenantContextHeaders = getTenantContextHeader(tenantContext)
	}

	err = vdc.client.setAccessControlWithHttpMethod(ctx, http.MethodPut, accessControl, vdc.Vdc.HREF, "vdc", vdc.Vdc.Name, tenantContextHeaders)
	if err != nil {
		return nil, fmt.Errorf("there was an error when setting VDC control access params - %s", err)
	}

	return vdc.GetControlAccess(ctx, useTenantContext)
}

// DeleteControlAccess makes stop sharing VDC with anyone
func (vdc *Vdc) DeleteControlAccess(ctx context.Context, useTenantContext bool) (*types.ControlAccessParams, error) {
	return vdc.SetControlAccess(ctx, false, "", nil, useTenantContext)
}

// checkSanityVdcControlAccess is a function that check some Vdc attributes and returns error if any is missing. It is useful for
// checking sanity of Vdc struct before running controlAccess methods.
func checkSanityVdcControlAccess(vdc *Vdc) error {
	if vdc.client == nil {
		return fmt.Errorf("client has not been set up on Vdc struct. Please initialize it before using this method")
	}
	if vdc.Vdc == nil || vdc.Vdc.Name == "" {
		return fmt.Errorf("types.Vdc struct has not been set up on Vdc struct or Vdc.Vdc.Name is missing. Please initialize it before using this method ")
	}
	return nil
}
