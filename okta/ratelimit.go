//go:generate stringer -type=rateLimitCategory

package okta

var rateLimitCategoryCtxKey contextKey

// Rate represents an the status of an individual rate limit.
type Rate struct {
	Limit     int
	Remaining int
	Reset     Timestamp
}

type rateLimitCategory int

const (
	rateLimitCoreCategory rateLimitCategory = iota
	rateLimitAppsCreateListCategory
	rateLimitAppsGetUpdateDeleteCategory
	rateLimitAuthnCategory
	rateLimitGroupsCreateListCategory
	rateLimitGroupsGetUpdateDeleteCategory
	rateLimitLogsCategory
	rateLimitSessionsCategory
	rateLimitUsersCreateListCategory
	rateLimitUsersGetByIDCategory
	rateLimitUsersGetByLoginNameCategory
	rateLimitUsersCreateUpdateDeleteByIDCategory

	categories // An array of this length will be able to contain all rate limit categories.
)
