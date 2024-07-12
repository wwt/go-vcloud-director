//go:build vapp || vdc || metadata || functional || ALL

/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	. "gopkg.in/check.v1"
)

// TODO: All tests here are deprecated in favor of those present in "metadata_v2_test". Remove this file once go-vcloud-director v3.0 is released.

func (vcd *TestVCD) Test_AddAndDeleteMetadataForVdc(check *C) {
	vcd.skipIfNotSysAdmin(check)
	if vcd.config.VCD.Vdc == "" {
		check.Skip("skipping test because VDC name is empty")
	}

	fmt.Printf("Running: %s\n", check.TestName())

	// Add metadata
	vdc, err := vcd.vdc.AddMetadata(ctx, "key", "value")
	check.Assert(err, IsNil)
	check.Assert(vdc, Not(Equals), Vdc{})

	AddToCleanupList("key", "vdcMetaData", vcd.org.Org.Name+"|"+vcd.config.VCD.Vdc, check.TestName())

	// Check if metadata was added correctly
	metadata, err := vcd.vdc.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 1)
	check.Assert(metadata.MetadataEntry[0].Key, Equals, "key")
	check.Assert(metadata.MetadataEntry[0].TypedValue.Value, Equals, "value")

	// Remove metadata
	vdc, err = vcd.vdc.DeleteMetadata(ctx, "key")
	check.Assert(err, IsNil)
	check.Assert(vdc, Not(Equals), Vdc{})

	metadata, err = vcd.vdc.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 0)
}

func (vcd *TestVCD) Test_AddAndDeleteMetadataOnVapp(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	if vcd.skipVappTests {
		check.Skip("Skipping test because vApp was not successfully created at setup")
	}
	// Add metadata
	task, err := vcd.vapp.AddMetadata(ctx, "key", "value")
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)
	check.Assert(task.Task.Status, Equals, "success")

	// Check if metadata was added correctly
	metadata, err := vcd.vapp.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 1)
	check.Assert(metadata.MetadataEntry[0].Key, Equals, "key")
	check.Assert(metadata.MetadataEntry[0].TypedValue.Value, Equals, "value")

	// Remove metadata
	task, err = vcd.vapp.DeleteMetadata(ctx, "key")
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)
	check.Assert(task.Task.Status, Equals, "success")

	metadata, err = vcd.vapp.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 0)
}

func (vcd *TestVCD) Test_AddAndDeleteMetadataOnVm(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	if vcd.skipVappTests {
		check.Skip("Skipping test because vApp was not successfully created at setup")
	}

	// Find VApp
	if vcd.vapp.VApp == nil {
		check.Skip("skipping test because no vApp is found")
	}

	vapp := vcd.findFirstVapp(ctx)
	vmType, vmName := vcd.findFirstVm(vapp)
	if vmName == "" {
		check.Skip("skipping test because no VM is found")
	}

	vm := NewVM(&vcd.client.Client)
	vm.VM = &vmType

	existingMetadata, err := vm.GetMetadata(ctx)
	check.Assert(err, IsNil)

	// Add metadata
	task, err := vm.AddMetadata(ctx, "key", "value")
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)
	check.Assert(task.Task.Status, Equals, "success")

	// Check if metadata was added correctly
	metadata, err := vm.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, len(existingMetadata.MetadataEntry)+1)
	found := false
	for _, entry := range metadata.MetadataEntry {
		if entry.Key == "key" && entry.TypedValue.Value == "value" {
			found = true
		}
	}
	check.Assert(found, Equals, true)

	// Remove metadata
	task, err = vm.DeleteMetadata(ctx, "key")
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)
	check.Assert(task.Task.Status, Equals, "success")
	metadata, err = vm.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, len(existingMetadata.MetadataEntry))
}

func (vcd *TestVCD) Test_AddAndDeleteMetadataOnVAppTemplate(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	cat, err := vcd.org.GetCatalogByName(ctx, vcd.config.VCD.Catalog.Name, false)
	if err != nil {
		check.Skip("Test_DeleteMetadataOnCatalogItem: Catalog not found. Test can't proceed")
		return
	}

	if vcd.config.VCD.Catalog.CatalogItem == "" {
		check.Skip("Test_DeleteMetadataOnCatalogItem: Catalog Item not given. Test can't proceed")
	}

	catItem, err := cat.GetCatalogItemByName(ctx, vcd.config.VCD.Catalog.CatalogItem, false)
	check.Assert(err, IsNil)
	check.Assert(catItem, NotNil)
	check.Assert(catItem.CatalogItem.Name, Equals, vcd.config.VCD.Catalog.CatalogItem)

	vAppTemplate, err := catItem.GetVAppTemplate(ctx)
	check.Assert(err, IsNil)
	check.Assert(vAppTemplate, NotNil)
	check.Assert(vAppTemplate.VAppTemplate.Name, Equals, vcd.config.VCD.Catalog.CatalogItem)

	// Check how much metaData exist
	metadata, err := vAppTemplate.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	existingMetaDataCount := len(metadata.MetadataEntry)

	// Add metadata
	_, err = vAppTemplate.AddMetadata(ctx, "key", "value")
	check.Assert(err, IsNil)

	// Check if metadata was added correctly
	metadata, err = vAppTemplate.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+1)
	var foundEntry *types.MetadataEntry
	for _, entry := range metadata.MetadataEntry {
		if entry.Key == "key" {
			foundEntry = entry
		}
	}
	check.Assert(foundEntry, NotNil)
	check.Assert(foundEntry.Key, Equals, "key")
	check.Assert(foundEntry.TypedValue.Value, Equals, "value")

	// Remove metadata
	err = vAppTemplate.DeleteMetadata(ctx, "key")
	check.Assert(err, IsNil)
	metadata, err = vAppTemplate.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 0)
}

func (vcd *TestVCD) Test_AddAndDeleteMetadataOnMediaRecord(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	//prepare media item
	skipWhenMediaPathMissing(vcd, check)
	itemName := "TestAddMediaMetaDataEntry"

	org, err := vcd.client.GetAdminOrgByName(ctx, vcd.org.Org.Name)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	catalog, err := org.GetCatalogByName(ctx, vcd.config.VCD.Catalog.Name, false)
	check.Assert(err, IsNil)
	check.Assert(catalog, NotNil)
	check.Assert(catalog.Catalog.Name, Equals, vcd.config.VCD.Catalog.Name)

	uploadTask, err := catalog.UploadMediaImage(ctx, itemName, "upload from test", vcd.config.Media.MediaPath, 1024)
	check.Assert(err, IsNil)
	check.Assert(uploadTask, NotNil)
	err = uploadTask.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)

	AddToCleanupList(itemName, "mediaCatalogImage", vcd.org.Org.Name+"|"+vcd.config.VCD.Catalog.Name, check.TestName())

	err = vcd.org.Refresh(ctx)
	check.Assert(err, IsNil)

	mediaRecord, err := catalog.QueryMedia(ctx, itemName)
	check.Assert(err, IsNil)
	check.Assert(mediaRecord, NotNil)
	check.Assert(mediaRecord.MediaRecord.Name, Equals, itemName)

	// Add metadata
	_, err = mediaRecord.AddMetadata(ctx, "key", "value")
	check.Assert(err, IsNil)

	// Check if metadata was added correctly
	metadata, err := mediaRecord.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 1)
	check.Assert(metadata.MetadataEntry[0].Key, Equals, "key")
	check.Assert(metadata.MetadataEntry[0].TypedValue.Value, Equals, "value")

	// Remove metadata
	err = mediaRecord.DeleteMetadata(ctx, "key")
	check.Assert(err, IsNil)
	metadata, err = mediaRecord.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 0)

	// Remove catalog item so far other tests don't fail
	task, err := mediaRecord.Delete(ctx)
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)
}

func (vcd *TestVCD) Test_MetadataOnAdminCatalogCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	var catalogName string = check.TestName()

	org, err := vcd.client.GetAdminOrgByName(ctx, vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	catalog, err := vcd.org.CreateCatalog(ctx, catalogName, catalogName)
	check.Assert(err, IsNil)
	check.Assert(catalog, NotNil)
	AddToCleanupList(catalogName, "catalog", org.AdminOrg.Name, catalogName)

	adminCatalog, err := org.GetAdminCatalogByName(ctx, catalogName, true)
	check.Assert(err, IsNil)
	check.Assert(adminCatalog, NotNil)

	testMetadataCRUDActionsDeprecated(adminCatalog, check, func() {
		metadata, err := catalog.GetMetadata(ctx)
		check.Assert(err, IsNil)
		check.Assert(metadata, NotNil)
		check.Assert(len(metadata.MetadataEntry), Equals, 1)
		var foundEntry *types.MetadataEntry
		for _, entry := range metadata.MetadataEntry {
			if entry.Key == "key" {
				foundEntry = entry
			}
		}
		check.Assert(foundEntry, NotNil)
		check.Assert(foundEntry.Key, Equals, "key")
		check.Assert(foundEntry.TypedValue.Value, Equals, "value")
	})
	err = catalog.Delete(ctx, true, true)
	check.Assert(err, IsNil)
}

func (vcd *TestVCD) Test_MetadataEntryForVdcCRUD(check *C) {
	vcd.skipIfNotSysAdmin(check)
	if vcd.config.VCD.Vdc == "" {
		check.Skip("skipping test because VDC name is empty")
	}

	fmt.Printf("Running: %s\n", check.TestName())

	testMetadataCRUDActionsDeprecated(vcd.vdc, check, nil)
}

func (vcd *TestVCD) Test_MetadataEntryOnVappCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	if vcd.skipVappTests {
		check.Skip("Skipping test because vApp was not successfully created at setup")
	}
	testMetadataCRUDActionsDeprecated(vcd.vapp, check, nil)
}

func (vcd *TestVCD) Test_MetadataEntryOnVmCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	if vcd.skipVappTests {
		check.Skip("Skipping test because vapp was not successfully created at setup")
	}

	// Find VApp
	if vcd.vapp.VApp == nil {
		check.Skip("skipping test because no vApp is found")
	}

	vapp := vcd.findFirstVapp(ctx)
	vmType, vmName := vcd.findFirstVm(vapp)
	if vmName == "" {
		check.Skip("skipping test because no VM is found")
	}

	vm := NewVM(&vcd.client.Client)
	vm.VM = &vmType

	testMetadataCRUDActionsDeprecated(vm, check, nil)
}

func (vcd *TestVCD) Test_MetadataEntryOnVAppTemplateCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	cat, err := vcd.org.GetCatalogByName(ctx, vcd.config.VCD.Catalog.Name, false)
	if err != nil {
		check.Skip("Test_DeleteMetadataOnCatalogItem: Catalog not found. Test can't proceed")
		return
	}

	if vcd.config.VCD.Catalog.CatalogItem == "" {
		check.Skip("Test_DeleteMetadataOnCatalogItem: Catalog Item not given. Test can't proceed")
	}

	catItem, err := cat.GetCatalogItemByName(ctx, vcd.config.VCD.Catalog.CatalogItem, false)
	check.Assert(err, IsNil)
	check.Assert(catItem, NotNil)
	check.Assert(catItem.CatalogItem.Name, Equals, vcd.config.VCD.Catalog.CatalogItem)

	vAppTemplate, err := catItem.GetVAppTemplate(ctx)
	check.Assert(err, IsNil)
	check.Assert(vAppTemplate, NotNil)
	check.Assert(vAppTemplate.VAppTemplate.Name, Equals, vcd.config.VCD.Catalog.CatalogItem)

	testMetadataCRUDActionsDeprecated(&vAppTemplate, check, nil)
}

func (vcd *TestVCD) Test_MetadataEntryOnMediaRecordCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	//prepare media item
	skipWhenMediaPathMissing(vcd, check)
	itemName := "TestAddMediaMetaData"

	org, err := vcd.client.GetAdminOrgByName(ctx, vcd.org.Org.Name)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	catalog, err := org.GetCatalogByName(ctx, vcd.config.VCD.Catalog.Name, false)
	check.Assert(err, IsNil)
	check.Assert(catalog, NotNil)
	check.Assert(catalog.Catalog.Name, Equals, vcd.config.VCD.Catalog.Name)

	uploadTask, err := catalog.UploadMediaImage(ctx, itemName, "upload from test", vcd.config.Media.MediaPath, 1024)
	check.Assert(err, IsNil)
	check.Assert(uploadTask, NotNil)
	err = uploadTask.WaitTaskCompletion(ctx)
	check.Assert(err, IsNil)

	AddToCleanupList(itemName, "mediaCatalogImage", vcd.org.Org.Name+"|"+vcd.config.VCD.Catalog.Name, check.TestName())

	err = vcd.org.Refresh(ctx)
	check.Assert(err, IsNil)

	mediaRecord, err := catalog.QueryMedia(ctx, itemName)
	check.Assert(err, IsNil)
	check.Assert(mediaRecord, NotNil)
	check.Assert(mediaRecord.MediaRecord.Name, Equals, itemName)

	testMetadataCRUDActionsDeprecated(mediaRecord, check, nil)
	task, err := mediaRecord.Delete(ctx)
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(context.Background())
	check.Assert(err, IsNil)
}

func (vcd *TestVCD) Test_MetadataOnAdminOrgCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	adminOrg, err := vcd.client.GetAdminOrgById(ctx, vcd.org.Org.ID)
	if err != nil {
		check.Skip("Test_AddMetadataOnAdminOrg: Organization (admin) not found. Test can't proceed")
		return
	}
	org, err := vcd.client.GetOrgById(ctx, vcd.org.Org.ID)
	if err != nil {
		check.Skip("Test_AddMetadataOnAdminOrg: Organization not found. Test can't proceed")
		return
	}

	testMetadataCRUDActionsDeprecated(adminOrg, check, func() {
		metadata, err := org.GetMetadata(ctx)
		check.Assert(err, IsNil)
		check.Assert(metadata, NotNil)
		check.Assert(len(metadata.MetadataEntry), Equals, 1)
		var foundEntry *types.MetadataEntry
		for _, entry := range metadata.MetadataEntry {
			if entry.Key == "key" {
				foundEntry = entry
			}
		}
		check.Assert(foundEntry, NotNil)
		check.Assert(foundEntry.Key, Equals, "key")
		check.Assert(foundEntry.TypedValue.Value, Equals, "value")
	})

}

func (vcd *TestVCD) Test_MetadataOnIndependentDiskCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	// Create disk
	diskCreateParamsDisk := &types.Disk{
		Name:        TestCreateDisk,
		SizeMb:      11,
		Description: TestCreateDisk,
	}

	diskCreateParams := &types.DiskCreateParams{
		Disk: diskCreateParamsDisk,
	}

	task, err := vcd.vdc.CreateDisk(ctx, diskCreateParams)
	check.Assert(err, IsNil)

	diskHREF := task.Task.Owner.HREF
	PrependToCleanupList(diskHREF, "disk", "", check.TestName())

	err = task.WaitTaskCompletion(context.Background())
	check.Assert(err, IsNil)

	disk, err := vcd.vdc.GetDiskByHref(ctx, diskHREF)
	check.Assert(err, IsNil)

	testMetadataCRUDActionsDeprecated(disk, check, nil)

	task, err = disk.Delete(ctx)
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion(context.Background())
	check.Assert(err, IsNil)
}

func (vcd *TestVCD) Test_MetadataOnVdcNetworkCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	net, err := vcd.vdc.GetOrgVdcNetworkByName(ctx, vcd.config.VCD.Network.Net1, false)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Network %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.Network.Net1))
		return
	}

	testMetadataCRUDActionsDeprecated(net, check, nil)
}

func (vcd *TestVCD) Test_MetadataOnCatalogItemCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	catalog, err := vcd.org.GetCatalogByName(ctx, vcd.config.VCD.Catalog.Name, false)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Catalog %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.Catalog.Name))
		return
	}

	catalogItem, err := catalog.GetCatalogItemByName(ctx, vcd.config.VCD.Catalog.CatalogItem, false)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Catalog item %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.Catalog.CatalogItem))
		return
	}

	testMetadataCRUDActionsDeprecated(catalogItem, check, nil)
}

func (vcd *TestVCD) Test_MetadataOnProviderVdcCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	vcd.skipIfNotSysAdmin(check)
	providerVdc, err := vcd.client.GetProviderVdcByName(ctx, vcd.config.VCD.NsxtProviderVdc.Name)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Provider VDC %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.NsxtProviderVdc.Name))
		return
	}

	testMetadataCRUDActionsDeprecated(providerVdc, check, nil)
}

func (vcd *TestVCD) Test_MetadataOnOpenApiOrgVdcNetworkCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	net, err := vcd.vdc.GetOpenApiOrgVdcNetworkByName(ctx, vcd.config.VCD.Network.Net1)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Network %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.Network.Net1))
		return
	}

	testMetadataCRUDActionsDeprecated(net, check, nil)
}

func (vcd *TestVCD) Test_MetadataByHrefCRUD(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())
	vcd.skipIfNotSysAdmin(check)
	storageProfileRef, err := vcd.vdc.FindStorageProfileReference(ctx, vcd.config.VCD.StorageProfile.SP1)
	if err != nil {
		check.Skip(fmt.Sprintf("%s: Storage Profile %s not found. Test can't proceed", check.TestName(), vcd.config.VCD.StorageProfile.SP1))
		return
	}

	storageProfileAdminHref := strings.ReplaceAll(storageProfileRef.HREF, "api", "api/admin")

	// Check how much metadata exists
	metadata, err := vcd.client.GetMetadataByHref(ctx, storageProfileAdminHref)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	existingMetaDataCount := len(metadata.MetadataEntry)

	// Add metadata
	err = vcd.client.AddMetadataEntryByHref(ctx, storageProfileAdminHref, types.MetadataStringValue, "key", "value")
	check.Assert(err, IsNil)

	// Check if metadata was added correctly
	metadata, err = vcd.client.GetMetadataByHref(ctx, storageProfileAdminHref)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+1)
	var foundEntry *types.MetadataEntry
	for _, entry := range metadata.MetadataEntry {
		if entry.Key == "key" {
			foundEntry = entry
		}
	}
	check.Assert(foundEntry, NotNil)
	check.Assert(foundEntry.Key, Equals, "key")
	check.Assert(foundEntry.TypedValue.Value, Equals, "value")

	// Check the same without admin privileges
	metadata, err = vcd.client.GetMetadataByHref(ctx, storageProfileRef.HREF)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+1)
	for _, entry := range metadata.MetadataEntry {
		if entry.Key == "key" {
			foundEntry = entry
		}
	}
	check.Assert(foundEntry, NotNil)
	check.Assert(foundEntry.Key, Equals, "key")
	check.Assert(foundEntry.TypedValue.Value, Equals, "value")

	// Merge updated metadata with a new entry
	err = vcd.client.MergeMetadataByHref(ctx, storageProfileAdminHref, types.MetadataStringValue, map[string]interface{}{
		"key":  "valueUpdated",
		"key2": "value2",
	})
	check.Assert(err, IsNil)

	// Check that the first key was updated and the second, created
	metadata, err = vcd.client.GetMetadataByHref(ctx, storageProfileRef.HREF)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+2)
	for _, entry := range metadata.MetadataEntry {
		switch entry.Key {
		case "key":
			check.Assert(entry.TypedValue.Value, Equals, "valueUpdated")
		case "key2":
			check.Assert(entry.TypedValue.Value, Equals, "value2")
		}
	}

	// Delete the metadata
	err = vcd.client.DeleteMetadataEntryByHref(ctx, storageProfileAdminHref, "key")
	check.Assert(err, IsNil)
	err = vcd.client.DeleteMetadataEntryByHref(ctx, storageProfileAdminHref, "key2")
	check.Assert(err, IsNil)
	// Check if metadata was deleted correctly
	metadata, err = vcd.client.GetMetadataByHref(ctx, storageProfileAdminHref)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, 0)
}

type metadataCompatibleDeprecated interface {
	GetMetadata(ctx context.Context) (*types.Metadata, error)
	AddMetadataEntry(ctx context.Context, typedValue, key, value string) error
	MergeMetadata(ctx context.Context, typedValue string, metadata map[string]interface{}) error
	DeleteMetadataEntry(ctx context.Context, key string) error
}

func testMetadataCRUDActionsDeprecated(resource metadataCompatibleDeprecated, check *C, extraCheck func()) {
	// Check how much metadata exists
	metadata, err := resource.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	existingMetaDataCount := len(metadata.MetadataEntry)

	// Add metadata
	err = resource.AddMetadataEntry(ctx, types.MetadataStringValue, "key", "value")
	check.Assert(err, IsNil)

	// Check if metadata was added correctly
	metadata, err = resource.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+1)
	var foundEntry *types.MetadataEntry
	for _, entry := range metadata.MetadataEntry {
		if entry.Key == "key" {
			foundEntry = entry
		}
	}
	check.Assert(foundEntry, NotNil)
	check.Assert(foundEntry.Key, Equals, "key")
	check.Assert(foundEntry.TypedValue.Value, Equals, "value")

	if extraCheck != nil {
		extraCheck()
	}

	// Merge updated metadata with a new entry
	err = resource.MergeMetadata(ctx, types.MetadataStringValue, map[string]interface{}{
		"key":  "valueUpdated",
		"key2": "value2",
	})
	check.Assert(err, IsNil)

	// Check that the first key was updated and the second, created
	metadata, err = resource.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount+2)
	for _, entry := range metadata.MetadataEntry {
		switch entry.Key {
		case "key":
			check.Assert(entry.TypedValue.Value, Equals, "valueUpdated")
		case "key2":
			check.Assert(entry.TypedValue.Value, Equals, "value2")
		}
	}

	err = resource.DeleteMetadataEntry(ctx, "key")
	check.Assert(err, IsNil)
	err = resource.DeleteMetadataEntry(ctx, "key2")
	check.Assert(err, IsNil)
	// Check if metadata was deleted correctly
	metadata, err = resource.GetMetadata(ctx)
	check.Assert(err, IsNil)
	check.Assert(metadata, NotNil)
	check.Assert(len(metadata.MetadataEntry), Equals, existingMetaDataCount)
}
