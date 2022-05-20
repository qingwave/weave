# Authentication

Weave support RBAC access control model, supported by [Casbin](https://github.com/casbin/casbin).

## Model
Model file, find in [](../config/auth_model.conf):
- Request Info: `sub, dom, obj, objName, act`
  eg, `bob, group1, containers, *, get` means bob in group group1 request containers resource with verb get
- Policy spec: `sub, dom, obj, objName, act`
  eg, `p, edit, tenant, tenant_sys_resource, *, "get"` means edit role in tenant group can get all resources in tenant_sys_resource resource group
- user role definition(`g`):  `user, role, group`, built-in roles contains `admin`, `edit`, `view` and cluster scope role `authenticated`, `unauthenticated`
  eg, `alice, admin, group1` means alice is admin role in group1
- resource role definition(`g2`): "resource, resource_group"
- matchers: support wildcard and slice, especially,  support authenticated role

## Policy
There are some default policies, see [auth_default_policy](../config/auth_default_policy.csv)

- `root` group is special built-in group contains administrators.
- `root.admin` can operate all resources in the system.  
- all user need have right to call auth(login/logout)
- `authenticated` users can get users info and containers(modify it ondemand)
