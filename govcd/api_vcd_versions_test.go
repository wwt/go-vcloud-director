//go:build api || functional || ALL

/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kr/pretty"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_APIVCDMaxVersionIs_Unauthenticated(check *C) {
	config, err := GetConfigStruct()
	check.Assert(err, IsNil)

	vcdClient, err := GetTestVCDFromYaml(config)
	check.Assert(err, IsNil)

	versionCheck := vcdClient.Client.APIVCDMaxVersionIs(context.Background(), ">= 27.0")
	check.Assert(versionCheck, Equals, true)
	check.Assert(vcdClient.Client.supportedVersions.VersionInfos, Not(Equals), 0)
}

func (vcd *TestVCD) Test_APIClientVersionIs_Unauthenticated(check *C) {
	config, err := GetConfigStruct()
	check.Assert(err, IsNil)

	vcdClient, err := GetTestVCDFromYaml(config)
	check.Assert(err, IsNil)

	versionCheck := vcdClient.Client.APIClientVersionIs(">= 27.0")
	check.Assert(versionCheck, Equals, true)
	check.Assert(vcdClient.Client.supportedVersions.VersionInfos, Not(Equals), 0)
}

// Test_APIVCDMaxVersionIs uses already authenticated vcdClient (in SetupSuite)
func (vcd *TestVCD) Test_APIVCDMaxVersionIs(check *C) {
	ctx := context.Background()

	// Minimum supported vCD 8.20 introduced API version 27.0
	versionCheck := vcd.client.Client.APIVCDMaxVersionIs(ctx, ">= 27.0")
	check.Assert(versionCheck, Equals, true)

	mockVcd := getMockVcdWithAPIVersion("27.0")

	var versionTests = []struct {
		version     string
		boolChecker Checker
		isSupported bool
	}{
		{"= 27.0", Equals, true},
		{">= 27.0", Equals, true},
		{">= 25.0, <= 30", Equals, true},
		{"> 27.0", Equals, false},
		{"< 27.0", Equals, false},
		{"invalid", Equals, false},
		{"", Equals, false},
	}

	for _, tt := range versionTests {
		versionCheck := mockVcd.Client.APIVCDMaxVersionIs(ctx, tt.version)
		check.Assert(versionCheck, tt.boolChecker, tt.isSupported)
	}
}

// Test_APIClientVersionIs uses already authenticated vcdClient (in SetupSuite)
func (vcd *TestVCD) Test_APIClientVersionIs(check *C) {

	// Check with currently set version
	versionCheck := vcd.client.Client.APIClientVersionIs(fmt.Sprintf("= %s", vcd.client.Client.APIVersion))
	check.Assert(versionCheck, Equals, true)

	versionCheck = vcd.client.Client.APIClientVersionIs(">= 27.0")
	check.Assert(versionCheck, Equals, true)

	mockVcd := getMockVcdWithAPIVersion("27.0")

	var versionTests = []struct {
		version     string
		boolChecker Checker
		isSupported bool
	}{
		{"= 27.0", Equals, true},
		{">= 27.0", Equals, true},
		{">= 25.0, <= 30", Equals, true},
		{"> 27.0", Equals, false},
		{"< 27.0", Equals, false},
		{"invalid", Equals, false},
		{"", Equals, false},
	}

	for _, tt := range versionTests {
		versionCheck := mockVcd.Client.APIClientVersionIs(tt.version)
		check.Assert(versionCheck, tt.boolChecker, tt.isSupported)
	}
}

func (vcd *TestVCD) Test_validateAPIVersion(check *C) {
	// valid version is checked automatically in SetUpSuite
	// we're checking only for a bad version here
	unsupportedVersion := "999.0"

	config, err := GetConfigStruct()
	check.Assert(err, IsNil)

	vcdClient, err := GetTestVCDFromYaml(config, WithAPIVersion(unsupportedVersion))
	check.Assert(err, IsNil)
	err = vcdClient.Client.validateAPIVersion(context.Background())
	check.Assert(err, ErrorMatches, "API version .* is not supported: version = .* is not supported")
}

func getMockVcdWithAPIVersion(version string) *VCDClient {
	return &VCDClient{
		Client: Client{
			APIVersion: version,
			supportedVersions: SupportedVersions{
				VersionInfos{
					VersionInfo{
						Version: version,
					},
				},
			},
		},
	}
}

func (vcd *TestVCD) Test_GetVcdVersion(check *C) {
	ctx := context.Background()

	version, versionTime, err := vcd.client.Client.GetVcdVersion(ctx)
	check.Assert(err, IsNil)
	check.Assert(version, Not(Equals), "")
	check.Assert(versionTime, Not(Equals), time.Time{})
	reVersion := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+`)
	check.Assert(reVersion.MatchString(version), Equals, true)

	fmt.Printf("VERSION %s\n", version)
	fmt.Printf("DATE    %s\n", versionTime)

	shortVersion, err := vcd.client.Client.GetVcdShortVersion(ctx)
	check.Assert(err, IsNil)
	check.Assert(shortVersion, Not(Equals), "")
	check.Assert(strings.HasPrefix(version, shortVersion), Equals, true)
	if testVerbose {
		fmt.Printf("SHORT VERSION %s\n", shortVersion)
	}

	fullVersion, err := vcd.client.Client.GetVcdFullVersion(ctx)
	check.Assert(err, IsNil)
	digits := fullVersion.Version.Segments()
	check.Assert(len(digits), Not(Equals), 0)
	check.Assert(shortVersion, Equals, fmt.Sprintf("%d.%d.%d", digits[0], digits[1], digits[2]))

	if testVerbose {
		fmt.Printf("FULL VERSION %# v\n", pretty.Formatter(fullVersion))
	}

	// Comparing the current version against itself, without build. Expected result: equal
	result, err := vcd.client.Client.VersionEqualOrGreater(ctx, version, 3)
	check.Assert(err, IsNil)
	check.Assert(result, Equals, true)

	// Comparing the current version against itself, with build. Expected result: equal
	result, err = vcd.client.Client.VersionEqualOrGreater(ctx, version, 4)
	check.Assert(err, IsNil)
	check.Assert(result, Equals, true)

	digits[3] -= 1
	smallerBuild := intListToVersion(digits, 4)
	// Comparing the current version against same version, with smaller build. Expected result: greater
	result, err = vcd.client.Client.VersionEqualOrGreater(ctx, smallerBuild, 4)
	check.Assert(err, IsNil)
	check.Assert(result, Equals, true)

	// Comparing the current version against same version, with bigger build. Expected result: less
	digits[3] += 2
	biggerBuild := intListToVersion(digits, 4)
	result, err = vcd.client.Client.VersionEqualOrGreater(ctx, biggerBuild, 4)
	check.Assert(err, IsNil)
	check.Assert(result, Equals, false)
}
