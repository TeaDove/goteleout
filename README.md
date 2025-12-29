# Teleout
Golang version of [teleout](https://github.com/teadove/teleout)
Pipe stdout and files(TBD) to telegram
This software uses telegram bots
Powered with love

# Examples
- `ls -la | goteleout -c` - send output of `ls -la` with monospace font
- `goteleout -q -f main.go` - send file `main.go` without notification

# Features
1. Send files, text messages directly to telegram
2. Pipe to teleout(`ls | goteleout` will work)
3. HTML parse mode supported
4. Easy install and use
5. GetMe - listen for updates and log chatId, userName etc. Userfull for getting this information.

# Manual
```shell
NAME:
   goteleout - pipes data to telegram, https://github.com/teadove/goteleout

USAGE:
   goteleout [options] "text to send"

GLOBAL OPTIONS:
   --code, -c           send text with <code> tag to make it monospace, automatically set parseMode=HTML and escapes content
   --quite, -q          send message without notifications
   --parse-mode string  sets parse mode, can be: HTML, Markdown, MarkdownV2
   --file, -f           specify files to send
   --help, -h           show help
```

# Installation
1. From source code
```
go install github.com/teadove/goteleout@latest
```
Or download your version from [release page](https://github.com/TeaDove/goteleout/releases), i.e. for Apple Silicon
```
wget -cO - https://github.com/TeaDove/goteleout/releases/download/1.1.3/goteleout-1.1.3-darwin-arm64 > goteleout
chmod u+x goteleout
# mv goteleout ~/.local/bin # or any other location in your PATH
```
2. Get bot token from [@BotFather](https://t.me/BotFather)
3. Put config data in `~/.config/goteleout.json` in format `{"token": <telegram-token>, "user": <telegram-user-id>}`
4. Run `ls -la | goteleout -c`

> don't worry, there are no sniffer and smth like that

inspired by https://termbin.com
