package mcptools

import (
	"fmt"
	"os"
	"strings"
)

// EmailValidator checks email format and domain against an allowed list.
// Allowed domains are loaded from the ALLOWED_EMAIL_DOMAINS env var
// as a comma-separated list (e.g. "school.com,students.school.com").
// If the env var is not set, only format is validated.
type EmailValidator struct {
	allowedDomains []string
}

func NewEmailValidator() *EmailValidator {
	var allowed []string
	if raw := os.Getenv("ALLOWED_EMAIL_DOMAINS"); raw != "" {
		for _, d := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(strings.ToLower(d)); trimmed != "" {
				allowed = append(allowed, trimmed)
			}
		}
	}
	return &EmailValidator{allowedDomains: allowed}
}

// Validate returns an error if the email fails format or domain checks.
func (v *EmailValidator) Validate(email string) error {
	// Format check
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("'%s' is not a valid email address", email)
	}

	domain := strings.ToLower(parts[1])

	// Must have at least one dot in the domain
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("'%s' is not a valid email address", email)
	}

	// Domain whitelist check
	if len(v.allowedDomains) > 0 {
		for _, allowed := range v.allowedDomains {
			if domain == allowed {
				return nil
			}
		}
		return fmt.Errorf(
			"email domain '@%s' is not allowed — accepted domains: %s",
			domain,
			strings.Join(v.allowedDomains, ", "),
		)
	}

	return nil
}
