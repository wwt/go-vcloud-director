//go:build network || nsxt || functional || openapi || ALL

/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_GetAllNsxtEdgeClusters(check *C) {
	skipNoNsxtConfiguration(vcd, check)

	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	nsxtVdc, err := vcd.org.GetVDCByNameOrId(ctx, vcd.config.VCD.Nsxt.Vdc, true)
	check.Assert(err, IsNil)

	edgeClusters, err := nsxtVdc.GetAllNsxtEdgeClusters(ctx, nil)
	check.Assert(err, IsNil)
	check.Assert(edgeClusters, NotNil)
	check.Assert(len(edgeClusters) > 0, Equals, true)

	queryParams := url.Values{}
	queryParams.Add("filter", fmt.Sprintf("orgVdcId==%s", nsxtVdc.Vdc.ID))
	allEdgeClusters, err := vcd.client.GetAllNsxtEdgeClusters(ctx, queryParams)
	check.Assert(err, IsNil)
	check.Assert(allEdgeClusters, NotNil)
	check.Assert(len(allEdgeClusters) > 0, Equals, true)
}

func (vcd *TestVCD) Test_GetNsxtEdgeClusterByName(check *C) {
	skipNoNsxtConfiguration(vcd, check)

	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	nsxtVdc, err := vcd.org.GetVDCByNameOrId(ctx, vcd.config.VCD.Nsxt.Vdc, true)
	check.Assert(err, IsNil)

	edgeCluster, err := nsxtVdc.GetNsxtEdgeClusterByName(ctx, vcd.config.VCD.Nsxt.NsxtEdgeCluster)
	check.Assert(err, IsNil)
	check.Assert(edgeCluster, NotNil)
	check.Assert(edgeCluster.NsxtEdgeCluster.Name, Equals, vcd.config.VCD.Nsxt.NsxtEdgeCluster)

}
