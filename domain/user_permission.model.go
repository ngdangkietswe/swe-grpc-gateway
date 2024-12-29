package domain

type Permission struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type UserPermission struct {
	Permissions []*Permission `json:"permissions"`
}
