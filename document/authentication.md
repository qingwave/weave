# Authentication

Weave support RBAC access control model.

## Role
There is a example for cluster-admin role
```json
{
      "id": 1,
      "name": "cluster-admin",
      "scope": "cluster",
      "namespace": "",
      "rules": [
        {
          "id": 1,
          "roleId": 1,
          "resource": "*",
          "operation": "*"
        }
      ]
}
```

Role spec:
- scope: contains `cluster` and `namespace`
- namespace: if the scope is `namespace`, the namespace should be set with group name
- rules: every role has many rules, one rule contains resource and operation
- operation: 
  - `*`: contains all operation
  - `edit`: contains view and create, update, delete and patch operations
  - `view`, 'contains get and list
  - `get`, get a resource, api likes GET /api/v1/users/1
  - `list`, get all resources, api likes GET /api/v1/users
  - `create`, create a resource, api likes POST /api/v1/users
  - `update`, update a resource with some spec, POST /api/v1/users/2
  - `patch`, update a resource with all spec, PUT /api/v1/users/2
  - `delete`, delete a resource, DELETE /apia/v1/users/2
- resource: non-resource, likes /, /index, /healthz
  - containers, containers/log, containers/exec
  - users, users/groups, groups, groups/users
  - auth, for login and logout
  - posts, posts/like, posts/comment
  - namespaces
  - roles, rbac roles
  - k8s resources, pods, deployments, services and so on
  - some sub resources, `log`, `exec`, `proxy` for containers and pos 

## Default setting

Default Groups
- root, root group contains all cluster admin users, binding with a cluster-admin rool
- system:authenticated, all authenticated users belong to authenticated group
- system:unauthenticated, anonymous users belong to unauthenticated group, can access register user api

Default user
- admin, admin user belong to root group, is the cluster admin, can access all apis
- demo, demo user without any roles, only can get non-resource api
