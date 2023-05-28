# whoosh

## How to run?

```bash
git clone git@github.com:atomicai/whoosh.git
git clone git@github.com:atomicai/whooshui.git
cd whoosh
docker-compose up -d rethinkdb rabbitmq_whoosh
go run cmd/app/main.go
docker-compose up -d whooshui_external
```
