package okta

import (
	"context"
	"fmt"
)

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

// GetByID fetches a group by ID.
//
// https://developer.okta.com/docs/api/resources/groups#get-group
func (s *GroupsService) GetByID(ctx context.Context, id string) (*Group, *Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, rateLimitGroupsGetUpdateDeleteCategory)
	path := fmt.Sprintf("groups/%s", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	groupOut := new(Group)
	resp, err := s.client.Do(ctx, req, groupOut)
	if err != nil {
		return nil, resp, err
	}

	return groupOut, resp, nil

}

// Add creates a new group.
//
// https://developer.okta.com/docs/api/resources/groups#add-group
func (s *GroupsService) Add(ctx context.Context, profile *GroupProfile) (*Group, *Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, rateLimitGroupsCreateListCategory)
	path := "groups"

	body := map[string]interface{}{"profile": profile}

	req, err := s.client.NewRequest("POST", path, body)
	if err != nil {
		return nil, nil, err
	}

	groupOut := new(Group)
	resp, err := s.client.Do(ctx, req, groupOut)
	if err != nil {
		return nil, resp, err
	}

	return groupOut, resp, nil

}

// UpdateWithProfile modifies a group using a GroupProfile object, it's a wrapper for Update().
//
// Note that delta updates are not supported. You must pass a full GroupProfile object.
func (s *GroupsService) UpdateWithProfile(ctx context.Context, id string, profile *GroupProfile) (*Group, *Response, error) {
	return s.Update(ctx, id, profile)

}

// UpdateWithGroup modifies a group using a Group object, from which the GroupProfile is extracted,
// it's a wrapper for Update().
//
// Note that delta updates are not supported. You must pass a full Group object.
func (s *GroupsService) UpdateWithGroup(ctx context.Context, id string, group *Group) (*Group, *Response, error) {
	profile := group.Profile

	return s.Update(ctx, id, &profile)
}

// Update modifies a group.
//
// Note that delta updates are not supported. You must pass a full GroupProfile object.
//
// https://developer.okta.com/docs/api/resources/groups#update-group
func (s *GroupsService) Update(ctx context.Context, id string, profile *GroupProfile) (*Group, *Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, rateLimitGroupsGetUpdateDeleteCategory)
	path := fmt.Sprintf("groups/%s", id)

	body := map[string]interface{}{"profile": profile}

	req, err := s.client.NewRequest("PUT", path, body)
	if err != nil {
		return nil, nil, err
	}

	groupOut := new(Group)
	resp, err := s.client.Do(ctx, req, groupOut)
	if err != nil {
		return nil, resp, err
	}

	return groupOut, resp, nil

}

// Remove deletes a group.
//
// https://developer.okta.com/docs/api/resources/groups#remove-group
func (s *GroupsService) Remove(ctx context.Context, id string) (*Response, error) {
	ctx = context.WithValue(ctx, rateLimitCategoryCtxKey, rateLimitGroupsGetUpdateDeleteCategory)
	path := fmt.Sprintf("groups/%s", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	groupOut := new(Group)
	resp, err := s.client.Do(ctx, req, groupOut)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
