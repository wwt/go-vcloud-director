//go:build nsxt || alb || functional || ALL

package govcd

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_GetAllAlbImportableClouds(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipNoNsxtAlbConfiguration(vcd, check)

	albController := spawnAlbController(vcd, check)

	// Test client function with explicit ALB Controller ID requirement
	clientImportableClouds, err := vcd.client.GetAllAlbImportableClouds(ctx, albController.NsxtAlbController.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(len(clientImportableClouds) > 0, Equals, true)

	// Test functions attached directly to NsxtAlbController
	controllerImportableClouds, err := albController.GetAllAlbImportableClouds(ctx, nil)
	check.Assert(err, IsNil)
	check.Assert(len(controllerImportableClouds) > 0, Equals, true)

	controllerImportableCloudByName, err := albController.GetAlbImportableCloudByName(ctx, vcd.config.VCD.Nsxt.NsxtAlbImportableCloud)
	check.Assert(err, IsNil)
	check.Assert(controllerImportableCloudByName, NotNil)
	check.Assert(controllerImportableCloudByName.NsxtAlbImportableCloud.ID, Equals, controllerImportableClouds[0].NsxtAlbImportableCloud.ID)

	// Cleanup
	err = albController.Delete(ctx)
	check.Assert(err, IsNil)
}
