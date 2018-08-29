package okta

import "time"

/*
{
  "id": "00ub0oNGTSWTBKOLGLNR",
  "status": "ACTIVE",
  "created": "2013-06-24T16:39:18.000Z",
  "activated": "2013-06-24T16:39:19.000Z",
  "statusChanged": "2013-06-24T16:39:19.000Z",
  "lastLogin": "2013-06-24T17:39:19.000Z",
  "lastUpdated": "2013-06-27T16:35:28.000Z",
  "passwordChanged": "2013-06-24T16:39:19.000Z",
  "profile": {
    "login": "isaac.brock@example.com",
    "firstName": "Isaac",
    "lastName": "Brock",
    "nickName": "issac",
    "displayName": "Isaac Brock",
    "email": "isaac.brock@example.com",
    "secondEmail": "isaac@example.org",
    "profileUrl": "http://www.example.com/profile",
    "preferredLanguage": "en-US",
    "userType": "Employee",
    "organization": "Okta",
    "title": "Director",
    "division": "R&D",
    "department": "Engineering",
    "costCenter": "10",
    "employeeNumber": "187",
    "mobilePhone": "+1-555-415-1337",
    "primaryPhone": "+1-555-514-1337",
    "streetAddress": "301 Brannan St.",
    "city": "San Francisco",
    "state": "CA",
    "zipCode": "94107",
    "countryCode": "US"
  },
  "credentials": {
    "password": {},
    "recovery_question": {
      "question": "Who's a major player in the cowboy scene?"
    },
    "provider": {
      "type": "OKTA",
      "name": "OKTA"
    }
  },
  "_links": {
    "resetPassword": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/lifecycle/reset_password"
    },
    "resetFactors": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/lifecycle/reset_factors"
    },
    "expirePassword": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/lifecycle/expire_password"
    },
    "forgotPassword": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/credentials/forgot_password"
    },
    "changeRecoveryQuestion": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/credentials/change_recovery_question"
    },
    "deactivate": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/lifecycle/deactivate"
    },
    "changePassword": {
      "href": "https://{yourOktaDomain}/api/v1/users/00ub0oNGTSWTBKOLGLNR/credentials/change_password"
    }
  }
}
*/

// User represents a user in Okta
//
// https://developer.okta.com/docs/api/resources/users#user-model
type User struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	Created         time.Time `json:"created"`
	Activated       time.Time `json:"activated"`
	StatusChanged   time.Time `json:"statusChanged"`
	LastLogin       time.Time `json:"lastLogin"`
	LastUpdated     time.Time `json:"lastUpdated"`
	PasswordChanged time.Time `json:"passwordChanged"`

	Profile UserProfile `json:"profile"`

	Credentials UserCredentials `json:"credentials"`

	Links struct {
		ResetPassword struct {
			Link string `json:"href"`
		} `json:"resetPassword"`
		ResetFactors struct {
			Link string `json:"href"`
		} `json:"resetFactors"`
		ExpirePassword struct {
			Link string `json:"href"`
		} `json:"expirePassword"`
		ForgotPassword struct {
			Link string `json:"href"`
		} `json:"forgotPassword"`
		ChangeRecoveryQuestion struct {
			Link string `json:"href"`
		} `json:"changeRecoveryQuestion"`
		Deactivate struct {
			Link string `json:"href"`
		} `json:"deactivate"`
		ChangePassword struct {
			Link string `json:"href"`
		} `json:"changePassword"`
	} `json:"_links"`
}

// UserProfile represents the profile object in Okta.
//
// https://developer.okta.com/docs/api/resources/users#profile-object
type UserProfile struct {
	Login             string `json:"login"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	NickName          string `json:"nickName"`
	DisplayName       string `json:"displayName"`
	Email             string `json:"email"`
	SecondEmail       string `json:"secondEmail"`
	ProfileURL        string `json:"profileUrl"`
	PreferredLanguage string `json:"preferredLanguage"`
	UserType          string `json:"userType"`
	Organization      string `json:"organization"`
	Title             string `json:"title"`
	Division          string `json:"division"`
	Department        string `json:"department"`
	CostCenter        string `json:"costCenter"`
	EmployeeNumber    string `json:"employeeNumber"`
	MobilePhone       string `json:"mobilePhone"`
	PrimaryPhone      string `json:"primaryPhone"`
	StreetAddress     string `json:"streetAddress"`
	City              string `json:"city"`
	State             string `json:"state"`
	ZipCode           string `json:"zipCode"`
	CountryCode       string `json:"countryCode"`
}

// UserCredentials represents the credentials object in Okta.
//
// https://developer.okta.com/docs/api/resources/users#credentials-object
type UserCredentials struct {
	Password struct {
		Value string `json:"value,omitempty"`
		Hash  struct {
			Algorithm  string `json:"algorithm"`
			WorkFactor int    `json:"workFactor"`
			Salt       string `json:"salt"`
			Value      string `json:"value"`
		} `json:"hash,omitempty"`
	} `json:"password"`
	RecoveryQuestion struct {
		Question string `json:"question"`
	} `json:"recovery_question"`
	Provider struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"provider"`
}
