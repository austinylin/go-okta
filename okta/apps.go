package okta

import (
	"context"
	"fmt"
	"net/url"
)

// AppsService is the service providing access to the App Resource in the Okta API
type AppsService service

// GetByID fetches a single application by its ID
//
// https://developer.okta.com/docs/api/resources/apps#get-application
func (s *AppsService) GetByID(ctx context.Context, id string) (*App, *Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, appsGetUpdateDeleteCategory)
	path := fmt.Sprintf("apps/%s", id)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	app := new(App)
	resp, err := s.client.Do(ctx, req, app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// AddBookmarkApp creates a new bookmark application, it wraps Add().
//
//	https://developer.okta.com/docs/api/resources/apps#add-bookmark-application
func (s *AppsService) AddBookmarkApp(ctx context.Context, label string, activate bool, url *url.URL) (*App, *Response, error) {
	appIn := new(App)
	appIn.SignOnMode = AppSignOnModeBookmark
	appIn.Name = AppNameBookmark
	appIn.Label = label
	appIn.Settings = map[string]map[string]interface{}{
		"app": {
			"requestIntegration": false,
			"url":                url.String(),
		},
	}

	appOut, resp, err := s.Add(ctx, appIn, activate)
	return appOut, resp, err
}

// AddSAMLApp creates a new SAML application, it wraps Add(). Caveats:
// 	- Okta Docs: Fields that require certificate uploads can’t be enabled through the API, such as Single Log Out and Assertion Encryption. These must be updated through the UI.
//  - Implementation Limitation: Override attributes aren't supported.
//
//	https://developer.okta.com/docs/api/resources/apps#add-custom-saml-application
func (s *AppsService) AddSAMLApp(
	ctx context.Context,
	label string,
	activate bool,
	params *AppAddSAMLAppParams,
) (*App, *Response, error) {

	// Okta Docs: Either (or both) “responseSigned” or “assertionSigned” must be TRUE.
	if !params.ResponseSigned && !params.AssertionSigned {
		return nil, nil, fmt.Errorf("Invalid paramaters, either `ResponseSigned` or `AssertionSigned` must be true")
	}

	// Defaults
	switch {
	case params.SignatureAlgorithm == "":
		params.SignatureAlgorithm = "RSA_SHA256"
		fallthrough
	case params.DigestAlgorithm == "":
		params.DigestAlgorithm = "SHA256"
		fallthrough
	case params.SubjectNameIDTemplate == "":
		params.SubjectNameIDTemplate = "${user.userName}"
		fallthrough
	case params.SubjectNameIDFormat == "":
		params.SubjectNameIDFormat = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
		fallthrough
	case params.AuthnContextClassRef == "":
		params.AuthnContextClassRef = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
		fallthrough
	case params.IdpIssuer == "":
		params.IdpIssuer = "http://www.okta.com/${org.externalKey}"
		fallthrough
	default:
	}

	// Default Namespace for Attribute Statements
	if len(params.AttributeStatements) > 0 {
		for _, elem := range params.AttributeStatements {
			if elem.Namespace == "" {
				elem.Namespace = "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"
			}
		}
	} else {
		params.AttributeStatements = make([]AppSAMLAttributeStatement, 0)
	}

	appIn := new(App)
	appIn.SignOnMode = AppSignOnModeSAML2
	appIn.Name = "" // Omited for custom SAML apps
	appIn.Label = label
	appIn.Visibility = NewAppVisability()
	appIn.Settings = map[string]map[string]interface{}{
		"signOn": {
			"defaultRelayState":     params.DefaultRelayState,
			"ssoAcsUrl":             params.SsoAcsURL.String(),
			"recipient":             params.Recipient.String(),
			"destination":           params.Destination.String(),
			"audience":              params.Audience,
			"idpIssuer":             params.IdpIssuer,
			"subjectNameIdTemplate": params.SubjectNameIDTemplate,
			"subjectNameIdFormat":   params.SubjectNameIDFormat,
			"responseSigned":        params.ResponseSigned,
			"assertionSigned":       params.AssertionSigned,
			"signatureAlgorithm":    params.SignatureAlgorithm,
			"digestAlgorithm":       params.DigestAlgorithm,
			"honorForceAuthn":       params.HonorForceAuthn,
			"authnContextClassRef":  params.AuthnContextClassRef,
			"attributeStatements":   params.AttributeStatements,
		},
	}

	appOut, resp, err := s.Add(ctx, appIn, activate)
	return appOut, resp, err
}

// Add creates a new application. Most people will want to call one of the helper methods instead.
//
// https://developer.okta.com/docs/api/resources/apps#add-application
func (s *AppsService) Add(ctx context.Context, appIn *App, activate bool) (*App, *Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, appsCreateListCategory)
	path := fmt.Sprintf("apps?activate=%t", activate)
	req, err := s.client.NewRequest("POST", path, appIn)
	if err != nil {
		return nil, nil, err
	}

	appOut := new(App)
	resp, err := s.client.Do(ctx, req, appOut)
	if err != nil {
		return nil, resp, err
	}

	return appOut, resp, nil
}
