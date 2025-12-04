package iam

// Member represents a member of the organization
type Member struct {
	UUID     string  `json:"uuid"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	TenantID string  `json:"tenant_id"`
	Profile  *string `json:"profile,omitempty"`
}

// CreateMember represents the request to create a new member
type CreateMember struct {
	Email       string   `json:"email"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// Privileges represents the roles and permissions granted to a member
type Privileges struct {
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// EditGrant represents the request to edit grants (roles/permissions) for a member
type EditGrant struct {
	Operation   OperationType `json:"operation"`
	Roles       []string      `json:"roles,omitempty"`
	Permissions []string      `json:"permissions,omitempty"`
}

// OperationType represents the type of operation (add or remove)
type OperationType string

const (
	OperationAdd    OperationType = "add"
	OperationRemove OperationType = "remove"
)

// BatchUpdateMembers represents the request to batch update members
type BatchUpdateMembers struct {
	MemberIDs       []string      `json:"member_ids"`
	Operation       OperationType `json:"operation"`
	RoleNames       []string      `json:"role_names,omitempty"`
	PermissionNames []string      `json:"permission_names,omitempty"`
}

// Role represents a role in the organization
type Role struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Origin      string  `json:"origin"`
}

// CreateRole represents the request to create a new role
type CreateRole struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	BasedRole   *string  `json:"based_role,omitempty"`
}

// RolePermissions represents a role with its permissions
type RolePermissions struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Origin      string   `json:"origin"`
	Permissions []string `json:"permissions"`
}

// EditPermissions represents the request to edit role permissions
type EditPermissions struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// RolesMember represents a member with their roles
type RolesMember struct {
	MemberUUID string   `json:"member_uuid"`
	Roles      []string `json:"roles,omitempty"`
}

// Permission represents a permission action on a resource
type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Product represents a product with its permissions
type Product struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// AccessControl represents access control configuration
type AccessControl struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	TenantID    *string `json:"tenant_id,omitempty"`
	Enabled     bool    `json:"enabled"`
	EnforceMFA  bool    `json:"enforce_mfa"`
}

// AccessControlCreate represents the request to create access control
type AccessControlCreate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AccessControlStatus represents the request to update access control status
type AccessControlStatus struct {
	Status     *bool `json:"status,omitempty"`
	EnforceMFA *bool `json:"enforce_mfa,omitempty"`
}

// ServiceAccountDetail represents a service account with details
type ServiceAccountDetail struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Email       string  `json:"email"`
	Tenant      Tenant  `json:"tenant"`
}

// ServiceAccountCreate represents the request to create a service account
type ServiceAccountCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

// ServiceAccountEdit represents the request to edit a service account
type ServiceAccountEdit struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Tenant represents tenant information
type Tenant struct {
	UUID      string `json:"uuid"`
	LegalName string `json:"legal_name"`
}

// APIKeyServiceAccountDetail represents an API key for a service account
type APIKeyServiceAccountDetail struct {
	UUID                  string   `json:"uuid"`
	Name                  *string  `json:"name,omitempty"`
	Description           *string  `json:"description,omitempty"`
	KeyPairID             *string  `json:"key_pair_id,omitempty"`
	KeyPairSecret         *string  `json:"key_pair_secret,omitempty"`
	Scopes                []string `json:"scopes,omitempty"`
	ScopesPendingApproval []string `json:"scopes_pending_approval,omitempty"`
	StartValidity         *string  `json:"start_validity,omitempty"`
	EndValidity           *string  `json:"end_validity,omitempty"`
	RevokedAt             *string  `json:"revoked_at,omitempty"`
	RevokedBy             *string  `json:"revoked_by,omitempty"`
	APIKey                *string  `json:"api_key,omitempty"`
}

// APIKeyServiceAccountCreate represents the request to create an API key for a service account
type APIKeyServiceAccountCreate struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

// APIKeyServiceAccountEditInput represents the request to edit an API key
type APIKeyServiceAccountEditInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

// Scope represents a scope for API products
type Scope struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

// ApiProducts represents API products with scopes
type ApiProducts struct {
	UUID   string  `json:"uuid"`
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
}

// ScopeGroup represents a group of scopes
type ScopeGroup struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	APIProducts []ApiProducts `json:"api_products"`
}
