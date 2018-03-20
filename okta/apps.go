package okta

import (
	"context"
	"fmt"
)

// AppsService is the service providing access to the App Resource in the Okta API
type AppsService service

// GetByID fetches a single application by its ID
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
