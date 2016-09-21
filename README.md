falcon-agent for mac
===

This is a mac monitor agent. Just like zabbix-agent and tcollector.


## Installation

It is a golang classic project

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/common.git
cd ..
#https://github.com/toolkits
#toolkits放入$GOPATH/src/github.com目录(https://github.com/GitHamburg/agent-mac/blob/master/resources/toolkits.tar)

#set agent-mac
#把agent-mac放在非gopath路径下,否则会报:local import "./cron" in non-local package异常
git clone https://github.com/GitHamburg/agent-mac.git
cd agent-mac
./control build
#注意在cfg.json中的hostname/ip,设置本机ip地址;在addr处,加入服务端ip地址
./falcon-agent start -c cfg.json

# goto http://localhost:1988
```

I use [linux-dash](https://github.com/afaqurk/linux-dash) as the page theme.

## Configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore

# Deployment

http://ulricqin.com/project/ops-updater/

