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
	coreCategory rateLimitCategory = iota
	appsCreateListCategory
	appsGetUpdateDeleteCategory
	authnCategory
	groupsCreateListCategory
	groupsGetUpdateDeleteCategory
	logsCategory
	sessionsCategory
	usersCreateListCategory
	usersGetByIDCategory
	usersGetByLoginNameCategory
	usersCreateUpdateDeleteByIDCategory

	categories // An array of this length will be able to contain all rate limit categories.
)
