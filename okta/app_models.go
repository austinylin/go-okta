package okta

// App represents an application in Okta
type App struct {
	ID            string           `json:"id,omitempty"`
	Name          string           `json:"name,omitempty"`
	Label         string           `json:"label,omitempty"`
	Created       Timestamp        `json:"created,omitempty"`
	LastUpdated   Timestamp        `json:"lastUpdated,omitempty"`
	Status        string           `json:"status,omitempty"`
	Features      []string         `json:"features,omitempty"`
	SignOnMode    AppSignOnMode    `json:"signOnMode,omitempty"`
	Accessibility AppAccessibility `json:"accessibility,omitempty"`
	// Visibility not implemented
	Credentials AppCredential `json:"credentials,omitempty"`
	Settings    interface{}   `json:"settings,omitempty"`
	Profile     interface{}   `json:"profile,omitempty"`
}

// AppAccessibility determines accessibility settings for the application.
//
// https://developer.okta.com/docs/api/resources/apps#accessibility-object
type AppAccessibility struct {
	SelfService      bool   `json:"selfService,omitempty"`
	ErrorRedirectURL string `json:"errorRedirectUrl,omitempty"`
	LoginRedirectURL string `json:"loginRedirectUrl,omitempty"`
}

// AppSignOnMode is a type for the SignOnMode enum
//
// https://developer.okta.com/docs/api/resources/apps#signon-modes
type AppSignOnMode string

// AppSignOnMode Constants
const (
	Bookmark            AppSignOnMode = "BOOKMARK"
	BasicAuth                         = "BASIC_AUTH"
	BrowserPlugin                     = "BROWSER_PLUGIN"
	SecurePasswordStore               = "SECURE_PASSWORD_STORE"
	SAML2                             = "SAML_2_0"
	WSFederation                      = "WS_FEDERATION"
	AutoLogin                         = "AUTO_LOGIN"
	OpenIDConnect                     = "OPENID_CONNECT"
	AppSignOnModeCustom               = "Custom"
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
	Template string
	// Type has possible values of: "NONE", "BUILT_IN", "CUSTOM"
	Type       string `json:"type,omitempty"`
	UserSuffix string
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
	Value string `json:"value"`
}
