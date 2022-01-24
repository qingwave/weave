# Weave
Simple but functional Go+Vue application starter, supported by gin, gorm, redis, postgres, vue, element-plus, websocket and much more.

## Features
Server support features:
- Restful api, write by gin
- MVC structure
- Postgres storage, via gorm
- Swagger doc, support by swag
- Structured log, support by logrus
- Prometheus monitor
- PProf debug
- Graceful shutdown
- Authentication, support jwt
- Request rate limit, server level or user ip
- Redis cache

Frontend support features:
- Vue3 supported
- UI with element-plus
- Build with vite
- Charts integration, support by echarts
- WebShell supported
- Windi CSS

TODOs
- [x] Redis cache
- [x] Request rate limit
- [x] Authentication
- [x] WebSocket
- [ ] Trace
- [x] UI
- [x] WebShell
- [ ] Dark theme
- [ ] Mobile UI 
## Run
### Run server
install dependencies
```bash
make postgres
make redis
```

run locally
```bash
make run
```

See http://localhost:8080

### Test api
Register user
```bash
curl -XPOST http://localhost:8080/api/auth/user -d '{"name": "zhang3", "email": "zhang3@t.com","password": "123456"}'
```

Login, get jwt token
```bash
curl -XPOST http://localhost:8080/api/auth/token -d '{"name": "zhang3", "password": "123456"}'
```
Response as follows, set token in `Authorization Header`
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "xxx",
    "describe": "set token in Authorization Header, [Authorization: Bearer {token}]"
  }
}
```

Get users
```bash
token=xxx
curl -X 'GET' 'http://localhost:8080/api/v1/users' -H "Authorization: Bearer $token"
```

Container web shell
1. login in swagger `http://localhost:8080/swagger/index.html#/auth/post_login`
2. create container `http://localhost:8080/swagger/index.html#/user/post_api_v1_users`
```json
{
  "cmd": [
    "top"
  ],
  "image": "alpine",
  "name": "myapp"
}
```
3. open cloudshell `http://localhost:8080/api/v1/containers/{:containerid}/terminal`

### Run UI
```bash
cd web
npm i
npm run dev 
```

Explore in http://127.0.0.1:8081

Login page
![login](./document/img/login.png)

Dashboard page
![dashboard](./document/img/dashboard.png)

App page
![app](./document/img/app.png)

Webshell page
![webshell](./document/img/webshell.png)
