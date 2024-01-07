A microservice that allows to read public telegram channels as rss feeds.

Tested with miniflux as addition to `docker-compose.yaml`:
```yaml
services:
...
  tg2rss:
    image: zojl/tg2rss:dev
```

After that just specify `http://tg2rss/yourchannel` as xml link (e.g. http://tg2rss/telegram)