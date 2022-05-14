package authorization

import (
	"os"
	"strings"
	"weave/pkg/config"

	"github.com/Knetic/govaluate"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Enforcer *casbin.Enforcer

	DefaultRoles  = []string{AdminRole, EditRole, ViewRole}
	DefaultGroups = []string{RootGroup, TenantGroup}
)

const (
	policyName      = `p`
	groupPolicyName = `g`

	AdminRole = `admin`
	EditRole  = `edit`
	ViewRole  = `view`

	RootGroup   = `root`
	TenantGroup = `tenant`
)

func InitAuthorization(db *gorm.DB, conf config.AuthenticationConfig) error {
	a, err := gormadapter.NewAdapterByDBUseTableName(db, conf.AuthTablePrefix, conf.AuthTableName)
	if err != nil {
		return err
	}

	m, err := model.NewModelFromFile(conf.AuthModelConfigFullName)
	if err != nil {
		return err
	}

	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}
	Enforcer.EnableAutoSave(true)
	Enforcer.AddFunction("wMatch", matchFunc(weaveMatch))

	if err := Enforcer.LoadPolicy(); err != nil {
		return err
	}

	if (conf.LoadDefaultPolicyAlways || len(Enforcer.GetPolicy()) == 0) && isFileExists(conf.AuthDefaultPolicyConfigFullName) {
		fa := fileadapter.NewAdapter(conf.AuthDefaultPolicyConfigFullName)
		fm := m.Copy()
		if err := fa.LoadPolicy(fm); err != nil {
			return err
		}

		// add policy
		for ptype, assertion := range fm[policyName] {
			for _, p := range assertion.Policy {
				if !Enforcer.HasPolicy(p) {
					if _, err := Enforcer.AddNamedPolicy(ptype, p); err != nil {
						logrus.Infof("add policy %v failed: %v", p, err)
					}
				}
			}
		}

		// add group policy
		for ptype, assertion := range fm[groupPolicyName] {
			for _, p := range assertion.Policy {
				if !Enforcer.HasPolicy(p) {
					if _, err := Enforcer.AddNamedGroupingPolicy(ptype, p); err != nil {
						logrus.Infof("add group policy %v failed: %v", p, err)
					}
				}
			}
		}
	}

	return nil
}

func Enforce(user string, namespace string, resource string, resourceName string, verb string) bool {
	res, err := Enforcer.Enforce(user, namespace, resource, resourceName, verb)
	if err != nil {
		logrus.Warnf("failed to enforce user: %s, namespace: %s, resource: %s/%s, verb: %s, %v", user, namespace, resource, resourceName, verb, err)
		return false
	}
	logrus.Infof("enforce user: %s, namespace: %s, resource: %s/%s, verb: %s, %t", user, namespace, resource, resourceName, verb, res)

	return res
}

func isFileExists(name string) bool {
	f, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func weaveMatch(args ...interface{}) bool {
	if len(args) < 2 {
		return false
	}

	key1 := args[0].(string)
	key2 := args[1].(string)

	if key2 == "*" {
		return true
	}

	for _, key := range strings.Split(key2, ",") {
		if key1 == key {
			return true
		}
	}

	return key1 == key2
}

func matchFunc(f func(args ...interface{}) bool) govaluate.ExpressionFunction {
	return func(arguments ...interface{}) (interface{}, error) {
		return f(arguments...), nil
	}
}

func IsRootAdmin(user string) bool {
	ok, _ := Enforcer.HasRoleForUser(user, AdminRole, RootGroup)
	return ok
}
