/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// Certificate is a structure defining a certificate in VCD
// It is called "Certificate Library" in the UI, and "Certificate Library item" in the API
type Certificate struct {
	CertificateLibrary *types.CertificateLibraryItem
	Href               string
	client             *Client
}

// GetCertificateFromLibraryById Returns certificate from library of certificates
func getCertificateFromLibraryById(ctx context.Context, client *Client, id string, additionalHeader map[string]string) (*Certificate, error) {
	endpoint, err := getEndpointByVersion(ctx, client)
	if err != nil {
		return nil, err
	}
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf("empty certificate ID")
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint, id)
	if err != nil {
		return nil, err
	}

	certificate := &Certificate{
		CertificateLibrary: &types.CertificateLibraryItem{},
		client:             client,
		Href:               urlRef.String(),
	}

	err = client.OpenApiGetItem(ctx, minimumApiVersion, urlRef, nil, certificate.CertificateLibrary, additionalHeader)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func getEndpointByVersion(ctx context.Context, client *Client) (string, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointSSLCertificateLibrary
	newerApiVersion, err := client.VersionEqualOrGreater(ctx, "10.3", 3)
	if err != nil {
		return "", err
	}
	if !newerApiVersion {
		// in previous version exist only API with mistype in name
		endpoint = types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointSSLCertificateLibraryOld
	}
	return endpoint, err
}

// GetCertificateFromLibraryById Returns certificate from library of certificates from System Context
func (client *Client) GetCertificateFromLibraryById(ctx context.Context, id string) (*Certificate, error) {
	return getCertificateFromLibraryById(ctx, client, id, nil)
}

// GetCertificateFromLibraryById Returns certificate from library of certificates from Org context
func (adminOrg *AdminOrg) GetCertificateFromLibraryById(ctx context.Context, id string) (*Certificate, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getCertificateFromLibraryById(ctx, adminOrg.client, id, getTenantContextHeader(tenantContext))
}

// addCertificateToLibrary uploads certificates with configuration details
func addCertificateToLibrary(ctx context.Context, client *Client, certificateConfig *types.CertificateLibraryItem,
	additionalHeader map[string]string) (*Certificate, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointSSLCertificateLibrary
	apiVersion, err := client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	typeResponse := &Certificate{
		CertificateLibrary: &types.CertificateLibraryItem{},
		client:             client,
		Href:               urlRef.String(),
	}

	err = client.OpenApiPostItem(ctx, apiVersion, urlRef, nil,
		certificateConfig, typeResponse.CertificateLibrary, additionalHeader)
	if err != nil {
		return nil, err
	}

	return typeResponse, nil
}

// AddCertificateToLibrary uploads certificates with configuration details
func (adminOrg *AdminOrg) AddCertificateToLibrary(ctx context.Context, certificateConfig *types.CertificateLibraryItem) (*Certificate, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return addCertificateToLibrary(ctx, adminOrg.client, certificateConfig, getTenantContextHeader(tenantContext))
}

// AddCertificateToLibrary uploads certificates with configuration details
func (client *Client) AddCertificateToLibrary(ctx context.Context, certificateConfig *types.CertificateLibraryItem) (*Certificate, error) {
	return addCertificateToLibrary(ctx, client, certificateConfig, nil)
}

// getAllCertificateFromLibrary retrieves all certificates. Query parameters can be supplied to perform additional
// filtering
func getAllCertificateFromLibrary(ctx context.Context, client *Client, queryParameters url.Values, additionalHeader map[string]string) ([]*Certificate, error) {
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointSSLCertificateLibrary
	apiVersion, err := client.getOpenApiHighestElevatedVersion(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	responses := []*types.CertificateLibraryItem{{}}
	err = client.OpenApiGetAllItems(ctx, apiVersion, urlRef, queryParameters, &responses, additionalHeader)
	if err != nil {
		return nil, err
	}

	var wrappedCertificates []*Certificate
	for _, response := range responses {
		urlRef, err := client.OpenApiBuildEndpoint(endpoint, response.Id)
		if err != nil {
			return nil, err
		}
		wrappedCertificate := &Certificate{
			CertificateLibrary: response,
			client:             client,
			Href:               urlRef.String(),
		}
		wrappedCertificates = append(wrappedCertificates, wrappedCertificate)
	}

	return wrappedCertificates, nil
}

// GetAllCertificatesFromLibrary retrieves all available certificates from certificate library.
// Query parameters can be supplied to perform additional filtering
func (client *Client) GetAllCertificatesFromLibrary(ctx context.Context, queryParameters url.Values) ([]*Certificate, error) {
	return getAllCertificateFromLibrary(ctx, client, queryParameters, nil)
}

// CountMatchingCertificates searches among all certificates and return the number of certificates
// with the text that matches the given PEM
func (client *Client) CountMatchingCertificates(ctx context.Context, pem string) (int, error) {
	matchingCertificates, err := client.MatchingCertificatesInLibrary(ctx, pem)
	if err != nil {
		return 0, err
	}
	return len(matchingCertificates), nil
}

// MatchingCertificatesInLibrary searches among all certificates and return all certificates
// with the text that matches the given PEM
func (client *Client) MatchingCertificatesInLibrary(ctx context.Context, pem string) ([]*Certificate, error) {
	certificates, err := client.GetAllCertificatesFromLibrary(ctx, nil)
	if err != nil {
		return nil, err
	}
	var matchingCertificates []*Certificate
	for _, cert := range certificates {
		isSame, err := cert.SameAs(pem)
		if err != nil {
			return nil, err
		}
		if isSame {
			matchingCertificates = append(matchingCertificates, cert)
		}
	}
	return matchingCertificates, nil
}

// GetAllCertificatesFromLibrary r retrieves all available certificates from certificate library.
// Query parameters can be supplied to perform additional filtering
func (adminOrg *AdminOrg) GetAllCertificatesFromLibrary(ctx context.Context, queryParameters url.Values) ([]*Certificate, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getAllCertificateFromLibrary(ctx, adminOrg.client, queryParameters, getTenantContextHeader(tenantContext))
}

// getCertificateFromLibraryByName retrieves certificate from certificate library by given name
// When the alias contains commas, semicolons or asterisks, the encoding is rejected by the API in VCD.
// For this reason, when one or more commas, semicolons or asterisks are present we run the search brute force,
// by fetching all certificates and comparing the alias.
// Also, url.QueryEscape as well as url.Values.Encode() both encode the space as a + character. So we use
// search brute force too. Reference to issue:
// https://github.com/golang/go/issues/4013
// https://github.com/czos/goamz/pull/11/files
func getCertificateFromLibraryByName(ctx context.Context, client *Client, name string, additionalHeader map[string]string) (*Certificate, error) {
	slowSearch, params := shouldDoSlowSearch("alias", name)

	var foundCertificates []*Certificate
	certificates, err := getAllCertificateFromLibrary(ctx, client, params, additionalHeader)
	if err != nil {
		return nil, err
	}
	if len(certificates) == 0 {
		return nil, ErrorEntityNotFound
	}
	foundCertificates = append(foundCertificates, certificates[0])

	if slowSearch {
		foundCertificates = nil
		for _, certificate := range certificates {
			if certificate.CertificateLibrary.Alias == name {
				foundCertificates = append(foundCertificates, certificate)
			}
		}
		if len(foundCertificates) == 0 {
			return nil, ErrorEntityNotFound
		}
		if len(foundCertificates) > 1 {
			return nil, fmt.Errorf("more than one certificate found with name '%s'", name)
		}
	}

	if len(certificates) > 1 && !slowSearch {
		{
			return nil, fmt.Errorf("more than one certificate found with name '%s'", name)
		}
	}
	return foundCertificates[0], nil
}

// GetCertificateFromLibraryByName retrieves certificate from certificate library by given name
func (client *Client) GetCertificateFromLibraryByName(ctx context.Context, name string) (*Certificate, error) {
	return getCertificateFromLibraryByName(ctx, client, name, nil)
}

// GetCertificateFromLibraryByName retrieves certificate from certificate library by given name
func (adminOrg *AdminOrg) GetCertificateFromLibraryByName(ctx context.Context, name string) (*Certificate, error) {
	tenantContext, err := adminOrg.getTenantContext()
	if err != nil {
		return nil, err
	}
	return getCertificateFromLibraryByName(ctx, adminOrg.client, name, getTenantContextHeader(tenantContext))
}

// Update updates existing Certificate. Allows changing only alias and description
func (certificate *Certificate) Update(ctx context.Context) (*Certificate, error) {
	endpoint, err := getEndpointByVersion(ctx, certificate.client)
	if err != nil {
		return nil, err
	}
	minimumApiVersion, err := certificate.client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if certificate.CertificateLibrary.Id == "" {
		return nil, fmt.Errorf("cannot update certificate without id")
	}

	urlRef, err := certificate.client.OpenApiBuildEndpoint(endpoint, certificate.CertificateLibrary.Id)
	if err != nil {
		return nil, err
	}

	returnCertificate := &Certificate{
		CertificateLibrary: &types.CertificateLibraryItem{},
		client:             certificate.client,
	}

	err = certificate.client.OpenApiPutItem(ctx, minimumApiVersion, urlRef, nil, certificate.CertificateLibrary,
		returnCertificate.CertificateLibrary, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating certificate: %s", err)
	}

	return returnCertificate, nil
}

// Delete deletes certificate from Certificate library
func (certificate *Certificate) Delete(ctx context.Context) error {
	endpoint, err := getEndpointByVersion(ctx, certificate.client)
	if err != nil {
		return err
	}
	minimumApiVersion, err := certificate.client.checkOpenApiEndpointCompatibility(ctx, endpoint)
	if err != nil {
		return err
	}

	if certificate.CertificateLibrary.Id == "" {
		return fmt.Errorf("cannot delete certificate without id")
	}

	urlRef, err := certificate.client.OpenApiBuildEndpoint(endpoint, certificate.CertificateLibrary.Id)
	if err != nil {
		return err
	}

	err = certificate.client.OpenApiDeleteItem(ctx, minimumApiVersion, urlRef, nil, nil)

	if err != nil {
		return fmt.Errorf("error deleting certificate: %s", err)
	}

	return nil
}

// getCertificateText returns the stripped text of the certificate, without the
// starting and ending markers
func getCertificateText(pem string) (string, error) {
	reText, err := regexp.Compile(
		`(?s)` + // treats newline as any other character
			`-----BEGIN CERTIFICATE-----` + // the 'begin certificate' marker
			`(.+)` + // any sequence of characters after the 'begin certificate' marker
			`-----END CERTIFICATE-----`) // the 'end certificate' marker
	if err != nil {
		return "", err
	}

	text := reText.FindStringSubmatch(pem)
	if len(text) < 2 {
		return "", fmt.Errorf("start marker 'BEGIN CERTIFICATE' or end marker 'END CERTIFICATE' not found in input certificate")
	}
	return text[1], nil
}

// SameAs returns true if the certificate text matches the text of the provided PEM
// (without the BEGIN CERTIFICATE and END CERTIFICATE markers)
func (certificate *Certificate) SameAs(pem string) (bool, error) {
	internalValue, err := getCertificateText(certificate.CertificateLibrary.Certificate)
	if err != nil {
		return false, err
	}
	compareValue, err := getCertificateText(pem)
	if err != nil {
		return false, err
	}
	return internalValue == compareValue, nil
}
