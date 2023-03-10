# Teleout
Golang version of [teleout](https://github.com/teadove/teleout)
Pipe stdout and files(TBD) to telegram
This software uses telegram bots
Powered with love

# Examples
- `ls -la | goteleout -u 418878871 -c` - send output of `ls -la` to chat `418878871` with monospace font
- `teleout -q -f main.py` - send file `main.py` without notification to default user

# Features
1. Send files, text messages directly to telegram
2. Pipe to teleout(`ls | teleout` will work)
3. HTML parse mode supported
4. Easy install and use
5. GetMe - listen for updates and log chatId, userName etc. Userfull for getting this information.

# Manual
```shell
NAME:
   goteleout - pipes data to telegram, https://github.com/teadove/goteleout

USAGE:
   goteleout [options] "text to send"

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --code, -c              send text with <code> tag to make it monospace (default: false)
   --quite, -q             send message without notifications (default: false)
   --html                  do no escape html tags (default: false)
   --token value           telegram api token [$TELEGRAM_TOKEN]
   --user value, -u value  telegram user id
   --settings-file value   file to store default settings
   --verbose, -v           (default: false)
   --get-me                will listen for updates, and reply with user_id and chat_id (default: false)
   --help, -h              show help (default: false)
```

# Installation
1. From source code
```
go install github.com/teadove/goteleout
```
Or download your version from [release page](https://github.com/TeaDove/goteleout/releases), i.e. for Apple Silicon
```
wget -cO - https://github.com/TeaDove/goteleout/releases/download/1.0.0/goteleout-1.0.0-darwin-arm64 > goteleout
chmod u+x goteleout
# mv goteleout ~/.local/bin # or any other location in your PATH
```
2. Get bot token from [@BotFather](https://t.me/BotFather)
3. Put config data in `~/.config/teleout.json` in format `{"token": <telegram-token>, "user": <telegram-user-id>}`
4. Run `ls -la | goteleout -c`

> don't worry, there are no sniffer and smth like that

> for feedbacks, write me [here](https://t.me/teas_feedbacks_bot)<br>
inspired by https://termbin.com
