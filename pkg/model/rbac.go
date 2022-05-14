package model

const (
	AddPolicyAction    = "add"
	UpdatePolicyAction = "update"
	RemovePolicyAction = "remove"

	DefaultPolicyType   = "p"
	UserGroupPolicyType = "g"
	ResourcePolicyType  = "g2"
)

type Policy struct {
	Type      string   `json:"type"`
	Action    string   `json:"action"`
	Policy    []string `json:"policy"`
	OldPolicy []string `json:"oldPolicy"`
}
