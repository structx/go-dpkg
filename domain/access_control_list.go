package domain

// AccessControlEntry application model
type AccessControlEntry struct {
	Subject    string     `json:"subject"`
	Resource   string     `json:"resource"`
	Permission Permission `json:"permission"`
}

// EntryResponse application model response to entry request
type EntryResponse struct {
	Granted bool `json:"granted"`
}
