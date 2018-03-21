package okta

import "net/url"

// App represents an application in Okta
type App struct {
	ID            string           `json:"id,omitempty"`
	Name          AppName          `json:"name,omitempty"`
	Label         string           `json:"label,omitempty"`
	Created       Timestamp        `json:"created,omitempty"`
	LastUpdated   Timestamp        `json:"lastUpdated,omitempty"`
	Status        string           `json:"status,omitempty"`
	Features      []string         `json:"features,omitempty"`
	SignOnMode    AppSignOnMode    `json:"signOnMode"`
	Accessibility AppAccessibility `json:"accessibility"`
	Visibility    AppVisability    `json:"visibility"`
	Credentials   AppCredential    `json:"credentials"`
	Settings      interface{}      `json:"settings,omitempty"`
	Profile       interface{}      `json:"profile,omitempty"`
}

// AppName is a type for the AppName enum.
// Note that name in the okta context is used to delinate the type of app.
// Shared apps, which can be used by multiple Okta Customers, aren't implemented.
//
// https://developer.okta.com/docs/api/resources/apps#app-names--settings
type AppName string

// AppName Constants
// Note that name in the okta context is used to delinate the type of app.
// Shared apps, which can be used by multiple Okta Customers, aren't implemented.
//
// https://developer.okta.com/docs/api/resources/apps#app-names--settings
const (
	AppNameBookmark AppName = "bookmark"
	AppNameSAML2            = "Custom SAML 2.0"
	// AppNameOAuth2           = "oidc_client"
	// AppNameSWA              = "Custom SWA"
)

// AppAccessibility determines accessibility settings for the application.
//
// https://developer.okta.com/docs/api/resources/apps#accessibility-object
type AppAccessibility struct {
	SelfService      bool   `json:"selfService"`
	ErrorRedirectURL string `json:"errorRedirectUrl"`
	LoginRedirectURL string `json:"loginRedirectUrl"`
}

// AppSignOnMode is a type for the SignOnMode enum
//
// https://developer.okta.com/docs/api/resources/apps#signon-modes
type AppSignOnMode string

// AppSignOnMode Constants
//
// https://developer.okta.com/docs/api/resources/apps#signon-modes
const (
	AppSignOnModeBookmark            AppSignOnMode = "BOOKMARK"
	AppSignOnModeBasicAuth                         = "BASIC_AUTH"
	AppSignOnModeBrowserPlugin                     = "BROWSER_PLUGIN"
	AppSignOnModeSecurePasswordStore               = "SECURE_PASSWORD_STORE"
	AppSignOnModeSAML2                             = "SAML_2_0"
	AppSignOnModeWSFederation                      = "WS_FEDERATION"
	AppSignOnModeAutoLogin                         = "AUTO_LOGIN"
	AppSignOnModeOpenIDConnect                     = "OPENID_CONNECT"
	AppSignOnModeCustom                            = "Custom"
)

// AppAuthenticationScheme is the type for the AppAuthenticationScheme enum
//
// https://developer.okta.com/docs/api/resources/apps#authentication-schemes
type AppAuthenticationScheme string

// AppAuthenticationScheme Constants
//
// https://developer.okta.com/docs/api/resources/apps#authentication-schemes
const (
	SharedUsernameAndPassword AppAuthenticationScheme = "SHARED_USERNAME_AND_PASSWORD"
	ExternalPasswordSync                              = "EXTERNAL_PASSWORD_SYNC"
	EditUsernameAndPassword                           = "EDIT_USERNAME_AND_PASSWORD"
	EditPasswordOnly                                  = "EDIT_PASSWORD_ONLY"
	AdminSetsCredentials                              = "ADMIN_SETS_CREDENTIALS"
)

// AppCredential specifies credentials and scheme for the application’s signOnMode
//
// https://developer.okta.com/docs/api/resources/apps#application-credentials-object
type AppCredential struct {
	Scheme           AppAuthenticationScheme        `json:"scheme,omitempty"`
	UserNameTemplate AppCredentialsUserNameTemplate `json:"userNameTemplate,omitempty"`
	Signing          AppCredentialSigningCredential `json:"signing,omitempty"`
	UserName         string                         `json:"username,omitempty"`
	Password         AppPassword                    `json:"password,omitempty"`
	OAuthClient      AppCredentialOAuthCredential   `json:"oauthClient,omitempty"`
}

// AppCredentialsUserNameTemplate represents the template used to generate the username when an
// app is assigend to a user.
//
// https://developer.okta.com/docs/api/resources/apps#username-template-object
type AppCredentialsUserNameTemplate struct {
	Template string `json:"template,omitempty"`
	// Type has possible values of: "NONE", "BUILT_IN", "CUSTOM"
	Type       string `json:"type,omitempty"`
	UserSuffix string `json:"userSuffix,omitempty"`
}

// AppCredentialSigningCredential determines the key used for signing assertions for the signOnMode.
//
// https://developer.okta.com/docs/api/resources/apps#signing-credential-object
type AppCredentialSigningCredential struct {
	KID string `json:"kid,omitempty"`
}

// AppCredentialOAuthCredential determines how to authenticate the OAuth 2.0 client.
//
// https://developer.okta.com/docs/api/resources/apps#oauth-credential-object
type AppCredentialOAuthCredential struct {
	ClientID                string `json:"client_id,omitempty"`
	ClientSecret            string `json:"client_secret,omitempty"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`
	AutoKeyRotation         bool   `json:"autoKeyRotation,omitempty"`
}

// AppPassword represents a password for user:app combination.
//
// It has one attribute, value which is write only.
//
// https://developer.okta.com/docs/api/resources/apps#password-object
type AppPassword struct {
	// Value is a write only property. An empty object represents a password exists.
	Value string `json:"value,omitempty"`
}

// AppSAMLAttributeStatement represents Attribute Statements for SAML apps.
//
// https://developer.okta.com/docs/api/resources/apps#attribute-statements-object
type AppSAMLAttributeStatement struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Values    []string `json:"values"`
}

// AppAddSAMLAppParams is a helper struct for calling AddSAMLApp().
type AppAddSAMLAppParams struct {
	DefaultRelayState     string
	SsoAcsURL             *url.URL
	Recipient             *url.URL
	Destination           *url.URL
	Audience              string
	IdpIssuer             string
	SubjectNameIDTemplate string
	SubjectNameIDFormat   string
	ResponseSigned        bool
	AssertionSigned       bool
	SignatureAlgorithm    string
	DigestAlgorithm       string
	HonorForceAuthn       bool
	AuthnContextClassRef  string
	AttributeStatements   []AppSAMLAttributeStatement
}

// AppVisability represents where an app is shown.
//
// https://developer.okta.com/docs/api/resources/apps#visibility-object
type AppVisability struct {
	AutoSubmitToolbar bool              `json:"autoSubmitToolbar"`
	Hide              AppVisabilityHide `json:"hide"`
	// AppLinks
}

// AppVisabilityHide is a helper struct.
//
// https://developer.okta.com/docs/api/resources/apps#hide-object
type AppVisabilityHide struct {
	IOS bool `json:"iOS"`
	Web bool `json:"web"`
}
