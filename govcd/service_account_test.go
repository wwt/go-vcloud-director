//go:build api || functional || ALL

/*
 * Copyright 2023 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_ServiceAccount(check *C) {
	isApiTokenEnabled, err := vcd.client.Client.VersionEqualOrGreater(ctx, "10.4.0", 3)
	check.Assert(err, IsNil)
	if !isApiTokenEnabled {
		check.Skip("This test requires VCD 10.4.0 or greater")
	}

	serviceAccount, err := vcd.client.CreateServiceAccount(ctx,
		vcd.config.VCD.Org,
		check.TestName(),
		"urn:vcloud:role:vApp%20Author",
		"12345678-1234-1234-1234-1234567890ab",
		"",
		"",
	)
	check.Assert(err, IsNil)
	check.Assert(serviceAccount, NotNil)
	check.Assert(serviceAccount.ServiceAccount.Status, Equals, "CREATED")

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointServiceAccounts
	check.Assert(err, IsNil)

	AddToCleanupListOpenApi(check.TestName(), check.TestName(), endpoint+serviceAccount.ServiceAccount.ID)

	err = serviceAccount.Authorize(ctx)
	check.Assert(err, IsNil)

	err = serviceAccount.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccount.ServiceAccount.Status, Equals, "REQUESTED")

	err = serviceAccount.Grant(ctx)
	check.Assert(err, IsNil)

	err = serviceAccount.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccount.ServiceAccount.Status, Equals, "GRANTED")

	_, err = serviceAccount.GetInitialApiToken(ctx)
	check.Assert(err, IsNil)

	err = serviceAccount.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccount.ServiceAccount.Status, Equals, "ACTIVE")

	err = serviceAccount.Revoke(ctx)
	check.Assert(err, IsNil)

	err = serviceAccount.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccount.ServiceAccount.Status, Equals, "CREATED")

	err = serviceAccount.Delete(ctx)
	check.Assert(err, IsNil)

	org, err := vcd.client.GetOrgByName(ctx, vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	notFound, err := org.GetServiceAccountById(ctx, serviceAccount.ServiceAccount.ID)
	check.Assert(err, NotNil)
	check.Assert(notFound, IsNil)
}

func (vcd *TestVCD) Test_ServiceAccount_SysOrg(check *C) {
	isApiTokenEnabled, err := vcd.client.Client.VersionEqualOrGreater(ctx, "10.4.0", 3)
	check.Assert(err, IsNil)
	if !isApiTokenEnabled {
		check.Skip("This test requires VCD 10.4.0 or greater")
	}

	if !vcd.org.client.IsSysAdmin {
		check.Skip("This test requires System Administrator role")
	}

	serviceAccountSysOrg, err := vcd.client.CreateServiceAccount(ctx,
		vcd.config.Provider.SysOrg,
		check.TestName(),
		"urn:vcloud:role:System%20Administrator",
		"12345678-1234-1234-1234-1234567890ab",
		"",
		"",
	)
	check.Assert(err, IsNil)
	check.Assert(serviceAccountSysOrg, NotNil)
	check.Assert(serviceAccountSysOrg.ServiceAccount.Status, Equals, "CREATED")

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointServiceAccounts
	check.Assert(err, IsNil)

	AddToCleanupListOpenApi(check.TestName(), check.TestName(), endpoint+serviceAccountSysOrg.ServiceAccount.ID)

	err = serviceAccountSysOrg.Authorize(ctx)
	check.Assert(err, IsNil)

	err = serviceAccountSysOrg.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccountSysOrg.ServiceAccount.Status, Equals, "REQUESTED")

	err = serviceAccountSysOrg.Grant(ctx)
	check.Assert(err, IsNil)

	err = serviceAccountSysOrg.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccountSysOrg.ServiceAccount.Status, Equals, "GRANTED")

	_, err = serviceAccountSysOrg.GetInitialApiToken(ctx)
	check.Assert(err, IsNil)

	err = serviceAccountSysOrg.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccountSysOrg.ServiceAccount.Status, Equals, "ACTIVE")

	err = serviceAccountSysOrg.Revoke(ctx)
	check.Assert(err, IsNil)

	err = serviceAccountSysOrg.Refresh(ctx)
	check.Assert(err, IsNil)
	check.Assert(serviceAccountSysOrg.ServiceAccount.Status, Equals, "CREATED")

	err = serviceAccountSysOrg.Delete(ctx)
	check.Assert(err, IsNil)

	sysorg, err := vcd.client.GetOrgByName(ctx, vcd.config.Provider.SysOrg)
	check.Assert(err, IsNil)
	check.Assert(sysorg, NotNil)

	notFound, err := sysorg.GetServiceAccountById(ctx, serviceAccountSysOrg.ServiceAccount.ID)
	check.Assert(err, NotNil)
	check.Assert(notFound, IsNil)
}
