//go:build nsxt || alb || functional || ALL

package govcd

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_AlbVirtualService(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipNoNsxtAlbConfiguration(vcd, check)
	skipOpenApiEndpointTest(ctx, vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointAlbEdgeGateway)

	// Setup prerequisite components
	controller, cloud, seGroup, edge, seGroupAssignment, albPool := setupAlbVirtualServicePrerequisites(check, vcd)

	// Setup Org user and connection
	adminOrg, err := vcd.client.GetAdminOrgByName(ctx, vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	orgUserVcdClient, orgUser, err := newOrgUserConnection(adminOrg, "alb-virtual-service-testing", "CHANGE-ME", vcd.config.Provider.Url, true)
	check.Assert(err, IsNil)

	printVerbose("# Running tests as Sysadmin user\n")
	// Run tests with System user
	testMinimalVirtualServiceConfigHTTP(check, edge, albPool, seGroup, vcd, vcd.client)
	testVirtualServiceConfigWithCertHTTPS(check, edge, albPool, seGroup, vcd, vcd.client)
	testMinimalVirtualServiceConfigL4(check, edge, albPool, seGroup, vcd, vcd.client)
	testMinimalVirtualServiceConfigL4TLS(check, edge, albPool, seGroup, vcd, vcd.client)

	printVerbose("# Running tests as Org user\n")
	// Run tests with Org admin user
	testMinimalVirtualServiceConfigHTTP(check, edge, albPool, seGroup, vcd, orgUserVcdClient)
	testVirtualServiceConfigWithCertHTTPS(check, edge, albPool, seGroup, vcd, orgUserVcdClient)
	testMinimalVirtualServiceConfigL4(check, edge, albPool, seGroup, vcd, orgUserVcdClient)
	testMinimalVirtualServiceConfigL4TLS(check, edge, albPool, seGroup, vcd, orgUserVcdClient)

	// Test 10.4.1 Transparent mode on VCD >= 10.4.1
	if vcd.client.Client.APIVCDMaxVersionIs(ctx, ">= 37.1") {
		printVerbose("# Running 10.4.1+ tests as Sysadmin user\n")

		printVerbose("## Creating ALB Pool with Member Group (VCD 10.4.1+) as Sysadmin\n")
		ipSet, poolWithMemberGroup := setupAlbPoolFirewallGroupMembers(check, vcd, edge)

		testMinimalVirtualServiceConfigHTTPTransparent(check, edge, poolWithMemberGroup, seGroup, vcd, vcd.client, true)
		testMinimalVirtualServiceConfigHTTPTransparent(check, edge, poolWithMemberGroup, seGroup, vcd, vcd.client, false)

		printVerbose("# Running 10.4.1+ tests as Org user\n")

		printVerbose("## Creating ALB Pool with Member Group (VCD 10.4.1+) as Org user\n")
		testMinimalVirtualServiceConfigHTTPTransparent(check, edge, poolWithMemberGroup, seGroup, vcd, orgUserVcdClient, true)
		testMinimalVirtualServiceConfigHTTPTransparent(check, edge, poolWithMemberGroup, seGroup, vcd, orgUserVcdClient, false)

		// cleanup ipset and pool membership
		err = poolWithMemberGroup.Delete(ctx)
		check.Assert(err, IsNil)

		err = ipSet.Delete(ctx)
		check.Assert(err, IsNil)
	}

	// teardown prerequisites
	tearDownAlbVirtualServicePrerequisites(check, albPool, seGroupAssignment, edge, seGroup, cloud, controller)

	// cleanup Org user
	err = orgUser.Delete(ctx, true)
	check.Assert(err, IsNil)
}

func testMinimalVirtualServiceConfigHTTP(check *C, edge *NsxtEdgeGateway, pool *NsxtAlbPool, seGroup *NsxtAlbServiceEngineGroup, vcd *TestVCD, client *VCDClient) {
	virtualServiceConfig := &types.NsxtAlbVirtualService{
		Name:    check.TestName(),
		Enabled: takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTP",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart: takeIntAddress(80),
			},
		},
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	virtualServiceConfigUpdated := &types.NsxtAlbVirtualService{
		Name:        check.TestName(),
		Description: "Updated",
		Enabled:     takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTP",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart:  takeIntAddress(443),
				PortEnd:    takeIntAddress(449),
				SslEnabled: takeBoolPointer(false),
			},
			{
				PortStart:  takeIntAddress(2000),
				PortEnd:    takeIntAddress(2010),
				SslEnabled: takeBoolPointer(false),
			},
		},
		// Use Primary IP of Edge Gateway as virtual service IP
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
		//HealthStatus:          "",
		//HealthMessage:         "",
		//DetailedHealthMessage: "",
	}

	testAlbVirtualServiceConfig(check, vcd, "MinimalHTTP", virtualServiceConfig, virtualServiceConfigUpdated, client)
}

func testMinimalVirtualServiceConfigHTTPTransparent(check *C, edge *NsxtEdgeGateway, poolWithMemberGroup *NsxtAlbPool, seGroup *NsxtAlbServiceEngineGroup, vcd *TestVCD, client *VCDClient, trueOnCreate bool) {
	createTransparentMode := trueOnCreate
	updateTransparentMode := !createTransparentMode

	virtualServiceConfig := &types.NsxtAlbVirtualService{
		Name:                   check.TestName(),
		Enabled:                addrOf(true),
		TransparentModeEnabled: &createTransparentMode,
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTP",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: poolWithMemberGroup.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart: addrOf(80),
			},
		},
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	virtualServiceConfigUpdated := &types.NsxtAlbVirtualService{
		Name:                   check.TestName(),
		Description:            "Updated",
		Enabled:                addrOf(true),
		TransparentModeEnabled: &updateTransparentMode,
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTP",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: poolWithMemberGroup.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart:  addrOf(443),
				PortEnd:    addrOf(449),
				SslEnabled: addrOf(false),
			},
			{
				PortStart:  addrOf(2000),
				PortEnd:    addrOf(2010),
				SslEnabled: addrOf(false),
			},
		},
		// Use Primary IP of Edge Gateway as virtual service IP
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
		//HealthStatus:          "",
		//HealthMessage:         "",
		//DetailedHealthMessage: "",
	}

	testAlbVirtualServiceConfig(check, vcd, fmt.Sprintf("MinimalHTTPWithTransparentModeOnCreate%t", createTransparentMode), virtualServiceConfig, virtualServiceConfigUpdated, client)
}

func testMinimalVirtualServiceConfigL4(check *C, edge *NsxtEdgeGateway, pool *NsxtAlbPool, seGroup *NsxtAlbServiceEngineGroup, vcd *TestVCD, client *VCDClient) {
	virtualServiceConfig := &types.NsxtAlbVirtualService{
		Name:    check.TestName(),
		Enabled: takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "L4",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart: takeIntAddress(80),
			},
		},
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	virtualServiceConfigUpdated := &types.NsxtAlbVirtualService{
		Name:        check.TestName(),
		Description: "Updated",
		Enabled:     takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "L4",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart: takeIntAddress(443),
				TcpUdpProfile: &types.NsxtAlbVirtualServicePortTcpUdpProfile{
					SystemDefined: true,
					Type:          "TCP_PROXY",
				},
			},
			{
				PortStart: takeIntAddress(8443),
				PortEnd:   takeIntAddress(8445),
				TcpUdpProfile: &types.NsxtAlbVirtualServicePortTcpUdpProfile{
					SystemDefined: true,
					Type:          "TCP_FAST_PATH",
				},
			},
			{
				PortStart: takeIntAddress(9000),
				TcpUdpProfile: &types.NsxtAlbVirtualServicePortTcpUdpProfile{
					SystemDefined: true,
					Type:          "UDP_FAST_PATH",
				},
			},
			{
				PortStart: takeIntAddress(10000),
			},
		},
		// Use Primary IP of Edge Gateway as virtual service IP
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	testAlbVirtualServiceConfig(check, vcd, "L4", virtualServiceConfig, virtualServiceConfigUpdated, client)
}

func testMinimalVirtualServiceConfigL4TLS(check *C, edge *NsxtEdgeGateway, pool *NsxtAlbPool, seGroup *NsxtAlbServiceEngineGroup, vcd *TestVCD, client *VCDClient) {
	adminOrg, err := vcd.client.GetAdminOrgByName(ctx, vcd.org.Org.Name)
	check.Assert(err, IsNil)
	check.Assert(adminOrg, NotNil)

	certificateConfigWithPrivateKey := &types.CertificateLibraryItem{
		Alias:                check.TestName(),
		Certificate:          certificate,
		PrivateKey:           privateKey,
		PrivateKeyPassphrase: "test",
	}
	openApiEndpoint, err := getEndpointByVersion(ctx, &vcd.client.Client)
	check.Assert(err, IsNil)
	createdCertificate, err := adminOrg.AddCertificateToLibrary(ctx, certificateConfigWithPrivateKey)
	check.Assert(err, IsNil)
	PrependToCleanupListOpenApi(createdCertificate.CertificateLibrary.Alias, check.TestName(), openApiEndpoint+createdCertificate.CertificateLibrary.Id)

	virtualServiceConfig := &types.NsxtAlbVirtualService{
		Name:    check.TestName(),
		Enabled: takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "L4_TLS",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		CertificateRef:        &types.OpenApiReference{ID: createdCertificate.CertificateLibrary.Id},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart:  takeIntAddress(80),
				SslEnabled: takeBoolPointer(true),
			},
		},
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	virtualServiceConfigUpdated := &types.NsxtAlbVirtualService{
		Name:        check.TestName(),
		Description: "Updated",
		Enabled:     takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "L4_TLS",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		CertificateRef:        &types.OpenApiReference{ID: createdCertificate.CertificateLibrary.Id},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart:  takeIntAddress(443),
				SslEnabled: takeBoolPointer(true),
				TcpUdpProfile: &types.NsxtAlbVirtualServicePortTcpUdpProfile{
					SystemDefined: true,
					Type:          "TCP_PROXY", // The only possible type with L4_TLS
				},
			},
		},
		// Use Primary IP of Edge Gateway as virtual service IP
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	testAlbVirtualServiceConfig(check, vcd, "L4-TLS", virtualServiceConfig, virtualServiceConfigUpdated, client)

	err = createdCertificate.Delete(ctx)
	check.Assert(err, IsNil)
}

func testVirtualServiceConfigWithCertHTTPS(check *C, edge *NsxtEdgeGateway, pool *NsxtAlbPool, seGroup *NsxtAlbServiceEngineGroup, vcd *TestVCD, client *VCDClient) {
	adminOrg, err := vcd.client.GetAdminOrgByName(ctx, vcd.org.Org.Name)
	check.Assert(err, IsNil)
	check.Assert(adminOrg, NotNil)

	certificateConfigWithPrivateKey := &types.CertificateLibraryItem{
		Alias:                check.TestName(),
		Certificate:          certificate,
		PrivateKey:           privateKey,
		PrivateKeyPassphrase: "test",
	}

	openApiEndpoint, err := getEndpointByVersion(ctx, &vcd.client.Client)
	check.Assert(err, IsNil)
	createdCertificate, err := adminOrg.AddCertificateToLibrary(ctx, certificateConfigWithPrivateKey)
	check.Assert(err, IsNil)
	PrependToCleanupListOpenApi(createdCertificate.CertificateLibrary.Alias, check.TestName(), openApiEndpoint+createdCertificate.CertificateLibrary.Id)

	virtualServiceConfig := &types.NsxtAlbVirtualService{
		Name:    check.TestName(),
		Enabled: takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTPS",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		CertificateRef:        &types.OpenApiReference{ID: createdCertificate.CertificateLibrary.Id},
		ServicePorts:          []types.NsxtAlbVirtualServicePort{{PortStart: takeIntAddress(80), SslEnabled: takeBoolPointer(true)}},
		VirtualIpAddress:      edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	virtualServiceConfigUpdated := &types.NsxtAlbVirtualService{
		Name:        check.TestName(),
		Description: "Updated",
		Enabled:     takeBoolPointer(true),
		ApplicationProfile: types.NsxtAlbVirtualServiceApplicationProfile{
			SystemDefined: true,
			Type:          "HTTPS",
		},
		GatewayRef:            types.OpenApiReference{ID: edge.EdgeGateway.ID},
		LoadBalancerPoolRef:   types.OpenApiReference{ID: pool.NsxtAlbPool.ID},
		ServiceEngineGroupRef: types.OpenApiReference{ID: seGroup.NsxtAlbServiceEngineGroup.ID},
		CertificateRef:        &types.OpenApiReference{ID: createdCertificate.CertificateLibrary.Id},
		ServicePorts: []types.NsxtAlbVirtualServicePort{
			{
				PortStart: takeIntAddress(80),
			},
			{
				PortStart:  takeIntAddress(443),
				SslEnabled: takeBoolPointer(true),
			},
		},
		// Use Primary IP of Edge Gateway as virtual service IP
		VirtualIpAddress: edge.EdgeGateway.EdgeGatewayUplinks[0].Subnets.Values[0].PrimaryIP,
	}

	testAlbVirtualServiceConfig(check, vcd, "WithCertHTTPS", virtualServiceConfig, virtualServiceConfigUpdated, client)

	err = createdCertificate.Delete(ctx)
	check.Assert(err, IsNil)
}

func testAlbVirtualServiceConfig(check *C, vcd *TestVCD, name string, setupConfig *types.NsxtAlbVirtualService, updateConfig *types.NsxtAlbVirtualService, client *VCDClient) {
	fmt.Printf("# Running ALB Virtual Service test with config %s ('System' user: %t) ", name, client.Client.IsSysAdmin)

	edge, err := vcd.nsxtVdc.GetNsxtEdgeGatewayByName(ctx, vcd.config.VCD.Nsxt.EdgeGateway)
	check.Assert(err, IsNil)

	createdVirtualService, err := client.CreateNsxtAlbVirtualService(ctx, setupConfig)
	check.Assert(err, IsNil)

	// Verify mandatory fields
	check.Assert(createdVirtualService.NsxtAlbVirtualService.ID, NotNil)
	check.Assert(createdVirtualService.NsxtAlbVirtualService.Name, NotNil)
	check.Assert(createdVirtualService.NsxtAlbVirtualService.GatewayRef.ID, NotNil)

	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbVirtualServices + createdVirtualService.NsxtAlbVirtualService.ID
	PrependToCleanupListOpenApi(createdVirtualService.NsxtAlbVirtualService.Name, check.TestName(), openApiEndpoint)

	// Get By ID
	virtualServiceById, err := client.GetAlbVirtualServiceById(ctx, createdVirtualService.NsxtAlbVirtualService.ID)
	check.Assert(err, IsNil)
	check.Assert(virtualServiceById.NsxtAlbVirtualService.ID, Equals, createdVirtualService.NsxtAlbVirtualService.ID)

	// Get By Name
	virtualServiceByName, err := client.GetAlbVirtualServiceByName(ctx, edge.EdgeGateway.ID, createdVirtualService.NsxtAlbVirtualService.Name)
	check.Assert(err, IsNil)
	check.Assert(virtualServiceByName.NsxtAlbVirtualService.ID, Equals, createdVirtualService.NsxtAlbVirtualService.ID)

	//Get All Virtual Service summaries
	allVirtualServiceSummaries, err := client.GetAllAlbVirtualServiceSummaries(ctx, edge.EdgeGateway.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(len(allVirtualServiceSummaries) > 0, Equals, true)

	// Get All Pools
	allVirtualServices, err := client.GetAllAlbVirtualServices(ctx, edge.EdgeGateway.ID, nil)
	check.Assert(err, IsNil)
	check.Assert(len(allVirtualServices) > 0, Equals, true)

	check.Assert(len(allVirtualServiceSummaries), Equals, len(allVirtualServices))

	// Attempt an update if config is provided
	if updateConfig != nil {
		updateConfig.ID = createdVirtualService.NsxtAlbVirtualService.ID
		updatedPool, err := createdVirtualService.Update(ctx, updateConfig)
		check.Assert(err, IsNil)
		check.Assert(createdVirtualService.NsxtAlbVirtualService.ID, Equals, updatedPool.NsxtAlbVirtualService.ID)
		check.Assert(updatedPool.NsxtAlbVirtualService.Name, NotNil)
		check.Assert(updatedPool.NsxtAlbVirtualService.GatewayRef.ID, NotNil)
	}

	err = createdVirtualService.Delete(ctx)
	check.Assert(err, IsNil)
	fmt.Printf("Done.\n")
}

func setupAlbVirtualServicePrerequisites(check *C, vcd *TestVCD) (*NsxtAlbController, *NsxtAlbCloud, *NsxtAlbServiceEngineGroup, *NsxtEdgeGateway, *NsxtAlbServiceEngineGroupAssignment, *NsxtAlbPool) {
	controller, cloud, seGroup, edge, assignedSeGroup := setupAlbPoolPrerequisites(check, vcd)

	poolConfig := &types.NsxtAlbPool{
		Name:       check.TestName(),
		Enabled:    takeBoolPointer(true),
		GatewayRef: types.OpenApiReference{ID: edge.EdgeGateway.ID},
	}

	albPool, err := vcd.client.CreateNsxtAlbPool(ctx, poolConfig)
	check.Assert(err, IsNil)

	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbPools + albPool.NsxtAlbPool.ID
	PrependToCleanupListOpenApi(albPool.NsxtAlbPool.Name, check.TestName(), openApiEndpoint)

	return controller, cloud, seGroup, edge, assignedSeGroup, albPool
}

func setupAlbPoolFirewallGroupMembers(check *C, vcd *TestVCD, edge *NsxtEdgeGateway) (*NsxtFirewallGroup, *NsxtAlbPool) {
	// creates ip set
	ipSetConfig := &types.NsxtFirewallGroup{
		Name:        check.TestName(),
		OwnerRef:    &types.OpenApiReference{ID: edge.EdgeGateway.ID},
		Description: "Test IP Set",
		Type:        "IP_SET",
		IpAddresses: []string{"1.1.1.1"},
	}

	ipSet, err := vcd.nsxtVdc.CreateNsxtFirewallGroup(ctx, ipSetConfig)
	check.Assert(err, IsNil)

	// add ip set to cleanup list
	openApiEndpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointFirewallGroups + ipSet.NsxtFirewallGroup.ID
	PrependToCleanupListOpenApi(ipSet.NsxtFirewallGroup.Name, check.TestName(), openApiEndpoint)

	poolConfig := &types.NsxtAlbPool{
		Name:       check.TestName() + "-member-group",
		Enabled:    takeBoolPointer(true),
		GatewayRef: types.OpenApiReference{ID: edge.EdgeGateway.ID},
		MemberGroupRef: &types.OpenApiReference{
			ID: ipSet.NsxtFirewallGroup.ID,
		},
	}

	albPool, err := vcd.client.CreateNsxtAlbPool(ctx, poolConfig)
	check.Assert(err, IsNil)

	openApiEndpoint = types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointAlbPools + albPool.NsxtAlbPool.ID
	PrependToCleanupListOpenApi(albPool.NsxtAlbPool.Name, check.TestName(), openApiEndpoint)

	return ipSet, albPool
}

func tearDownAlbVirtualServicePrerequisites(check *C, albPool *NsxtAlbPool, assignment *NsxtAlbServiceEngineGroupAssignment, edge *NsxtEdgeGateway, seGroup *NsxtAlbServiceEngineGroup, cloud *NsxtAlbCloud, controller *NsxtAlbController) {
	err := albPool.Delete(ctx)
	check.Assert(err, IsNil)
	err = assignment.Delete(ctx)
	check.Assert(err, IsNil)
	err = edge.DisableAlb(ctx)
	check.Assert(err, IsNil)
	err = seGroup.Delete(ctx)
	check.Assert(err, IsNil)
	err = cloud.Delete(ctx)
	check.Assert(err, IsNil)
	err = controller.Delete(ctx)
	check.Assert(err, IsNil)
}
