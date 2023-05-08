/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"github.com/vmware/go-vcloud-director/v2/util"
)

type CatalogItem struct {
	CatalogItem *types.CatalogItem
	client      *Client
}

func NewCatalogItem(cli *Client) *CatalogItem {
	return &CatalogItem{
		CatalogItem: new(types.CatalogItem),
		client:      cli,
	}
}

func (catalogItem *CatalogItem) GetVAppTemplate(ctx context.Context) (VAppTemplate, error) {

	cat := NewVAppTemplate(catalogItem.client)

	_, err := catalogItem.client.ExecuteRequest(ctx, catalogItem.CatalogItem.Entity.HREF, http.MethodGet,
		"", "error retrieving vApp template: %s", nil, cat.VAppTemplate)

	// The request was successful
	return *cat, err

}

// Delete deletes the Catalog Item, returning an error if the vCD call fails.
// Link to API call: https://code.vmware.com/apis/220/vcloud#/doc/doc/operations/DELETE-CatalogItem.html
func (catalogItem *CatalogItem) Delete(ctx context.Context) error {
	util.Logger.Printf("[TRACE] Deleting catalog item: %#v", catalogItem.CatalogItem)
	catalogItemHREF := catalogItem.client.VCDHREF
	catalogItemHREF.Path += "/catalogItem/" + catalogItem.CatalogItem.ID[23:]

	util.Logger.Printf("[TRACE] Url for deleting catalog item: %#v and name: %s", catalogItemHREF, catalogItem.CatalogItem.Name)

	return catalogItem.client.ExecuteRequestWithoutResponse(ctx, catalogItemHREF.String(), http.MethodDelete,
		"", "error deleting Catalog item: %s", nil)
}

// queryCatalogItemList returns a list of Catalog Item for the given parent
func queryCatalogItemList(ctx context.Context, client *Client, parentField, parentValue string) ([]*types.QueryResultCatalogItemType, error) {

	catalogItemType := types.QtCatalogItem
	if client.IsSysAdmin {
		catalogItemType = types.QtAdminCatalogItem
	}

	filterText := fmt.Sprintf("%s==%s", parentField, url.QueryEscape(parentValue))

	results, err := client.cumulativeQuery(ctx, catalogItemType, nil, map[string]string{
		"type":   catalogItemType,
		"filter": filterText,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying catalog items %s", err)
	}

	if client.IsSysAdmin {
		return results.Results.AdminCatalogItemRecord, nil
	} else {
		return results.Results.CatalogItemRecord, nil
	}
}

// QueryCatalogItemList returns a list of Catalog Item for the given catalog
func (catalog *Catalog) QueryCatalogItemList(ctx context.Context) ([]*types.QueryResultCatalogItemType, error) {
	return queryCatalogItemList(ctx, catalog.client, "catalog", catalog.Catalog.ID)
}

// QueryCatalogItemList returns a list of Catalog Item for the given VDC
func (vdc *Vdc) QueryCatalogItemList(ctx context.Context) ([]*types.QueryResultCatalogItemType, error) {
	return queryCatalogItemList(ctx, vdc.client, "vdc", vdc.Vdc.ID)
}

// QueryCatalogItemList returns a list of Catalog Item for the given Admin VDC
func (vdc *AdminVdc) QueryCatalogItemList(ctx context.Context) ([]*types.QueryResultCatalogItemType, error) {
	return queryCatalogItemList(ctx, vdc.client, "vdc", vdc.AdminVdc.ID)
}

// queryVappTemplateListWithParentField returns a list of vApp templates for the given parent
func queryVappTemplateListWithParentField(ctx context.Context, client *Client, parentField, parentValue string) ([]*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateListWithFilter(ctx, client, map[string]string{
		parentField: parentValue,
	})
}

// queryVappTemplateListWithFilter returns a list of vApp templates filtered by the given filter map.
// The filter map will build a filter like filterKey==filterValue;filterKey2==filterValue2;...
func queryVappTemplateListWithFilter(ctx context.Context, client *Client, filter map[string]string) ([]*types.QueryResultVappTemplateType, error) {
	vappTemplateType := types.QtVappTemplate
	if client.IsSysAdmin {
		vappTemplateType = types.QtAdminVappTemplate
	}
	filterEncoded := ""
	for k, v := range filter {
		filterEncoded += fmt.Sprintf("%s==%s;", url.QueryEscape(k), url.QueryEscape(v))
	}
	results, err := client.cumulativeQuery(ctx, vappTemplateType, nil, map[string]string{
		"type":   vappTemplateType,
		"filter": filterEncoded[:len(filterEncoded)-1], // Removes the trailing ';'
	})
	if err != nil {
		return nil, fmt.Errorf("error querying vApp templates %s", err)
	}

	if client.IsSysAdmin {
		return results.Results.AdminVappTemplateRecord, nil
	} else {
		return results.Results.VappTemplateRecord, nil
	}
}

// QueryVappTemplateList returns a list of vApp templates for the given VDC
func (vdc *Vdc) QueryVappTemplateList(ctx context.Context) ([]*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateListWithParentField(ctx, vdc.client, "vdcName", vdc.Vdc.Name)
}

// QueryVappTemplateWithName returns one vApp template for the given VDC with the given name.
// Returns an error if it finds more than one.
func (vdc *Vdc) QueryVappTemplateWithName(ctx context.Context, vAppTemplateName string) (*types.QueryResultVappTemplateType, error) {
	vAppTemplates, err := queryVappTemplateListWithFilter(ctx, vdc.client, map[string]string{
		"vdcName": vdc.Vdc.Name,
		"name":    vAppTemplateName,
	})
	if err != nil {
		return nil, err
	}
	if len(vAppTemplates) != 1 {
		if len(vAppTemplates) == 0 {
			return nil, ErrorEntityNotFound
		}
		return nil, fmt.Errorf("found %d vApp Templates with name %s in VDC %s", len(vAppTemplates), vAppTemplateName, vdc.Vdc.Name)
	}
	return vAppTemplates[0], nil
}

// QueryVappTemplateList returns a list of vApp templates for the given VDC
func (vdc *AdminVdc) QueryVappTemplateList(ctx context.Context) ([]*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateListWithParentField(ctx, vdc.client, "vdcName", vdc.AdminVdc.Name)
}

// QueryVappTemplateWithName returns one vApp template for the given VDC with the given name.
// Returns an error if it finds more than one.
func (vdc *AdminVdc) QueryVappTemplateWithName(ctx context.Context, vAppTemplateName string) (*types.QueryResultVappTemplateType, error) {
	vAppTemplates, err := queryVappTemplateListWithFilter(ctx, vdc.client, map[string]string{
		"vdcName": vdc.AdminVdc.Name,
		"name":    vAppTemplateName,
	})
	if err != nil {
		return nil, err
	}
	if len(vAppTemplates) != 1 {
		if len(vAppTemplates) == 0 {
			return nil, ErrorEntityNotFound
		}
		return nil, fmt.Errorf("found %d vApp Templates with name %s in VDC %s", len(vAppTemplates), vAppTemplateName, vdc.AdminVdc.Name)
	}
	return vAppTemplates[0], nil
}

// QueryVappTemplateList returns a list of vApp templates for the given catalog
func (catalog *Catalog) QueryVappTemplateList(ctx context.Context) ([]*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateListWithParentField(ctx, catalog.client, "catalogName", catalog.Catalog.Name)
}

// QueryVappTemplateWithName returns one vApp template for the given Catalog with the given name.
// Returns an error if it finds more than one.
func (catalog *Catalog) QueryVappTemplateWithName(ctx context.Context, vAppTemplateName string) (*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateWithName(ctx, catalog.client, catalog.Catalog.Name, vAppTemplateName)
}

// QueryVappTemplateWithName returns one vApp template for the given Catalog with the given name.
// Returns an error if it finds more than one.
func (catalog *AdminCatalog) QueryVappTemplateWithName(ctx context.Context, vAppTemplateName string) (*types.QueryResultVappTemplateType, error) {
	return queryVappTemplateWithName(ctx, catalog.client, catalog.AdminCatalog.Name, vAppTemplateName)
}

// queryVappTemplateWithName returns one vApp template for the given Catalog with the given name.
// Returns an error if it finds more than one.
func queryVappTemplateWithName(ctx context.Context, client *Client, catalogName, vAppTemplateName string) (*types.QueryResultVappTemplateType, error) {
	vAppTemplates, err := queryVappTemplateListWithFilter(ctx, client, map[string]string{
		"catalogName": catalogName,
		"name":        vAppTemplateName,
	})
	if err != nil {
		return nil, err
	}
	if len(vAppTemplates) != 1 {
		if len(vAppTemplates) == 0 {
			return nil, ErrorEntityNotFound
		}
		return nil, fmt.Errorf("found %d vApp Templates with name %s in Catalog %s", len(vAppTemplates), vAppTemplateName, catalogName)
	}
	return vAppTemplates[0], nil
}

// queryCatalogItemFilteredList returns a list of Catalog Items with an optional filter
func queryCatalogItemFilteredList(ctx context.Context, client *Client, filter map[string]string) ([]*types.QueryResultCatalogItemType, error) {
	catalogItemType := types.QtCatalogItem
	if client.IsSysAdmin {
		catalogItemType = types.QtAdminCatalogItem
	}

	filterText := ""
	for k, v := range filter {
		if filterText != "" {
			filterText += ";"
		}
		filterText += fmt.Sprintf("%s==%s", k, url.QueryEscape(v))
	}

	notEncodedParams := map[string]string{
		"type": catalogItemType,
	}
	if filterText != "" {
		notEncodedParams["filter"] = filterText
	}
	results, err := client.cumulativeQuery(ctx, catalogItemType, nil, notEncodedParams)
	if err != nil {
		return nil, fmt.Errorf("error querying catalog items %s", err)
	}

	if client.IsSysAdmin {
		return results.Results.AdminCatalogItemRecord, nil
	} else {
		return results.Results.CatalogItemRecord, nil
	}
}

// QueryCatalogItemList returns a list of Catalog Item for the given admin catalog
func (catalog *AdminCatalog) QueryCatalogItemList(ctx context.Context) ([]*types.QueryResultCatalogItemType, error) {
	return queryCatalogItemList(ctx, catalog.client, "catalog", catalog.AdminCatalog.ID)
}

// QueryCatalogItem returns a named Catalog Item for the given catalog
func (catalog *AdminCatalog) QueryCatalogItem(ctx context.Context, name string) (*types.QueryResultCatalogItemType, error) {
	return queryCatalogItem(ctx, catalog.client, "catalog", catalog.AdminCatalog.ID, name)
}

// queryCatalogItem returns a named Catalog Item for the given parent
func queryCatalogItem(ctx context.Context, client *Client, parentField, parentValue, name string) (*types.QueryResultCatalogItemType, error) {

	result, err := queryCatalogItemFilteredList(ctx, client, map[string]string{parentField: parentValue, "name": name})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrorEntityNotFound
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("more than one item (%d) found with name %s", len(result), name)
	}
	return result[0], nil
}

// queryResultCatalogItemToCatalogItem converts a catalog item as retrieved from a query into a regular one
func queryResultCatalogItemToCatalogItem(client *Client, qr *types.QueryResultCatalogItemType) *CatalogItem {
	var catalogItem = NewCatalogItem(client)
	catalogItem.CatalogItem = &types.CatalogItem{
		HREF:        qr.HREF,
		Type:        qr.Type,
		ID:          extractUuid(qr.HREF),
		Name:        qr.Name,
		DateCreated: qr.CreationDate,
		Entity: &types.Entity{
			HREF: qr.Entity,
			Type: qr.EntityType,
			Name: qr.EntityName,
		},
	}
	return catalogItem
}

// LaunchSync starts synchronisation of a subscribed Catalog item
func (item *CatalogItem) LaunchSync(ctx context.Context) (*Task, error) {
	util.Logger.Printf("[TRACE] LaunchSync '%s' \n", item.CatalogItem.Name)
	err := WaitResource(func() (*types.TasksInProgress, error) {
		if item.CatalogItem.Tasks == nil {
			return nil, nil
		}
		err := item.Refresh(ctx)
		if err != nil {
			return nil, err
		}
		return item.CatalogItem.Tasks, nil
	})
	if err != nil {
		return nil, err
	}
	return elementLaunchSync(ctx, item.client, item.CatalogItem.HREF, "catalog item")
}

// Refresh retrieves a fresh copy of the catalog Item
func (item *CatalogItem) Refresh(ctx context.Context) error {
	_, err := item.client.ExecuteRequest(ctx, item.CatalogItem.HREF, http.MethodGet,
		"", "error retrieving catalog item: %s", nil, item.CatalogItem)
	if err != nil {
		return err
	}
	return nil
}
