/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	semver "github.com/hashicorp/go-version"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"github.com/vmware/go-vcloud-director/v2/util"
)

// VCDClientOption defines signature for customizing VCDClient using
// functional options pattern.
type VCDClientOption func(*VCDClient) error

type VCDClient struct {
	Client      Client  // Client for the underlying VCD instance
	sessionHREF url.URL // HREF for the session API
	QueryHREF   url.URL // HREF for the query API
}

func (vcdClient *VCDClient) vcdloginurl(ctx context.Context) error {
	if err := vcdClient.Client.validateAPIVersion(ctx); err != nil {
		return fmt.Errorf("could not find valid version for login: %s", err)
	}

	// find login address matching the API version
	var neededVersion VersionInfo
	for _, versionInfo := range vcdClient.Client.supportedVersions.VersionInfos {
		if versionInfo.Version == vcdClient.Client.APIVersion {
			neededVersion = versionInfo
			break
		}
	}

	loginUrl, err := url.Parse(neededVersion.LoginUrl)
	if err != nil {
		return fmt.Errorf("couldn't find a LoginUrl for version %s", vcdClient.Client.APIVersion)
	}
	vcdClient.sessionHREF = *loginUrl
	return nil
}

// vcdCloudApiAuthorize performs the authorization to VCD using open API
func (vcdClient *VCDClient) vcdCloudApiAuthorize(ctx context.Context, user, pass, org string) (*http.Response, error) {
	var missingItems []string
	if user == "" {
		missingItems = append(missingItems, "user")
	}
	if pass == "" {
		missingItems = append(missingItems, "password")
	}
	if org == "" {
		missingItems = append(missingItems, "org")
	}
	if len(missingItems) > 0 {
		return nil, fmt.Errorf("authorization is not possible because of these missing items: %v", missingItems)
	}

	util.Logger.Println("[TRACE] Connecting to VCD using cloudapi")
	// This call can only be used by tenants
	rawUrl := vcdClient.sessionHREF.Scheme + "://" + vcdClient.sessionHREF.Host + "/cloudapi/1.0.0/sessions"

	// If we are connecting as provider, we need to qualify the request.
	if strings.EqualFold(org, "system") {
		rawUrl += "/provider"
	}
	util.Logger.Printf("[TRACE] URL %s\n", rawUrl)
	loginUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL %s", rawUrl)
	}
	vcdClient.sessionHREF = *loginUrl
	req := vcdClient.Client.NewRequest(ctx, map[string]string{}, http.MethodPost, *loginUrl, nil)
	// Set Basic Authentication Header
	req.SetBasicAuth(user+"@"+org, pass)
	// Add the Accept header. The version must be at least 33.0 for cloudapi to work
	req.Header.Add("Accept", "application/*;version="+vcdClient.Client.APIVersion)
	resp, err := vcdClient.Client.Http.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			util.Logger.Printf("error closing response Body [vcdCloudApiAuthorize]: %s", err)
		}
	}(resp.Body)

	// Catch HTTP 401 (Status Unauthorized) to return an error as otherwise this library would return
	// odd errors while doing lookup of resources and confuse user.
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("received response HTTP %d (Unauthorized). Please check if your credentials are valid",
			resp.StatusCode)
	}

	// Store the authorization header
	vcdClient.Client.VCDToken = resp.Header.Get(BearerTokenHeader)
	vcdClient.Client.VCDAuthHeader = BearerTokenHeader
	vcdClient.Client.IsSysAdmin = strings.EqualFold(org, "system")
	// Get query href
	vcdClient.QueryHREF = vcdClient.Client.VCDHREF
	vcdClient.QueryHREF.Path += "/query"
	return resp, nil
}

// NewVCDClient initializes VMware VMware Cloud Director client with reasonable defaults.
// It accepts functions of type VCDClientOption for adjusting defaults.
func NewVCDClient(vcdEndpoint url.URL, insecure bool, options ...VCDClientOption) *VCDClient {
	minVcdApiVersion := "36.0" // supported by 10.3+
	userDefinedApiVersion := os.Getenv("GOVCD_API_VERSION")
	if userDefinedApiVersion != "" {
		_, err := semver.NewVersion(userDefinedApiVersion)
		if err != nil {
			// We do not have error in return of this function signature.
			// To avoid breaking API the only thing we can do is panic.
			panic(fmt.Sprintf("unable to initialize VCD client from environment variable GOVCD_API_VERSION. Version '%s' is not valid: %s", userDefinedApiVersion, err))
		}
		minVcdApiVersion = userDefinedApiVersion
	}

	// Setting defaults
	// #nosec G402 -- InsecureSkipVerify: insecure - This allows connecting to VCDs with self-signed certificates
	vcdClient := &VCDClient{
		Client: Client{
			APIVersion: minVcdApiVersion,
			// UserAgent cannot embed exact version by default because this is source code and is supposed to be used by programs,
			// but any client can customize or disable it at all using WithHttpUserAgent() configuration options function.
			UserAgent: "go-vcloud-director",
			VCDHREF:   vcdEndpoint,
			Http: http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: insecure,
					},
					Proxy:               http.ProxyFromEnvironment,
					TLSHandshakeTimeout: 120 * time.Second, // Default timeout for TSL hand shake
				},
				Timeout: 600 * time.Second, // Default value for http request+response timeout
			},
			MaxRetryTimeout: 60, // Default timeout in seconds for retries calls in functions
		},
	}

	// Override defaults with functional options
	for _, option := range options {
		err := option(vcdClient)
		if err != nil {
			// We do not have error in return of this function signature.
			// To avoid breaking API the only thing we can do is panic.
			panic(fmt.Sprintf("unable to initialize VCD client: %s", err))
		}
	}
	return vcdClient
}

// Authenticate is a helper function that performs a login in VMware Cloud Director.
func (vcdClient *VCDClient) Authenticate(ctx context.Context, username, password, org string) error {
	_, err := vcdClient.GetAuthResponse(ctx, username, password, org)
	return err
}

// GetAuthResponse performs authentication and returns the full HTTP response
// The purpose of this function is to preserve information that is useful
// for token-based authentication
func (vcdClient *VCDClient) GetAuthResponse(ctx context.Context, username, password, org string) (*http.Response, error) {
	// LoginUrl
	err := vcdClient.vcdloginurl(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding LoginUrl: %s", err)
	}

	// Choose correct auth mechanism based on what type of authentication is used. The end result
	// for each of the below functions is to set authorization token vcdCli.Client.VCDToken.
	var resp *http.Response
	switch {
	case vcdClient.Client.UseSamlAdfs:
		err = vcdClient.authorizeSamlAdfs(ctx, username, password, org, vcdClient.Client.CustomAdfsRptId)
		if err != nil {
			return nil, fmt.Errorf("error authorizing SAML: %s", err)
		}
	default:
		// Authorize
		resp, err = vcdClient.vcdCloudApiAuthorize(ctx, username, password, org)
		if err != nil {
			return nil, fmt.Errorf("error authorizing: %s", err)
		}
	}

	vcdClient.LogSessionInfo(ctx)
	return resp, nil
}

// SetToken will set the authorization token in the client, without using other credentials
// Up to version 29, token authorization uses the header key x-vcloud-authorization
// In version 30+ it also uses X-Vmware-Vcloud-Access-Token:TOKEN coupled with
// X-Vmware-Vcloud-Token-Type:"bearer"
func (vcdClient *VCDClient) SetToken(ctx context.Context, org, authHeader, token string) error {
	if authHeader == ApiTokenHeader {
		util.Logger.Printf("[DEBUG] Attempt authentication using API token")
		apiToken, err := vcdClient.GetBearerTokenFromApiToken(ctx, org, token)
		if err != nil {
			util.Logger.Printf("[DEBUG] Authentication using API token was UNSUCCESSFUL: %s", err)
			return err
		}
		token = apiToken.AccessToken
		authHeader = BearerTokenHeader
		vcdClient.Client.UsingAccessToken = true
		util.Logger.Printf("[DEBUG] Authentication using API token was SUCCESSFUL")
	}
	if !vcdClient.Client.UsingAccessToken {
		vcdClient.Client.UsingBearerToken = true
	}
	vcdClient.Client.VCDAuthHeader = authHeader
	vcdClient.Client.VCDToken = token

	err := vcdClient.vcdloginurl(ctx)
	if err != nil {
		return fmt.Errorf("error finding LoginUrl: %s", err)
	}

	vcdClient.Client.IsSysAdmin = strings.EqualFold(org, "system")
	// Get query href
	vcdClient.QueryHREF = vcdClient.Client.VCDHREF
	vcdClient.QueryHREF.Path += "/query"

	// The client is now ready to connect using the token, but has not communicated with the vCD yet.
	// To make sure that it is working, we run a request for the org list.
	// This list should work always: when run as system administrator, it retrieves all organizations.
	// When run as org user, it only returns the organization the user is authorized to.
	// In both cases, we discard the list, as we only use it to certify that the token works.
	orgListHREF := vcdClient.Client.VCDHREF
	orgListHREF.Path += "/org"

	orgList := new(types.OrgList)

	_, err = vcdClient.Client.ExecuteRequest(ctx, orgListHREF.String(), http.MethodGet,
		"", "error connecting to vCD using token: %s", nil, orgList)
	if err != nil {
		return err
	}
	vcdClient.LogSessionInfo(ctx)
	return nil
}

// Disconnect performs a disconnection from the VMware Cloud Director API endpoint.
func (vcdClient *VCDClient) Disconnect(ctx context.Context) error {
	if vcdClient.Client.VCDToken == "" && vcdClient.Client.VCDAuthHeader == "" {
		return fmt.Errorf("cannot disconnect, client is not authenticated")
	}
	req := vcdClient.Client.NewRequest(ctx, map[string]string{}, http.MethodDelete, vcdClient.sessionHREF, nil)
	// Add the Accept header for vCA
	req.Header.Add("Accept", "application/xml;version="+vcdClient.Client.APIVersion)
	// Set Authorization Header
	req.Header.Add(vcdClient.Client.VCDAuthHeader, vcdClient.Client.VCDToken)
	if _, err := checkResp(vcdClient.Client.Http.Do(req)); err != nil {
		return fmt.Errorf("error processing session delete for VMware Cloud Director: %s", err)
	}
	return nil
}

// WithMaxRetryTimeout allows default vCDClient MaxRetryTimeout value override
func WithMaxRetryTimeout(timeoutSeconds int) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.MaxRetryTimeout = timeoutSeconds
		return nil
	}
}

// WithAPIVersion allows to override default API version. Please be cautious
// about changing the version as the default specified is the most tested.
func WithAPIVersion(version string) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.APIVersion = version
		return nil
	}
}

// WithHttpTimeout allows to override default http timeout
func WithHttpTimeout(timeout int64) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.Http.Timeout = time.Duration(timeout) * time.Second
		return nil
	}
}

// WithSamlAdfs specifies if SAML auth is used for authenticating to vCD instead of local login.
// The following conditions must be met so that SAML authentication works:
// * SAML IdP (Identity Provider) is Active Directory Federation Service (ADFS)
// * WS-Trust authentication endpoint "/adfs/services/trust/13/usernamemixed" must be enabled on
// ADFS server
// By default vCD SAML Entity ID will be used as Relaying Party Trust Identifier unless
// customAdfsRptId is specified
func WithSamlAdfs(useSaml bool, customAdfsRptId string) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.UseSamlAdfs = useSaml
		vcdClient.Client.CustomAdfsRptId = customAdfsRptId
		return nil
	}
}

// WithHttpUserAgent allows to specify HTTP user-agent which can be useful for statistics tracking.
// By default User-Agent is set to "go-vcloud-director". It can be unset by supplying empty value.
func WithHttpUserAgent(userAgent string) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.UserAgent = userAgent
		return nil
	}
}

// WithHttpHeader allows to specify custom HTTP header values.
// Typical usage of this function is to inject a tenant context into the client.
//
// WARNING: Using this function in an environment with concurrent operations may result in negative side effects,
// such as operations as system administrator and as tenant using the same client.
// This setting is justified when we want to start a session where the additional header is always needed.
// For cases where we need system administrator and tenant operations in the same environment we can either
// a) use two separate clients
// or b) use the `additionalHeader` parameter in *newRequest* functions
func WithHttpHeader(options map[string]string) VCDClientOption {
	return func(vcdClient *VCDClient) error {
		vcdClient.Client.customHeader = make(http.Header)
		for k, v := range options {
			vcdClient.Client.customHeader.Add(k, v)
		}
		return nil
	}
}
