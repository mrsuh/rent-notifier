# Rent notifier

## Installation
```
sh bin/deploy.sh
set parameters in a config/config.yml
```

## Compilation
```sh
sh bin/compile.sh
```

## Run
```sh
bin/fixture config/config.yml fixtures
bin/api config/config.yml
bin/telegram config/config.yml
```

```sh
curl -X POST —data '{"url": "https://yourdomain.ru/{your_bot_token}/webhook"}' -H "Content-Type: application/json" "https://api.telegram.org/bot{your_bot_token}/setWebhook"
```