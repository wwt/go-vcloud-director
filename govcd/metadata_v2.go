/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/http"
	"strings"
)

// NOTE: This "v2" is not v2 in terms of API versioning, it's just a way to separate the functions that handle
// metadata in a complete way (v2, this file) and the deprecated functions that were incomplete (v1, they lacked
// "visibility" and "domain" handling).
//
// The idea is that once a new major version of go-vcloud-director is released, one can just remove "v1" file and perform
// a minor refactoring of the code here (probably renaming functions). Also, the code in "v2" is organized differently,
// as this is classified using "CRUD blocks" (meaning that all Create functions are together, same for Read... etc),
// which makes the code more readable.

// ------------------------------------------------------------------------------------------------
// GET metadata by key
// ------------------------------------------------------------------------------------------------

// GetMetadataByKeyAndHref returns metadata from the given resource reference, corresponding to the given key and domain.
func (vcdClient *VCDClient) GetMetadataByKeyAndHref(ctx context.Context, href, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, &vcdClient.Client, href, key, isSystem)
}

// GetMetadataByKey returns VM metadata corresponding to the given key and domain.
func (vm *VM) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, vm.client, vm.VM.HREF, key, isSystem)
}

// GetMetadataByKey returns VDC metadata corresponding to the given key and domain.
func (vdc *Vdc) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, vdc.client, vdc.Vdc.HREF, key, isSystem)
}

// GetMetadataByKey returns AdminVdc metadata corresponding to the given key and domain.
func (adminVdc *AdminVdc) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, key, isSystem)
}

// GetMetadataByKey returns ProviderVdc metadata corresponding to the given key and domain.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, key, isSystem)
}

// GetMetadataByKey returns VApp metadata corresponding to the given key and domain.
func (vapp *VApp) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, vapp.client, vapp.VApp.HREF, key, isSystem)
}

// GetMetadataByKey returns VAppTemplate metadata corresponding to the given key and domain.
func (vAppTemplate *VAppTemplate) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, key, isSystem)
}

// GetMetadataByKey returns MediaRecord metadata corresponding to the given key and domain.
func (mediaRecord *MediaRecord) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, key, isSystem)
}

// GetMetadataByKey returns Media metadata corresponding to the given key and domain.
func (media *Media) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, media.client, media.Media.HREF, key, isSystem)
}

// GetMetadataByKey returns Catalog metadata corresponding to the given key and domain.
func (catalog *Catalog) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, catalog.client, catalog.Catalog.HREF, key, isSystem)
}

// GetMetadataByKey returns AdminCatalog metadata corresponding to the given key and domain.
func (adminCatalog *AdminCatalog) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, key, isSystem)
}

// GetMetadataByKey returns the Org metadata corresponding to the given key and domain.
func (org *Org) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, org.client, org.Org.HREF, key, isSystem)
}

// GetMetadataByKey returns the AdminOrg metadata corresponding to the given key and domain.
// Note: Requires system administrator privileges.
func (adminOrg *AdminOrg) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, key, isSystem)
}

// GetMetadataByKey returns the metadata corresponding to the given key and domain.
func (disk *Disk) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, disk.client, disk.Disk.HREF, key, isSystem)
}

// GetMetadataByKey returns OrgVDCNetwork metadata corresponding to the given key and domain.
func (orgVdcNetwork *OrgVDCNetwork) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, orgVdcNetwork.client, orgVdcNetwork.OrgVDCNetwork.HREF, key, isSystem)
}

// GetMetadataByKey returns CatalogItem metadata corresponding to the given key and domain.
func (catalogItem *CatalogItem) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	return getMetadataByKey(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, key, isSystem)
}

// GetMetadataByKey returns OpenApiOrgVdcNetwork metadata corresponding to the given key and domain.
// NOTE: This function cannot retrieve metadata if the network belongs to a VDC Group.
// TODO: This function is currently using XML API underneath as OpenAPI metadata is still not supported.
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) GetMetadataByKey(ctx context.Context, key string, isSystem bool) (*types.MetadataValue, error) {
	href := fmt.Sprintf("%s/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), extractUuid(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID))
	return getMetadataByKey(ctx, openApiOrgVdcNetwork.client, href, key, isSystem)
}

// ------------------------------------------------------------------------------------------------
// GET all metadata
// ------------------------------------------------------------------------------------------------

// GetMetadataByHref returns metadata from the given resource reference.
func (vcdClient *VCDClient) GetMetadataByHref(ctx context.Context, href string) (*types.Metadata, error) {
	return getMetadata(ctx, &vcdClient.Client, href)
}

// GetMetadata returns VM metadata.
func (vm *VM) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, vm.client, vm.VM.HREF)
}

// GetMetadata returns VDC metadata.
func (vdc *Vdc) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, vdc.client, vdc.Vdc.HREF)
}

// GetMetadata returns AdminVdc metadata.
func (adminVdc *AdminVdc) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, adminVdc.client, adminVdc.AdminVdc.HREF)
}

// GetMetadata returns ProviderVdc metadata.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF)
}

// GetMetadata returns VApp metadata.
func (vapp *VApp) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, vapp.client, vapp.VApp.HREF)
}

// GetMetadata returns VAppTemplate metadata.
func (vAppTemplate *VAppTemplate) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF)
}

// GetMetadata returns MediaRecord metadata.
func (mediaRecord *MediaRecord) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF)
}

// GetMetadata returns Media metadata.
func (media *Media) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, media.client, media.Media.HREF)
}

// GetMetadata returns Catalog metadata.
func (catalog *Catalog) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, catalog.client, catalog.Catalog.HREF)
}

// GetMetadata returns AdminCatalog metadata.
func (adminCatalog *AdminCatalog) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF)
}

// GetMetadata returns the Org metadata of the corresponding organization seen as administrator
func (org *Org) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, org.client, org.Org.HREF)
}

// GetMetadata returns the AdminOrg metadata of the corresponding organization seen as administrator
func (adminOrg *AdminOrg) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, adminOrg.client, adminOrg.AdminOrg.HREF)
}

// GetMetadata returns the metadata of the corresponding independent disk
func (disk *Disk) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, disk.client, disk.Disk.HREF)
}

// GetMetadata returns OrgVDCNetwork metadata.
func (orgVdcNetwork *OrgVDCNetwork) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, orgVdcNetwork.client, orgVdcNetwork.OrgVDCNetwork.HREF)
}

// GetMetadata returns CatalogItem metadata.
func (catalogItem *CatalogItem) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, catalogItem.client, catalogItem.CatalogItem.HREF)
}

// GetMetadata returns OpenApiOrgVdcNetwork metadata.
// NOTE: This function cannot retrieve metadata if the network belongs to a VDC Group.
// TODO: This function is currently using XML API underneath as OpenAPI metadata is still not supported.
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	href := fmt.Sprintf("%s/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), extractUuid(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID))
	return getMetadata(ctx, openApiOrgVdcNetwork.client, href)
}

// ------------------------------------------------------------------------------------------------
// ADD metadata async
// ------------------------------------------------------------------------------------------------

// AddMetadataEntryWithVisibilityByHrefAsync adds metadata to the given resource reference with the given key, value, type and visibility
// and returns the task.
func (vcdClient *VCDClient) AddMetadataEntryWithVisibilityByHrefAsync(ctx context.Context, href, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, &vcdClient.Client, href, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given VM with the given key, value, type and visibility
// // and returns the task.
func (vm *VM) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, vm.client, vm.VM.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given AdminVdc with the given key, value, type and visibility
// and returns the task.
func (adminVdc *AdminVdc) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given ProviderVdc with the given key, value, type and visibility
// and returns the task.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given VApp with the given key, value, type and visibility
// and returns the task.
func (vapp *VApp) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, vapp.client, vapp.VApp.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given VAppTemplate with the given key, value, type and visibility
// and returns the task.
func (vAppTemplate *VAppTemplate) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given MediaRecord with the given key, value, type and visibility
// and returns the task.
func (mediaRecord *MediaRecord) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given Media with the given key, value, type and visibility
// and returns the task.
func (media *Media) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, media.client, media.Media.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given AdminCatalog with the given key, value, type and visibility
// and returns the task.
func (adminCatalog *AdminCatalog) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given AdminOrg with the given key, value, type and visibility
// and returns the task.
func (adminOrg *AdminOrg) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given Disk with the given key, value, type and visibility
// and returns the task.
func (disk *Disk) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, disk.client, disk.Disk.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given OrgVDCNetwork with the given key, value, type and visibility
// and returns the task.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibilityAsync adds metadata to the given Catalog Item with the given key, value, type and visibility
// and returns the task.
func (catalogItem *CatalogItem) AddMetadataEntryWithVisibilityAsync(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	return addMetadata(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, key, value, typedValue, visibility, isSystem)
}

// ------------------------------------------------------------------------------------------------
// ADD metadata
// ------------------------------------------------------------------------------------------------

// AddMetadataEntryWithVisibilityByHref adds metadata to the given resource reference with the given key, value, type and visibility
// and waits for completion.
func (vcdClient *VCDClient) AddMetadataEntryWithVisibilityByHref(ctx context.Context, href, key, value, typedValue, visibility string, isSystem bool) error {
	task, err := vcdClient.AddMetadataEntryWithVisibilityByHrefAsync(ctx, href, key, value, typedValue, visibility, isSystem)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver VM and waits for the task to finish.
func (vm *VM) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, vm.client, vm.VM.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver AdminVdc and waits for the task to finish.
func (adminVdc *AdminVdc) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver ProviderVdc and waits for the task to finish.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver VApp and waits for the task to finish.
func (vapp *VApp) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, vapp.client, vapp.VApp.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver VAppTemplate and waits for the task to finish.
func (vAppTemplate *VAppTemplate) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver MediaRecord and waits for the task to finish.
func (mediaRecord *MediaRecord) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver Media and waits for the task to finish.
func (media *Media) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, media.client, media.Media.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver AdminCatalog and waits for the task to finish.
func (adminCatalog *AdminCatalog) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver AdminOrg and waits for the task to finish.
func (adminOrg *AdminOrg) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver Disk and waits for the task to finish.
func (disk *Disk) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, disk.client, disk.Disk.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver OrgVDCNetwork and waits for the task to finish.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver CatalogItem and waits for the task to finish.
func (catalogItem *CatalogItem) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	return addMetadataAndWait(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, key, value, typedValue, visibility, isSystem)
}

// AddMetadataEntryWithVisibility adds metadata to the receiver OpenApiOrgVdcNetwork and waits for the task to finish.
// Note: It doesn't add metadata to networks that belong to a VDC Group.
// TODO: This function is currently using XML API underneath as OpenAPI metadata is still not supported.
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) AddMetadataEntryWithVisibility(ctx context.Context, key, value, typedValue, visibility string, isSystem bool) error {
	href := fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), extractUuid(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID))
	task, err := addMetadata(ctx, openApiOrgVdcNetwork.client, href, key, value, typedValue, visibility, isSystem)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// ------------------------------------------------------------------------------------------------
// MERGE metadata async
// ------------------------------------------------------------------------------------------------

// MergeMetadataWithVisibilityByHrefAsync updates the metadata entries present in the referenced entity and creates the ones not present, then
// returns the task.
func (vcdClient *VCDClient) MergeMetadataWithVisibilityByHrefAsync(ctx context.Context, href string, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, &vcdClient.Client, href, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges VM metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then returns the task.
func (vm *VM) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, vm.client, vm.VM.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges AdminVdc metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (adminVdc *AdminVdc) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges Provider VDC metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges VApp metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (vapp *VApp) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, vapp.client, vapp.VApp.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges VAppTemplate metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (vAppTemplate *VAppTemplate) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges MediaRecord metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (mediaRecord *MediaRecord) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges Media metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (media *Media) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, media.client, media.Media.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges AdminCatalog metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (adminCatalog *AdminCatalog) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges AdminOrg metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (adminOrg *AdminOrg) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges Disk metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (disk *Disk) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, disk.client, disk.Disk.HREF, metadata)
}

// MergeMetadataWithMetadataValuesAsync merges OrgVDCNetwork metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), metadata)
}

// MergeMetadataWithMetadataValuesAsync merges CatalogItem metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (catalogItem *CatalogItem) MergeMetadataWithMetadataValuesAsync(ctx context.Context, metadata map[string]types.MetadataValue) (Task, error) {
	return mergeAllMetadata(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, metadata)
}

// ------------------------------------------------------------------------------------------------
// MERGE metadata
// ------------------------------------------------------------------------------------------------

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver VM and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (vm *VM) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, vm.client, vm.VM.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver AdminVdc and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (adminVdc *AdminVdc) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver ProviderVdc and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver VApp and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (vApp *VApp) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, vApp.client, vApp.VApp.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver VAppTemplate and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (vAppTemplate *VAppTemplate) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver MediaRecord and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (mediaRecord *MediaRecord) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver Media and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (media *Media) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, media.client, media.Media.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver AdminCatalog and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (adminCatalog *AdminCatalog) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver AdminOrg and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (adminOrg *AdminOrg) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver Disk and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (disk *Disk) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, disk.client, disk.Disk.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver OrgVDCNetwork and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver CatalogItem and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func (catalogItem *CatalogItem) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	return mergeMetadataAndWait(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, metadata)
}

// MergeMetadataWithMetadataValues updates the metadata values that are already present in the receiver OpenApiOrgVdcNetwork and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
// Note: It doesn't merge metadata to networks that belong to a VDC Group.
// TODO: This function is currently using XML API underneath as OpenAPI metadata is still not supported.
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) MergeMetadataWithMetadataValues(ctx context.Context, metadata map[string]types.MetadataValue) error {
	href := fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), extractUuid(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID))
	task, err := mergeAllMetadata(ctx, openApiOrgVdcNetwork.client, href, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// ------------------------------------------------------------------------------------------------
// DELETE metadata async
// ------------------------------------------------------------------------------------------------

// DeleteMetadataEntryWithDomainByHrefAsync deletes metadata from the given resource reference, depending on key provided as input
// and returns a task.
func (vcdClient *VCDClient) DeleteMetadataEntryWithDomainByHrefAsync(ctx context.Context, href, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, &vcdClient.Client, href, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes VM metadata associated to the input key and returns the task.
func (vm *VM) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, vm.client, vm.VM.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes AdminVdc metadata associated to the input key and returns the task.
func (adminVdc *AdminVdc) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes ProviderVdc metadata associated to the input key and returns the task.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes VApp metadata associated to the input key and returns the task.
func (vapp *VApp) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, vapp.client, vapp.VApp.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes VAppTemplate metadata associated to the input key and returns the task.
func (vAppTemplate *VAppTemplate) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes MediaRecord metadata associated to the input key and returns the task.
func (mediaRecord *MediaRecord) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes Media metadata associated to the input key and returns the task.
func (media *Media) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, media.client, media.Media.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes AdminCatalog metadata associated to the input key and returns the task.
func (adminCatalog *AdminCatalog) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes AdminOrg metadata associated to the input key and returns the task.
func (adminOrg *AdminOrg) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes Disk metadata associated to the input key and returns the task.
func (disk *Disk) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, disk.client, disk.Disk.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes OrgVDCNetwork metadata associated to the input key and returns the task.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), key, isSystem)
}

// DeleteMetadataEntryWithDomainAsync deletes CatalogItem metadata associated to the input key and returns the task.
func (catalogItem *CatalogItem) DeleteMetadataEntryWithDomainAsync(ctx context.Context, key string, isSystem bool) (Task, error) {
	return deleteMetadata(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, key, isSystem)
}

// ------------------------------------------------------------------------------------------------
// DELETE metadata
// ------------------------------------------------------------------------------------------------

// DeleteMetadataEntryWithDomainByHref deletes metadata from the given resource reference, depending on key provided as input
// and waits for the task to finish.
func (vcdClient *VCDClient) DeleteMetadataEntryWithDomainByHref(ctx context.Context, href, key string, isSystem bool) error {
	task, err := vcdClient.DeleteMetadataEntryWithDomainByHrefAsync(ctx, href, key, isSystem)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntryWithDomain deletes VM metadata associated to the input key and waits for the task to finish.
func (vm *VM) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, vm.client, vm.VM.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes AdminVdc metadata associated to the input key and waits for the task to finish.
// Note: Requires system administrator privileges.
func (adminVdc *AdminVdc) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, adminVdc.client, getAdminURL(adminVdc.AdminVdc.HREF), key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes ProviderVdc metadata associated to the input key and waits for the task to finish.
// Note: Requires system administrator privileges.
func (providerVdc *ProviderVdc) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, providerVdc.client, providerVdc.ProviderVdc.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes VApp metadata associated to the input key and waits for the task to finish.
func (vApp *VApp) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, vApp.client, vApp.VApp.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes VAppTemplate metadata associated to the input key and waits for the task to finish.
func (vAppTemplate *VAppTemplate) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes MediaRecord metadata associated to the input key and waits for the task to finish.
func (mediaRecord *MediaRecord) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes Media metadata associated to the input key and waits for the task to finish.
func (media *Media) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, media.client, media.Media.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes AdminCatalog metadata associated to the input key and waits for the task to finish.
func (adminCatalog *AdminCatalog) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes AdminOrg metadata associated to the input key and waits for the task to finish.
func (adminOrg *AdminOrg) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes Disk metadata associated to the input key and waits for the task to finish.
func (disk *Disk) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, disk.client, disk.Disk.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes OrgVDCNetwork metadata associated to the input key and waits for the task to finish.
// Note: Requires system administrator privileges.
func (orgVdcNetwork *OrgVDCNetwork) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes CatalogItem metadata associated to the input key and waits for the task to finish.
func (catalogItem *CatalogItem) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	return deleteMetadataAndWait(ctx, catalogItem.client, catalogItem.CatalogItem.HREF, key, isSystem)
}

// DeleteMetadataEntryWithDomain deletes OpenApiOrgVdcNetwork metadata associated to the input key and waits for the task to finish.
// Note: It doesn't delete metadata from networks that belong to a VDC Group.
// TODO: This function is currently using XML API underneath as OpenAPI metadata is still not supported.
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) DeleteMetadataEntryWithDomain(ctx context.Context, key string, isSystem bool) error {
	href := fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), extractUuid(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID))
	task, err := deleteMetadata(ctx, openApiOrgVdcNetwork.client, href, key, isSystem)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// ------------------------------------------------------------------------------------------------
// Generic private functions
// ------------------------------------------------------------------------------------------------

// getMetadata is a generic function to retrieve metadata from VCD
func getMetadataByKey(ctx context.Context, client *Client, requestUri, key string, isSystem bool) (*types.MetadataValue, error) {
	metadata := &types.MetadataValue{}
	href := requestUri + "/metadata/"

	if isSystem {
		href += "SYSTEM/"
	}

	_, err := client.ExecuteRequest(ctx, href+key, http.MethodGet, types.MimeMetaData, "error retrieving metadata by key "+key+": %s", nil, metadata)
	return metadata, err
}

// getMetadata is a generic function to retrieve metadata from VCD
func getMetadata(ctx context.Context, client *Client, requestUri string) (*types.Metadata, error) {
	metadata := &types.Metadata{}

	_, err := client.ExecuteRequest(ctx, requestUri+"/metadata/", http.MethodGet, types.MimeMetaData, "error retrieving metadata: %s", nil, metadata)
	return metadata, err
}

// addMetadata adds metadata to an entity.
// If the metadata entry is of the SYSTEM domain (isSystem=true), one can set different types of Visibility:
// types.MetadataReadOnlyVisibility, types.MetadataHiddenVisibility but NOT types.MetadataReadWriteVisibility.
// If the metadata entry is of the GENERAL domain (isSystem=false), visibility is always types.MetadataReadWriteVisibility.
// In terms of typedValues, that must be one of:
// types.MetadataStringValue, types.MetadataNumberValue, types.MetadataDateTimeValue and types.MetadataBooleanValue.
func addMetadata(ctx context.Context, client *Client, requestUri, key, value, typedValue, visibility string, isSystem bool) (Task, error) {
	apiEndpoint := urlParseRequestURI(requestUri)
	newMetadata := &types.MetadataValue{
		Xmlns: types.XMLNamespaceVCloud,
		Xsi:   types.XMLNamespaceXSI,
		TypedValue: &types.MetadataTypedValue{
			XsiType: typedValue,
			Value:   value,
		},
		Domain: &types.MetadataDomainTag{
			Visibility: visibility,
			Domain:     "SYSTEM",
		},
	}

	if isSystem {
		apiEndpoint.Path += "/metadata/SYSTEM/" + key
	} else {
		apiEndpoint.Path += "/metadata/" + key
		newMetadata.Domain.Domain = "GENERAL"
		if visibility != types.MetadataReadWriteVisibility {
			newMetadata.Domain.Visibility = types.MetadataReadWriteVisibility
		}
	}

	domain := newMetadata.Domain.Visibility
	task, err := client.ExecuteTaskRequest(ctx, apiEndpoint.String(), http.MethodPut, types.MimeMetaDataValue, "error adding metadata: %s", newMetadata)

	// Workaround for ugly error returned by VCD: "API Error: 500: [ <uuid> ] visibility"
	if err != nil && strings.HasSuffix(err.Error(), "visibility") {
		err = fmt.Errorf("error adding metadata with key %s: visibility cannot be %s when domain is %s: %s", key, visibility, domain, err)
	}
	return task, err
}

// addMetadataAndWait adds metadata to an entity and waits for the task completion.
// The function supports passing a value that requires a typed value that must be one of:
// types.MetadataStringValue, types.MetadataNumberValue, types.MetadataDateTimeValue and types.MetadataBooleanValue.
// Visibility also needs to be one of: types.MetadataReadOnlyVisibility, types.MetadataHiddenVisibility or types.MetadataReadWriteVisibility
func addMetadataAndWait(ctx context.Context, client *Client, requestUri, key, value, typedValue, visibility string, isSystem bool) error {
	task, err := addMetadata(ctx, client, requestUri, key, value, typedValue, visibility, isSystem)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// mergeAllMetadata updates the metadata values that are already present in VCD and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// If the operation is successful, it returns the created task.
func mergeAllMetadata(ctx context.Context, client *Client, requestUri string, metadata map[string]types.MetadataValue) (Task, error) {
	var metadataToMerge []*types.MetadataEntry
	for key, value := range metadata {
		metadataToMerge = append(metadataToMerge, &types.MetadataEntry{
			Xmlns:      types.XMLNamespaceVCloud,
			Xsi:        types.XMLNamespaceXSI,
			Key:        key,
			TypedValue: value.TypedValue,
			Domain:     value.Domain,
		})
	}

	newMetadata := &types.Metadata{
		Xmlns:         types.XMLNamespaceVCloud,
		Xsi:           types.XMLNamespaceXSI,
		MetadataEntry: metadataToMerge,
	}

	apiEndpoint := urlParseRequestURI(requestUri)
	apiEndpoint.Path += "/metadata"

	return client.ExecuteTaskRequest(ctx, apiEndpoint.String(), http.MethodPost, types.MimeMetaData, "error adding metadata: %s", newMetadata)
}

// mergeAllMetadata updates the metadata values that are already present in VCD and creates the ones not present.
// The input metadata map has a "metadata key"->"metadata value" relation.
// This function waits until merge finishes.
func mergeMetadataAndWait(ctx context.Context, client *Client, requestUri string, metadata map[string]types.MetadataValue) error {
	task, err := mergeAllMetadata(ctx, client, requestUri, metadata)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// deleteMetadata deletes metadata associated to the input key from an entity referenced by its URI, then returns the
// task.
func deleteMetadata(ctx context.Context, client *Client, requestUri string, key string, isSystem bool) (Task, error) {
	apiEndpoint := urlParseRequestURI(requestUri)
	if isSystem {
		apiEndpoint.Path += "/metadata/SYSTEM/" + key
	} else {
		apiEndpoint.Path += "/metadata/" + key
	}

	return client.ExecuteTaskRequest(ctx, apiEndpoint.String(), http.MethodDelete, "", "error deleting metadata: %s", nil)
}

// deleteMetadata deletes metadata associated to the input key from an entity referenced by its URI.
func deleteMetadataAndWait(ctx context.Context, client *Client, requestUri string, key string, isSystem bool) error {
	task, err := deleteMetadata(ctx, client, requestUri, key, isSystem)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}
