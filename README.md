A microservice that allows to read public telegram channels as rss feeds.

Tested with miniflux as addition to `docker-compose.yaml`:
```yaml
services:
...
  tg2rss:
    image: zojl/tg2rss:dev
```

This example is too simple and doesn't support media proxying.
The more complex config can look like:
```yaml
services
...
  tg2rss:
    image: zojl/tg2rss:dev
    ports:
      - "8080:80"
    environment:
      - MEDIA_HOST=https://tg2rss.example.com
      - PROXY_MEDIA=true
      - LISTEN_PORT=80
      - MAX_TITLE_LENGTH=128
```

After that just specify `http://tg2rss/yourchannel` as xml link (e.g. http://tg2rss/telegram)

### Configuration:
App configurations can be changed editing `.env` file or with setting environment variables. It's optional.  
- MAX_TITLE_LENGTH — how long post title can be (default 128 chars)
- LISTEN_PORT — server listening port (default 80)
- PROXY_MEDIA — if false the service will use direct links to photos and videos (but this links will sometimes become obsolete), if true the service will return links to itself containing redirects to actual media files