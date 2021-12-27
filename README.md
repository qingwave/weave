
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
Run locally
```bash
make run
```

See http://localhost:8080

Test api
```bash
curl -XPOST http://localhost:8080/api/v1/users -d '{"name": "zhang3", "email": "zhang3@test.com"}'
```
