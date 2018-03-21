package okta

// NewAppVisability is a helper method to create a new AppVisability object with default settings.
func NewAppVisability() AppVisability {
	return AppVisability{
		AutoSubmitToolbar: false,
		Hide: AppVisabilityHide{
			IOS: false,
			Web: false,
		},
	}
}
