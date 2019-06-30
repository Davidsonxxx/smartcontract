# telegram-accountant-bot
Bot for managing criptocurrency accounts

In order it to work you need to create `config.json` with this content
```json
{
	"defaultLanguage" : "en-us",
	"extendedLog" : false,
	"updateIntervalSec" : 300,
	"availableLanguages" : [
		{"key": "en-us", "name": "English"}
	]
}
```
and `telegramApiToken.txt` that containts telegram API key for your bot.

## Install
Run this script to build
```
#!/bin/bash
bot_dir=github.com/gameraccoon
bot_name=telegram-accountant-bot
bot_exec=${bot_name}
go fmt ${bot_dir}/${bot_name}
go vet ${bot_dir}/${bot_name}
go test -v ${bot_dir}/${bot_name}/...
go install ${bot_dir}/${bot_name}
cp ${GOPATH}/bin/${bot_name} ./${bot_exec}
rm -rf "./data"
cp -r ${GOPATH}/src/${bot_dir}/${bot_name}/data ./
```
and this script to run
```
bot_name=telegram-accountant-bot
bot_exec=${bot_name}
mkdir -p logs
./${bot_exec} 2>> logs/log.txt 1>> logs/errors.txt & disown
```
