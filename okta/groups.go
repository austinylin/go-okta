package okta

// GroupsService is the service providing access to the Groups Resource in the Okta API
type GroupsService service

// Group represents an Okta Group.
//
// https://developer.okta.com/docs/api/resources/groups#group-model
type Group struct {
	ID                    string       `json:"id,omitempty"`
	Created               Timestamp    `json:"created,omitempty"`
	LastUpdated           Timestamp    `json:"lastUpdated,omitempty"`
	LastMembershipUpdated Timestamp    `json:"lastMembershipUpdated,omitempty"`
	ObjectClass           []string     `json:"objectClass,omitempty"`
	Type                  string       `json:"type,omitempty"`
	Profile               GroupProfile `json:"profile"`
}

// GroupProfile represents an Okta Group Profile.
//
// https://developer.okta.com/docs/api/resources/groups#profile-object
type GroupProfile struct {
	Name                       string `json:"name,omitempty"`
	Description                string `json:"description,omitempty"`
	SamAccountName             string `json:"samAccountName,omitempty"`
	DN                         string `json:"dn,omitempty"`
	WindowsDomainQualifiedName string `json:"windowsDomainQualifiedName,omitempty"`
	ExternalID                 string `json:"externalId,omitempty"`
}
