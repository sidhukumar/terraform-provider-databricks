package secrets

import (
	"testing"

	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/terraform-provider-databricks/qa"
	"github.com/stretchr/testify/assert"
)

func TestResourceSecretACLRead(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/secrets/acls/get?principal=something&scope=global",
				Response: ACLItem{
					Permission: "CAN_MANAGE",
				},
			},
		},
		Resource: ResourceSecretACL(),
		Read:     true,
		ID:       "global|||something",
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, "global|||something", d.Id(), "Id should not be empty")
	assert.Equal(t, "CAN_MANAGE", d.Get("permission"))
	assert.Equal(t, "something", d.Get("principal"))
	assert.Equal(t, "global", d.Get("scope"))
}

func TestResourceSecretACLRead_NotFound(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/secrets/acls/get?principal=something&scope=global",
				Response: apierr.APIErrorBody{
					ErrorCode: "NOT_FOUND",
					Message:   "Item not found",
				},
				Status: 404,
			},
		},
		Resource: ResourceSecretACL(),
		Read:     true,
		Removed:  true,
		ID:       "global|||something",
	}.ApplyNoError(t)
}

func TestResourceSecretACLRead_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/secrets/acls/get?principal=something&scope=global",
				Response: apierr.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceSecretACL(),
		Read:     true,
		ID:       "global|||something",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "global|||something", d.Id(), "Id should not be empty for error reads")
}

func TestResourceSecretACLCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/secrets/acls/put",
				ExpectedRequest: SecretACLRequest{
					Principal:  "something",
					Permission: "CAN_MANAGE",
					Scope:      "global",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/secrets/acls/get?principal=something&scope=global",
				Response: ACLItem{
					Permission: "CAN_MANAGE",
				},
			},
		},
		Resource: ResourceSecretACL(),
		State: map[string]any{
			"permission": "CAN_MANAGE",
			"principal":  "something",
			"scope":      "global",
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, "global|||something", d.Id())
}

func TestResourceSecretACLCreate_ScopeWithSlash(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/secrets/acls/put",
				ExpectedRequest: SecretACLRequest{
					Principal:  "something",
					Permission: "CAN_MANAGE",
					Scope:      "myapplication/branch",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/secrets/acls/get?principal=something&scope=myapplication%2Fbranch",
				Response: ACLItem{
					Permission: "CAN_MANAGE",
				},
			},
		},
		Resource: ResourceSecretACL(),
		State: map[string]any{
			"permission": "CAN_MANAGE",
			"principal":  "something",
			"scope":      "myapplication/branch",
		},
		Create: true,
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, "myapplication/branch|||something", d.Id())
}

func TestResourceSecretACLCreate_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{ // read log output for better stub url...
				Method:   "POST",
				Resource: "/api/2.0/secrets/acls/put",
				Response: apierr.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceSecretACL(),
		State: map[string]any{
			"permission": "CAN_MANAGE",
			"principal":  "something",
			"scope":      "global",
		},
		Create: true,
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "", d.Id(), "Id should be empty for error creates")
}

func TestResourceSecretACLDelete(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/secrets/acls/delete",
				ExpectedRequest: map[string]string{
					"scope":     "global",
					"principal": "something",
				},
			},
		},
		Resource: ResourceSecretACL(),
		Delete:   true,
		ID:       "global|||something",
	}.Apply(t)
	assert.NoError(t, err)
	assert.Equal(t, "global|||something", d.Id())
}

func TestResourceSecretACLDelete_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/secrets/acls/delete",
				Response: apierr.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceSecretACL(),
		Delete:   true,
		ID:       "global|||something",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "global|||something", d.Id())
}
