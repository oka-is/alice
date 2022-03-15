package engine

import (
	"github.com/wault-pw/alice/pkg/pack"
	"github.com/wault-pw/alice/server/policy"
)

type Opts struct {
	AllowOrigin     []string
	JwtKey          []byte
	CookieDomain    string
	CookieSecure    bool
	BackupUrl       string
	Ver             *pack.Ver
	UserPolicy      policy.IUserPolicy
	WorkspacePolicy policy.IWorkspacePolicy
}

func (o *Opts) SetDefaultPolicies() {
	o.UserPolicy = &policy.UserPolicy{}
	o.WorkspacePolicy = &policy.WorkspacePolicy{}
}
