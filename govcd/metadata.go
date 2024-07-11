/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// All functions here should not be used as they are deprecated in favor of those present in "metadata_v2".
// Remove this file once go-vcloud-director v3.0 is released.

// AddMetadataEntryByHref adds metadata typedValue and key/value pair provided as input to the given resource reference,
// then waits for the task to finish.
// Deprecated: Use VCDClient.AddMetadataEntryWithVisibilityByHref instead
func (vcdClient *VCDClient) AddMetadataEntryByHref(ctx context.Context, href, typedValue, key, value string) error {
	task, err := vcdClient.AddMetadataEntryByHrefAsync(ctx, href, typedValue, key, value)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryByHrefAsync adds metadata typedValue and key/value pair provided as input to the given resource reference
// and returns the task.
// Deprecated: Use VCDClient.AddMetadataEntryWithVisibilityByHrefAsync instead.
func (vcdClient *VCDClient) AddMetadataEntryByHrefAsync(ctx context.Context, href, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, &vcdClient.Client, typedValue, key, value, href)
}

// MergeMetadataByHrefAsync merges metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// and returns the task.
// Deprecated: Use VCDClient.MergeMetadataWithVisibilityByHrefAsync instead.
func (vcdClient *VCDClient) MergeMetadataByHrefAsync(ctx context.Context, href, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, &vcdClient.Client, typedValue, metadata, href)
}

// MergeMetadataByHref merges metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
func (vcdClient *VCDClient) MergeMetadataByHref(ctx context.Context, href, typedValue string, metadata map[string]interface{}) error {
	task, err := vcdClient.MergeMetadataByHrefAsync(ctx, href, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntryByHref deletes metadata from the given resource reference, depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use VCDClient.DeleteMetadataEntryWithDomainByHref
func (vcdClient *VCDClient) DeleteMetadataEntryByHref(ctx context.Context, href, key string) error {
	task, err := vcdClient.DeleteMetadataEntryByHrefAsync(ctx, href, key)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntryByHrefAsync deletes metadata from the given resource reference, depending on key provided as input
// and returns a task.
// Deprecated: Use VCDClient.DeleteMetadataEntryWithDomainByHrefAsync
func (vcdClient *VCDClient) DeleteMetadataEntryByHrefAsync(ctx context.Context, href, key string) (Task, error) {
	return deleteMetadata(ctx, &vcdClient.Client, href, "", key, false)
}

// AddMetadataEntry adds VM metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use VM.AddMetadataEntryWithVisibility instead
func (vm *VM) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := vm.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vm.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds VM metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use VM.AddMetadataEntryWithVisibilityAsync instead
func (vm *VM) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vm.client, typedValue, key, value, vm.VM.HREF)
}

// MergeMetadataAsync merges VM metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then returns the task.
// Deprecated: Use VM.MergeMetadataWithMetadataValuesAsync instead
func (vm *VM) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, vm.client, typedValue, metadata, vm.VM.HREF)
}

// MergeMetadata merges VM metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use VM.MergeMetadataWithMetadataValues
func (vm *VM) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := vm.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes VM metadata by key provided as input and waits for the task to finish.
// Deprecated: Use VM.DeleteMetadataEntryWithDomain instead
func (vm *VM) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := vm.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vm.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes VM metadata depending on key provided as input
// and returns the task.
// Deprecated: Use VM.DeleteMetadataEntryWithDomainAsync instead
func (vm *VM) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vm.client, vm.VM.HREF, vm.VM.Name, key, false)
}

// AddMetadataEntry adds VDC metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.AddMetadataEntryWithVisibility instead
func (vdc *Vdc) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := vdc.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vdc.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntry adds VDC metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use AdminVdc.AddMetadataEntryWithVisibility instead
func (adminVdc *AdminVdc) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := adminVdc.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryAsync adds VDC metadata typedValue and key/value pair provided as input and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.AddMetadataEntryWithVisibilityAsync instead
func (vdc *Vdc) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vdc.client, typedValue, key, value, getAdminURL(vdc.Vdc.HREF))
}

// AddMetadataEntryAsync adds AdminVdc metadata typedValue and key/value pair provided as input and returns the task.
// Deprecated: Use AdminVdc.AddMetadataEntryWithVisibilityAsync instead
func (adminVdc *AdminVdc) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, adminVdc.client, typedValue, key, value, adminVdc.AdminVdc.HREF)
}

// MergeMetadataAsync merges VDC metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.MergeMetadataWithMetadataValuesAsync
func (vdc *Vdc) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, vdc.client, typedValue, metadata, getAdminURL(vdc.Vdc.HREF))
}

// MergeMetadataAsync merges AdminVdc metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminVdc.MergeMetadataWithMetadataValuesAsync
func (adminVdc *AdminVdc) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, adminVdc.client, typedValue, metadata, adminVdc.AdminVdc.HREF)
}

// MergeMetadata merges VDC metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.MergeMetadataWithMetadataValues
func (vdc *Vdc) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := vdc.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// MergeMetadata merges AdminVdc metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminVdc.MergeMetadataWithMetadataValues
func (adminVdc *AdminVdc) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := adminVdc.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes VDC metadata by key provided as input and waits for
// the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.DeleteMetadataEntryWithDomain
func (vdc *Vdc) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := vdc.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vdc.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntry deletes AdminVdc metadata by key provided as input and waits for
// the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.DeleteMetadataEntryWithDomain
func (adminVdc *AdminVdc) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := adminVdc.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = adminVdc.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes VDC metadata depending on key provided as input and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.DeleteMetadataEntryWithDomainAsync
func (vdc *Vdc) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(vdc.client, getAdminURL(vdc.Vdc.HREF), vdc.Vdc.Name, key, false)
}

// DeleteMetadataEntryAsync deletes VDC metadata depending on key provided as input and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use AdminVdc.DeleteMetadataEntryWithDomainAsync
func (adminVdc *AdminVdc) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, adminVdc.client, adminVdc.AdminVdc.HREF, adminVdc.AdminVdc.Name, key, false)
}

// AddMetadataEntry adds Provider VDC metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.AddMetadataEntryWithVisibility instead
func (providerVdc *ProviderVdc) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := providerVdc.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = providerVdc.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds Provider VDC metadata typedValue and key/value pair provided as input and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.AddMetadataEntryWithVisibilityAsync instead
func (providerVdc *ProviderVdc) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, providerVdc.client, typedValue, key, value, providerVdc.ProviderVdc.HREF)
}

// MergeMetadataAsync merges Provider VDC metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.MergeMetadataWithMetadataValuesAsync
func (providerVdc *ProviderVdc) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, providerVdc.client, typedValue, metadata, providerVdc.ProviderVdc.HREF)
}

// MergeMetadata merges Provider VDC metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.MergeMetadataWithMetadataValues
func (providerVdc *ProviderVdc) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := providerVdc.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes Provider VDC metadata by key provided as input and waits for
// the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.DeleteMetadataEntryWithDomain
func (providerVdc *ProviderVdc) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := providerVdc.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = providerVdc.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes Provider VDC metadata depending on key provided as input and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use ProviderVdc.DeleteMetadataEntryWithDomainAsync
func (providerVdc *ProviderVdc) DeleteMetadataEntryAsync(key string) (Task, error) {
	return deleteMetadata(providerVdc.client, providerVdc.ProviderVdc.HREF, providerVdc.ProviderVdc.Name, key, false)
}

// AddMetadataEntry adds VApp metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use VApp.AddMetadataEntryWithVisibility instead
func (vapp *VApp) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := vapp.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vapp.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds VApp metadata typedValue and key/value pair provided as input and returns the task.
// Deprecated: Use VApp.AddMetadataEntryWithVisibilityAsync instead
func (vapp *VApp) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vapp.client, typedValue, key, value, vapp.VApp.HREF)
}

// MergeMetadataAsync merges VApp metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use VApp.MergeMetadataWithMetadataValuesAsync
func (vapp *VApp) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, vapp.client, typedValue, metadata, vapp.VApp.HREF)
}

// MergeMetadata merges VApp metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use VApp.MergeMetadataWithMetadataValues
func (vapp *VApp) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := vapp.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes VApp metadata by key provided as input and waits for
// the task to finish.
// Deprecated: Use VApp.DeleteMetadataEntryWithDomain instead
func (vapp *VApp) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := vapp.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vapp.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes VApp metadata depending on key provided as input and returns the task.
// Deprecated: Use VApp.DeleteMetadataEntryWithDomainAsync instead
func (vapp *VApp) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vapp.client, vapp.VApp.HREF, vapp.VApp.Name, key, false)
}

// AddMetadataEntry adds VAppTemplate metadata typedValue and key/value pair provided as input and
// waits for the task to finish.
// Deprecated: Use VAppTemplate.AddMetadataEntryWithVisibility instead
func (vAppTemplate *VAppTemplate) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := vAppTemplate.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vAppTemplate.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds VAppTemplate metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use VAppTemplate.AddMetadataEntryWithVisibilityAsync instead
func (vAppTemplate *VAppTemplate) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vAppTemplate.client, typedValue, key, value, vAppTemplate.VAppTemplate.HREF)
}

// MergeMetadataAsync merges VAppTemplate metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use VAppTemplate.MergeMetadataWithMetadataValuesAsync
func (vAppTemplate *VAppTemplate) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, vAppTemplate.client, typedValue, metadata, vAppTemplate.VAppTemplate.HREF)
}

// MergeMetadata merges VAppTemplate metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use VAppTemplate.MergeMetadataWithMetadataValues
func (vAppTemplate *VAppTemplate) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := vAppTemplate.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes VAppTemplate metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use VAppTemplate.DeleteMetadataEntryWithDomain instead
func (vAppTemplate *VAppTemplate) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := vAppTemplate.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = vAppTemplate.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes VAppTemplate metadata depending on key provided as input
// and returns the task.
// Deprecated: Use VAppTemplate.DeleteMetadataEntryWithDomainAsync instead
func (vAppTemplate *VAppTemplate) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, vAppTemplate.VAppTemplate.Name, key, false)
}

// AddMetadataEntry adds MediaRecord metadata typedValue and key/value pair provided as input and
// waits for the task to finish.
// Deprecated: Use MediaRecord.AddMetadataEntryWithVisibility instead
func (mediaRecord *MediaRecord) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := mediaRecord.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = mediaRecord.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds MediaRecord metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use MediaRecord.AddMetadataEntryWithVisibilityAsync instead
func (mediaRecord *MediaRecord) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, mediaRecord.client, typedValue, key, value, mediaRecord.MediaRecord.HREF)
}

// MergeMetadataAsync merges MediaRecord metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use MediaRecord.MergeMetadataWithMetadataValuesAsync
func (mediaRecord *MediaRecord) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, mediaRecord.client, typedValue, metadata, mediaRecord.MediaRecord.HREF)
}

// MergeMetadata merges MediaRecord metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use MediaRecord.MergeMetadataWithMetadataValues
func (mediaRecord *MediaRecord) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := mediaRecord.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes MediaRecord metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use MediaRecord.DeleteMetadataEntryWithDomain instead
func (mediaRecord *MediaRecord) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := mediaRecord.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = mediaRecord.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes MediaRecord metadata depending on key provided as input
// and returns the task.
// Deprecated: Use MediaRecord.DeleteMetadataEntryWithDomainAsync instead
func (mediaRecord *MediaRecord) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, mediaRecord.MediaRecord.Name, key, false)
}

// AddMetadataEntry adds Media metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use Media.AddMetadataEntryWithVisibility instead
func (media *Media) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := media.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = media.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds Media metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use Media.AddMetadataEntryWithVisibilityAsync instead
func (media *Media) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, media.client, typedValue, key, value, media.Media.HREF)
}

// MergeMetadataAsync merges Media metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use Media.MergeMetadataWithMetadataValuesAsync
func (media *Media) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, media.client, typedValue, metadata, media.Media.HREF)
}

// MergeMetadata merges Media metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use Media.MergeMetadataWithMetadataValues
func (media *Media) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := media.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes Media metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use Media.DeleteMetadataEntryWithDomain instead
func (media *Media) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := media.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = media.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes Media metadata depending on key provided as input
// and returns the task.
// Deprecated: Use Media.DeleteMetadataEntryWithDomainAsync instead
func (media *Media) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, media.client, media.Media.HREF, media.Media.Name, key, false)
}

// AddMetadataEntry adds AdminCatalog metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use AdminCatalog.AddMetadataEntryWithVisibility instead
func (adminCatalog *AdminCatalog) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := adminCatalog.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = adminCatalog.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMetadataEntryAsync adds AdminCatalog metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use AdminCatalog.AddMetadataEntryWithVisibilityAsync instead
func (adminCatalog *AdminCatalog) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, adminCatalog.client, typedValue, key, value, adminCatalog.AdminCatalog.HREF)
}

// MergeMetadataAsync merges AdminCatalog metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminCatalog.MergeMetadataWithMetadataValuesAsync
func (adminCatalog *AdminCatalog) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, adminCatalog.client, typedValue, metadata, adminCatalog.AdminCatalog.HREF)
}

// MergeMetadata merges AdminCatalog metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminCatalog.MergeMetadataWithMetadataValues
func (adminCatalog *AdminCatalog) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := adminCatalog.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes AdminCatalog metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use AdminCatalog.DeleteMetadataEntryWithDomain instead
func (adminCatalog *AdminCatalog) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := adminCatalog.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return err
	}

	err = adminCatalog.Refresh(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMetadataEntryAsync deletes AdminCatalog metadata depending on key provided as input
// and returns a task.
// Deprecated: Use AdminCatalog.DeleteMetadataEntryWithDomainAsync instead
func (adminCatalog *AdminCatalog) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, adminCatalog.client, adminCatalog.AdminCatalog.HREF, adminCatalog.AdminCatalog.Name, key, false)
}

// AddMetadataEntry adds AdminOrg metadata key/value pair provided as input to the corresponding organization seen as administrator
// and waits for completion.
// Deprecated: Use AdminOrg.AddMetadataEntryWithVisibility instead
func (adminOrg *AdminOrg) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := adminOrg.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryAsync adds AdminOrg metadata key/value pair provided as input to the corresponding organization seen as administrator
// and returns a task.
// Deprecated: Use AdminOrg.AddMetadataEntryWithVisibilityAsync instead
func (adminOrg *AdminOrg) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, adminOrg.client, typedValue, key, value, adminOrg.AdminOrg.HREF)
}

// MergeMetadataAsync merges AdminOrg metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminOrg.MergeMetadataWithMetadataValuesAsync
func (adminOrg *AdminOrg) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, adminOrg.client, typedValue, metadata, adminOrg.AdminOrg.HREF)
}

// MergeMetadata merges AdminOrg metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use AdminOrg.MergeMetadataWithMetadataValues
func (adminOrg *AdminOrg) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := adminOrg.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes metadata of the corresponding organization with the given key, and waits for completion
// Deprecated: Use AdminOrg.DeleteMetadataEntryWithDomain instead
func (adminOrg *AdminOrg) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := adminOrg.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for organization task: %s", err)
	}

	return nil
}

// DeleteMetadataEntryAsync deletes metadata of the corresponding organization with the given key, and returns
// a task.
// Deprecated: Use AdminOrg.DeleteMetadataEntryWithDomainAsync instead
func (adminOrg *AdminOrg) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, adminOrg.client, adminOrg.AdminOrg.HREF, adminOrg.AdminOrg.Name, key, false)
}

// AddMetadataEntry adds metadata key/value pair provided as input to the corresponding independent disk and waits for completion.
// Deprecated: Use Disk.AddMetadataEntryWithVisibility instead
func (disk *Disk) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := disk.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryAsync adds metadata key/value pair provided as input to the corresponding independent disk and returns a task.
// Deprecated: Use Disk.AddMetadataEntryWithVisibilityAsync instead
func (disk *Disk) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, disk.client, typedValue, key, value, disk.Disk.HREF)
}

// MergeMetadataAsync merges Disk metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use Disk.MergeMetadataWithMetadataValuesAsync
func (disk *Disk) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, disk.client, typedValue, metadata, disk.Disk.HREF)
}

// MergeMetadata merges Disk metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use Disk.MergeMetadataWithMetadataValues
func (disk *Disk) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := disk.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes metadata of the corresponding independent disk with the given key, and waits for completion
// Deprecated: Use Disk.DeleteMetadataEntryWithDomain instead
func (disk *Disk) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := disk.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for independent disk task: %s", err)
	}

	return nil
}

// DeleteMetadataEntryAsync deletes metadata of the corresponding independent disk with the given key, and returns
// a task.
// Deprecated: Use Disk.DeleteMetadataEntryWithDomainAsync instead
func (disk *Disk) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(disk.client, disk.Disk.HREF, disk.Disk.Name, key, false)
}

// AddMetadataEntry adds OrgVDCNetwork metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.AddMetadataEntryWithVisibility instead
func (orgVdcNetwork *OrgVDCNetwork) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := orgVdcNetwork.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryAsync adds OrgVDCNetwork metadata typedValue and key/value pair provided as input
// and returns the task.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.AddMetadataEntryWithVisibilityAsync instead
func (orgVdcNetwork *OrgVDCNetwork) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, orgVdcNetwork.client, typedValue, key, value, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF))
}

// MergeMetadataAsync merges OrgVDCNetwork metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.MergeMetadataWithMetadataValuesAsync
func (orgVdcNetwork *OrgVDCNetwork) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, orgVdcNetwork.client, typedValue, metadata, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF))
}

// MergeMetadata merges OrgVDCNetwork metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.MergeMetadataWithMetadataValues
func (orgVdcNetwork *OrgVDCNetwork) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := orgVdcNetwork.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntry adds CatalogItem metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use CatalogItem.AddMetadataEntryWithVisibility instead
func (catalogItem *CatalogItem) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := catalogItem.AddMetadataEntryAsync(ctx, typedValue, key, value)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// AddMetadataEntryAsync adds CatalogItem metadata typedValue and key/value pair provided as input
// and returns the task.
// Deprecated: Use CatalogItem.AddMetadataEntryWithVisibilityAsync instead
func (catalogItem *CatalogItem) AddMetadataEntryAsync(ctx context.Context, typedValue, key, value string) (Task, error) {
	return addMetadataDeprecated(ctx, catalogItem.client, typedValue, key, value, catalogItem.CatalogItem.HREF)
}

// MergeMetadataAsync merges CatalogItem metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use CatalogItem.MergeMetadataWithMetadataValuesAsync
func (catalogItem *CatalogItem) MergeMetadataAsync(ctx context.Context, typedValue string, metadata map[string]interface{}) (Task, error) {
	return mergeAllMetadataDeprecated(ctx, catalogItem.client, typedValue, metadata, catalogItem.CatalogItem.HREF)
}

// MergeMetadata merges CatalogItem metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// then waits for the task to complete.
// Deprecated: Use CatalogItem.MergeMetadataWithMetadataValues
func (catalogItem *CatalogItem) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := catalogItem.MergeMetadataAsync(ctx, typedValue, metadata)
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes CatalogItem metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use CatalogItem.DeleteMetadataEntryWithDomain instead
func (catalogItem *CatalogItem) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := catalogItem.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntryAsync deletes CatalogItem metadata depending on key provided as input
// and returns a task.
// Deprecated: Use CatalogItem.DeleteMetadataEntryWithDomainAsync instead
func (catalogItem *CatalogItem) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(catalogItem.client, catalogItem.CatalogItem.HREF, catalogItem.CatalogItem.Name, key, false)
}

// DeleteMetadataEntry deletes OrgVDCNetwork metadata depending on key provided as input
// and waits for the task to finish.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.DeleteMetadataEntryWithDomain instead
func (orgVdcNetwork *OrgVDCNetwork) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := orgVdcNetwork.DeleteMetadataEntryAsync(ctx, key)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntryAsync deletes OrgVDCNetwork metadata depending on key provided as input
// and returns a task.
// Note: Requires system administrator privileges.
// Deprecated: Use OrgVDCNetwork.DeleteMetadataEntryWithDomainAsync instead
func (orgVdcNetwork *OrgVDCNetwork) DeleteMetadataEntryAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, orgVdcNetwork.client, getAdminURL(orgVdcNetwork.OrgVDCNetwork.HREF), orgVdcNetwork.OrgVDCNetwork.Name, key, false)
}

// ----------------
// OpenAPI metadata functions

// AddMetadataEntry adds OpenApiOrgVdcNetwork metadata typedValue and key/value pair provided as input
// and waits for the task to finish.
// Deprecated: Use OpenApiOrgVdcNetwork.AddMetadataEntryWithVisibility instead
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) AddMetadataEntry(ctx context.Context, typedValue, key, value string) error {
	task, err := addMetadataDeprecated(ctx, openApiOrgVdcNetwork.client, typedValue, key, value, fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), strings.ReplaceAll(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID, "urn:vcloud:network:", "")))
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// MergeMetadata merges OpenApiOrgVdcNetwork metadata provided as a key-value map of type `typedValue` with the already present in VCD,
// and waits for the task to finish.
// Deprecated: Use OpenApiOrgVdcNetwork.MergeMetadataWithMetadataValues
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error {
	task, err := mergeAllMetadataDeprecated(ctx, openApiOrgVdcNetwork.client, typedValue, metadata, fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), strings.ReplaceAll(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID, "urn:vcloud:network:", "")))
	if err != nil {
		return err
	}
	return task.WaitTaskCompletion(ctx)
}

// DeleteMetadataEntry deletes OpenApiOrgVdcNetwork metadata depending on key provided as input
// and waits for the task to finish.
// Deprecated: Use OpenApiOrgVdcNetwork.DeleteMetadataEntryWithDomain
func (openApiOrgVdcNetwork *OpenApiOrgVdcNetwork) DeleteMetadataEntry(ctx context.Context, key string) error {
	task, err := deleteMetadata(ctx, openApiOrgVdcNetwork.client, fmt.Sprintf("%s/admin/network/%s", openApiOrgVdcNetwork.client.VCDHREF.String(), strings.ReplaceAll(openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.ID, "urn:vcloud:network:", "")), openApiOrgVdcNetwork.OpenApiOrgVdcNetwork.Name, key, false)
	if err != nil {
		return err
	}

	return task.WaitTaskCompletion(ctx)
}

// ----------------
// Generic private functions

// addMetadata adds metadata to an entity.
// The function supports passing a typedValue. Use one of the constants defined.
// Constants are types.MetadataStringValue, types.MetadataNumberValue, types.MetadataDateTimeValue and types.MetadataBooleanValue.
// Only tested with types.MetadataStringValue and types.MetadataNumberValue.
// Deprecated
func addMetadataDeprecated(ctx context.Context, client *Client, typedValue, key, value, requestUri string) (Task, error) {
	newMetadata := &types.MetadataValue{
		Xmlns: types.XMLNamespaceVCloud,
		Xsi:   types.XMLNamespaceXSI,
		TypedValue: &types.MetadataTypedValue{
			XsiType: typedValue,
			Value:   value,
		},
	}

	apiEndpoint := urlParseRequestURI(requestUri)
	apiEndpoint.Path += "/metadata/" + key

	// Return the task
	return client.ExecuteTaskRequest(ctx, apiEndpoint.String(), http.MethodPut,
		types.MimeMetaDataValue, "error adding metadata: %s", newMetadata)
}

// mergeAllMetadataDeprecated merges the metadata key-values provided as parameter with existing entity metadata
// Deprecated
func mergeAllMetadataDeprecated(ctx context.Context, client *Client, typedValue string, metadata map[string]interface{}, requestUri string) (Task, error) {
	var metadataToMerge []*types.MetadataEntry
	for key, value := range metadata {
		metadataToMerge = append(metadataToMerge, &types.MetadataEntry{
			Xmlns: types.XMLNamespaceVCloud,
			Xsi:   types.XMLNamespaceXSI,
			Key:   key,
			TypedValue: &types.MetadataTypedValue{
				XsiType: typedValue,
				Value:   value.(string),
			},
		})
	}

	newMetadata := &types.Metadata{
		Xmlns:         types.XMLNamespaceVCloud,
		Xsi:           types.XMLNamespaceXSI,
		MetadataEntry: metadataToMerge,
	}

	apiEndpoint := urlParseRequestURI(requestUri)
	apiEndpoint.Path += "/metadata"

	// Return the task
	return client.ExecuteTaskRequest(ctx, apiEndpoint.String(), http.MethodPost,
		types.MimeMetaData, "error adding metadata: %s", newMetadata)
}

// ----------------
// Deprecations

// Deprecated: use VM.DeleteMetadataEntry.
func (vm *VM) DeleteMetadata(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vm.client, vm.VM.HREF, vm.VM.Name, key, false)
}

// Deprecated: use VM.AddMetadataEntry.
func (vm *VM) AddMetadata(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vm.client, types.MetadataStringValue, key, value, vm.VM.HREF)
}

// Deprecated: use Vdc.DeleteMetadataEntry.
func (vdc *Vdc) DeleteMetadata(ctx context.Context, key string) (Vdc, error) {
	task, err := deleteMetadata(ctx, vdc.client, getAdminURL(vdc.Vdc.HREF), vdc.Vdc.Name, key, false)
	if err != nil {
		return Vdc{}, err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return Vdc{}, err
	}

	err = vdc.Refresh(ctx)
	if err != nil {
		return Vdc{}, err
	}

	return *vdc, nil
}

// Deprecated: use Vdc.AddMetadataEntry.
func (vdc *Vdc) AddMetadata(ctx context.Context, key string, value string) (Vdc, error) {
	task, err := addMetadataDeprecated(ctx, vdc.client, types.MetadataStringValue, key, value, getAdminURL(vdc.Vdc.HREF))
	if err != nil {
		return Vdc{}, err
	}

	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return Vdc{}, err
	}

	err = vdc.Refresh(ctx)
	if err != nil {
		return Vdc{}, err
	}

	return *vdc, nil
}

// Deprecated: use Vdc.AddMetadataEntryAsync.
func (vdc *Vdc) AddMetadataAsync(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vdc.client, types.MetadataStringValue, key, value, getAdminURL(vdc.Vdc.HREF))
}

// Deprecated: use Vdc.DeleteMetadataEntryAsync.
func (vdc *Vdc) DeleteMetadataAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vdc.client, getAdminURL(vdc.Vdc.HREF), vdc.Vdc.Name, key, false)
}

// Deprecated: use VApp.DeleteMetadataEntry.
func (vapp *VApp) DeleteMetadata(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vapp.client, vapp.VApp.HREF, vapp.VApp.Name, key, false)
}

// Deprecated: use VApp.AddMetadataEntry
func (vapp *VApp) AddMetadata(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vapp.client, types.MetadataStringValue, key, value, vapp.VApp.HREF)
}

// Deprecated: use VAppTemplate.AddMetadataEntry.
func (vAppTemplate *VAppTemplate) AddMetadata(ctx context.Context, key string, value string) (*VAppTemplate, error) {
	task, err := vAppTemplate.AddMetadataAsync(ctx, key, value)
	if err != nil {
		return nil, err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error completing add metadata for vApp template task: %s", err)
	}

	err = vAppTemplate.Refresh(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing vApp template: %s", err)
	}

	return vAppTemplate, nil
}

// Deprecated: use VAppTemplate.AddMetadataEntryAsync.
func (vAppTemplate *VAppTemplate) AddMetadataAsync(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, vAppTemplate.client, types.MetadataStringValue, key, value, vAppTemplate.VAppTemplate.HREF)
}

// Deprecated: use VAppTemplate.DeleteMetadataEntry.
func (vAppTemplate *VAppTemplate) DeleteMetadata(ctx context.Context, key string) error {
	task, err := vAppTemplate.DeleteMetadataAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for vApp template task: %s", err)
	}

	return nil
}

// Deprecated: use VAppTemplate.DeleteMetadataEntryAsync.
func (vAppTemplate *VAppTemplate) DeleteMetadataAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, vAppTemplate.client, vAppTemplate.VAppTemplate.HREF, vAppTemplate.VAppTemplate.Name, key, false)
}

// Deprecated: use Media.AddMetadataEntry.
func (media *Media) AddMetadata(ctx context.Context, key string, value string) (*Media, error) {
	task, err := media.AddMetadataAsync(ctx, key, value)
	if err != nil {
		return nil, err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error completing add metadata for media item task: %s", err)
	}

	err = media.Refresh(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing media item: %s", err)
	}

	return media, nil
}

// Deprecated: use Media.AddMetadataEntryAsync.
func (media *Media) AddMetadataAsync(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, media.client, types.MetadataStringValue, key, value, media.Media.HREF)
}

// Deprecated: use Media.DeleteMetadataEntry.
func (media *Media) DeleteMetadata(ctx context.Context, key string) error {
	task, err := media.DeleteMetadataAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for media item task: %s", err)
	}

	return nil
}

// Deprecated: use Media.DeleteMetadataEntryAsync.
func (media *Media) DeleteMetadataAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, media.client, media.Media.HREF, media.Media.Name, key, false)
}

// GetMetadata returns MediaItem metadata.
// Deprecated: Use MediaRecord.GetMetadata.
func (mediaItem *MediaItem) GetMetadata(ctx context.Context) (*types.Metadata, error) {
	return getMetadata(ctx, mediaItem.vdc.client, mediaItem.MediaItem.HREF, mediaItem.MediaItem.Name)
}

// AddMetadata adds metadata key/value pair provided as input.
// Deprecated: Use MediaRecord.AddMetadata.
func (mediaItem *MediaItem) AddMetadata(ctx context.Context, key string, value string) (*MediaItem, error) {
	task, err := mediaItem.AddMetadataAsync(ctx, key, value)
	if err != nil {
		return nil, err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error completing add metadata for media item task: %s", err)
	}

	err = mediaItem.Refresh(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing media item: %s", err)
	}

	return mediaItem, nil
}

// Deprecated: use MediaItem.AddMetadataEntryAsync.
func (mediaItem *MediaItem) AddMetadataAsync(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, mediaItem.vdc.client, types.MetadataStringValue, key, value, mediaItem.MediaItem.HREF)
}

// DeleteMetadata deletes metadata depending on key provided as input from media item.
// Deprecated: Use MediaRecord.DeleteMetadata.
func (mediaItem *MediaItem) DeleteMetadata(ctx context.Context, key string) error {
	task, err := mediaItem.DeleteMetadataAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for media item task: %s", err)
	}

	return nil
}

// DeleteMetadataAsync deletes metadata depending on key provided as input from MediaItem.
// Deprecated: Use MediaRecord.DeleteMetadataAsync.
func (mediaItem *MediaItem) DeleteMetadataAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, mediaItem.vdc.client, mediaItem.MediaItem.HREF, mediaItem.MediaItem.Name, key, false)
}

// Deprecated: use MediaRecord.AddMetadataEntry.
func (mediaRecord *MediaRecord) AddMetadata(ctx context.Context, key string, value string) (*MediaRecord, error) {
	task, err := mediaRecord.AddMetadataAsync(ctx, key, value)
	if err != nil {
		return nil, err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error completing add metadata for media item task: %s", err)
	}

	err = mediaRecord.Refresh(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing media item: %s", err)
	}

	return mediaRecord, nil
}

// Deprecated: use MediaRecord.AddMetadataEntryAsync.
func (mediaRecord *MediaRecord) AddMetadataAsync(ctx context.Context, key string, value string) (Task, error) {
	return addMetadataDeprecated(ctx, mediaRecord.client, types.MetadataStringValue, key, value, mediaRecord.MediaRecord.HREF)
}

// Deprecated: use MediaRecord.DeleteMetadataEntry.
func (mediaRecord *MediaRecord) DeleteMetadata(ctx context.Context, key string) error {
	task, err := mediaRecord.DeleteMetadataAsync(ctx, key)
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion(ctx)
	if err != nil {
		return fmt.Errorf("error completing delete metadata for media item task: %s", err)
	}

	return nil
}

// Deprecated: use MediaRecord.DeleteMetadataEntryAsync.
func (mediaRecord *MediaRecord) DeleteMetadataAsync(ctx context.Context, key string) (Task, error) {
	return deleteMetadata(ctx, mediaRecord.client, mediaRecord.MediaRecord.HREF, mediaRecord.MediaRecord.Name, key, false)
}
