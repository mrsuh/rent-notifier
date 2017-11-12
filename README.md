# Rent notifier

## Installation
```
sh bin/deploy.sh
set parameters in a config/config.yml
```

## Load fixtures
```sh
bin/fixture config/config.yml /path/to/fixtures/dir
```

## Compilation
```sh
sh bin/compile.sh
```

## Run
```sh
bin/api config/config.yml
bin/telegram config/config.yml
bin/vk config/config.yml
```

## Configuration
config.config.yml
```yaml
#server api
api.listen: '127.0.0.1:8080'
api.prefix: 'api/v1'

#server telegram
telegram.listen: '127.0.0.1:8081'
telegram.token: 'telegram_token'
telegram.prefix: 'telegram_prefix'

#server vk
vk.listen: '127.0.0.1:8082'
vk.token: 'vk_token'
vk.confirm_secret: 'secret'
vk.prefix: 'vk_prefix'

log.file: '/Users/newuser/web/go/src/rent-notifier/file.log'
#database
database.dsn: localhost:5000

```

## Set telegram webhook
```sh
curl -X POST —data '{"url": "https://yourdomain.ru/{your_bot_token}/webhook"}' -H "Content-Type: application/json" "https://api.telegram.org/bot{your_bot_token}/setWebhook"
```