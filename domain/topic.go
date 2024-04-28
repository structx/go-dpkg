package domain

// RequestResponse application request/response topics
type RequestResponse string

const (
	// VerifyUserAccess verify user access
	VerifyUserAccess RequestResponse = "verify_user_access"
	// VerifyServiceAccess verify service access
	VerifyServiceAccess RequestResponse = "verify_service_access"
)

// String stringify request/response topic
func (rr RequestResponse) String() string {
	return string(rr)
}
