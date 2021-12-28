
# Weave
Simple but functional golang webserver example write by gin and gorm.

## Features
Support features:
- Restful api, write by gin
- MVC structure
- Postgres storage, via gorm
- Swagger doc, support by swag
- Structured log, support by logrus
- Prometheus monitor
- PProf debug
- Graceful shutdown

TODO
- [ ] Redis cache
- [x] Request rate limit
- [ ] Authentication
## Run
### Run locally
```bash
make run
```

See http://localhost:8080

### Test api
Register user
```bash
curl -XPOST http://localhost:8080/register -d '{"name": "zhang3", "email": "zhang3@t.com","password": "123456"}'
```

Login, get jwt token
```bash
curl -XPOST http://localhost:8080/login -d '{"name": "zhang3", "password": "123456"}'
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
