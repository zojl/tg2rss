A microservice that allows to read public telegram channels as rss feeds.

Tested with miniflux as addition to `docker-compose.yaml`:
```yaml
services:
...
  tg2rss:
    image: zojl/tg2rss:dev
```

After that just specify `http://tg2rss/yourchannel` as xml link (e.g. http://tg2rss/telegram)

### Configuration:
App configurations can be changed editing `.env` file or with setting environment variables. It's optional.  
- MAX_TITLE_LENGTH — how long post title can be (default 128 chars)
- LISTEN_PORT — server listening port (default 80)