//go:build network || nsxt || functional || openapi || ALL

package govcd

import (
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
)

// Test_NsxEdgeBgpIpPrefixList tests CRUD operations for NSX-T Edge Gateway BGP IP Prefix Lists
func (vcd *TestVCD) Test_NsxEdgeBgpIpPrefixList(check *C) {
	skipNoNsxtConfiguration(vcd, check)
	skipOpenApiEndpointTest(ctx, vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointEdgeBgpConfigPrefixLists)

	org, err := vcd.client.GetOrgByName(ctx, vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	nsxtVdc, err := org.GetVDCByName(ctx, vcd.config.VCD.Nsxt.Vdc, false)
	check.Assert(err, IsNil)
	edge, err := nsxtVdc.GetNsxtEdgeGatewayByName(ctx, vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	// Switch Edge Gateway to use dedicated uplink for the time of this test and then turn it off
	err = switchEdgeGatewayDedication(edge, true) // Turn on Dedicated Tier 0 gateway
	check.Assert(err, IsNil)
	defer switchEdgeGatewayDedication(edge, false) // Turn off Dedicated Tier 0 gateway

	// Create a new BGP IP Prefix List
	bgpIpPrefixList := &types.EdgeBgpIpPrefixList{
		Name:        check.TestName(),
		Description: "test-description",
		Prefixes: []types.EdgeBgpConfigPrefixListPrefixes{
			{
				Network: "1.1.1.0/24",
				Action:  "PERMIT",
			},
			{
				Network:            "2.1.0.0/16",
				Action:             "PERMIT",
				LessThanEqualTo:    29,
				GreaterThanEqualTo: 24,
			},
		},
	}

	bgpIpPrefix, err := edge.CreateBgpIpPrefixList(ctx, bgpIpPrefixList)
	check.Assert(err, IsNil)
	check.Assert(bgpIpPrefix, NotNil)

	// Get all BGP IP Prefix Lists
	bgpIpPrefixLists, err := edge.GetAllBgpIpPrefixLists(ctx, nil)
	check.Assert(err, IsNil)
	check.Assert(bgpIpPrefixLists, NotNil)
	check.Assert(len(bgpIpPrefixLists), Equals, 1)
	check.Assert(bgpIpPrefixLists[0].EdgeBgpIpPrefixList.Name, Equals, bgpIpPrefixList.Name)

	// Get By Name
	bgpPrefixListByName, err := edge.GetBgpIpPrefixListByName(ctx, bgpIpPrefixList.Name)
	check.Assert(err, IsNil)
	check.Assert(bgpPrefixListByName, NotNil)

	// Get By Id
	bgpPrefixListById, err := edge.GetBgpIpPrefixListById(ctx, bgpIpPrefix.EdgeBgpIpPrefixList.ID)
	check.Assert(err, IsNil)
	check.Assert(bgpPrefixListById, NotNil)

	// Update
	bgpIpPrefixList.Name = check.TestName() + "-updated"
	bgpIpPrefixList.Description = "test-description-updated"
	bgpIpPrefixList.ID = bgpIpPrefixLists[0].EdgeBgpIpPrefixList.ID

	updatedBgpIpPrefixList, err := bgpIpPrefixLists[0].Update(ctx, bgpIpPrefixList)
	check.Assert(err, IsNil)
	check.Assert(updatedBgpIpPrefixList, NotNil)

	check.Assert(updatedBgpIpPrefixList.EdgeBgpIpPrefixList.ID, Equals, bgpIpPrefixLists[0].EdgeBgpIpPrefixList.ID)

	// Delete
	err = bgpIpPrefixLists[0].Delete(ctx)
	check.Assert(err, IsNil)

	// Try to get once again and ensure it is not there
	notFoundByName, err := edge.GetBgpIpPrefixListByName(ctx, bgpIpPrefixList.Name)
	check.Assert(err, NotNil)
	check.Assert(notFoundByName, IsNil)

	notFoundById, err := edge.GetBgpIpPrefixListById(ctx, bgpIpPrefix.EdgeBgpIpPrefixList.ID)
	check.Assert(err, NotNil)
	check.Assert(notFoundById, IsNil)

}
