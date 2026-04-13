package mcptools_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/periasamy/school-api/internal/mcptools"
)

func TestEmailValidator_NoDomainsConfigured(t *testing.T) {
	// No ALLOWED_EMAIL_DOMAINS set — any valid format should pass
	v := mcptools.NewEmailValidator()

	assert.NoError(t, v.Validate("alice@school.com"))
	assert.NoError(t, v.Validate("bob@gmail.com"))
	assert.NoError(t, v.Validate("user@company.org"))
}

func TestEmailValidator_InvalidFormat(t *testing.T) {
	v := mcptools.NewEmailValidator()

	assert.Error(t, v.Validate("notanemail"))
	assert.Error(t, v.Validate("@nodomain"))
	assert.Error(t, v.Validate("noatsign"))
	assert.Error(t, v.Validate("missing@dot"))
	assert.Error(t, v.Validate(""))
}

func TestEmailValidator_DomainWhitelist(t *testing.T) {
	t.Setenv("ALLOWED_EMAIL_DOMAINS", "school.com,students.school.com")
	v := mcptools.NewEmailValidator()

	// Allowed domains
	assert.NoError(t, v.Validate("alice@school.com"))
	assert.NoError(t, v.Validate("bob@students.school.com"))

	// Not allowed
	assert.Error(t, v.Validate("alice@gmail.com"))
	assert.Error(t, v.Validate("bob@yahoo.com"))
	assert.Error(t, v.Validate("user@school.org"))
}

func TestEmailValidator_DomainCaseInsensitive(t *testing.T) {
	t.Setenv("ALLOWED_EMAIL_DOMAINS", "school.com")
	v := mcptools.NewEmailValidator()

	assert.NoError(t, v.Validate("alice@SCHOOL.COM"))
	assert.NoError(t, v.Validate("alice@School.Com"))
}

func TestEmailValidator_WhitespaceInDomainList(t *testing.T) {
	t.Setenv("ALLOWED_EMAIL_DOMAINS", " school.com , students.school.com ")
	v := mcptools.NewEmailValidator()

	assert.NoError(t, v.Validate("alice@school.com"))
	assert.NoError(t, v.Validate("bob@students.school.com"))
}
