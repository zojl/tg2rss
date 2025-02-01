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
      - LISTEN_PORT=80
      - MAX_TITLE_LENGTH=128
      - PYROGRAM_BRIDGE_HOST=http://pyrogram.example.com
      - PROXY_MEDIA=true
      - MEDIA_HOST=http://tgrss.example.tld
      - HOST_SECRET=SuperSecretPrivateKey
      - SAFE_HOST=tg2rss
```

After that just specify `http://tg2rss/yourchannel` as xml link (e.g. http://tg2rss/telegram)

### Configuration:
App configurations can be changed editing `.env` file or with setting environment variables. It's optional.  
- MAX_TITLE_LENGTH — how long post title can be (default 128 chars)
- LISTEN_PORT — server listening port (default 80)
- PYROGRAM_BRIDGE_HOST — to avoid the "Unsupported post" messages the [vvzvvlad's pyrogram bridge](https://github.com/vvzvlad/pyrogram-bridge) can be used to read posts like using MTProto application instead of web preview. Caution: telegram can limit such requests to prevent self-botting, and may block your account for such activity.
#### Temporary media links solution
Telegram web preview returns links to media that becoming obsolete after couple of days. To prevent this there's a bit complicated proxying/redirecting system that requires several more environment variables
- PROXY_MEDIA — set it to true to make this work
- MEDIA_HOST — specify host for links, so links to photos and videos in feed should be MEDIA_HOST/:channel/:id
- HOST_SECRET — if set, jwt token will be added to links, preventing unauthorized service as telegram media scrapper usage
- SAFE_HOST — secret or internal host to make requests without tokens.

For example, you can set MEDIA_HOST to http://tgrss.example.tld and SAFE_HOST to tg2rss. So rss reader deployed in docker can call the service using http://tg2rss/your_channel link and media redirect will be available at http://tgrss.example.tld/your_channel/12345.jpg?token=eyJhb...

See examples in .env file in the project