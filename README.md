Yubitoast
=========
This little CLI daemon keeps an eye on your `gpg-agent` log 
and when it detects a request to your Yubikey for signing or authenticating it will pop a small notification.

Notifications are done using the wonderful cross platform: https://github.com/martinlindhe/notify  
Eventhough it's deprecated it seems to work better than it's alternatives ...

Usage
-----
```
$ yubitoast -h
  -logfile string
    	path to gpg-agent.log (default "/var/log/gpg-agent.log")
```

GPG-Agent
---------
You need to make sure your `~/.gnupg/gpg-agent.conf` contains the following 2 lines:
```
log-file /var/log/gpg-agent.log
debug ipc
```

If you set the `log-file` to something different then you need to specific the `-logfile` arg.

After changing your `gpg-agent.conf` you need to `gpgconf --kill gpg-agent` to restart the agent.

Development
-----------
Code is located in `./src` and we accept that import paths include the `/src` because putting code in the root sucks ...

### Build
```
go build -o ./yubitoast ./src/cmd/yubitoast/main.go
```