falcon-agent for mac
===

This is a mac monitor agent. Just like zabbix-agent and tcollector.


## Installation

It is a golang classic project

```bash
# set $GOPATH and $GOROOT

git clone https://github.com/GitHamburg/agent-mac.git
cd agent-mac
./control build
./agent-mac start -c cfg.json

# goto http://localhost:1988
```

I use [linux-dash](https://github.com/afaqurk/linux-dash) as the page theme.

## Configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore

# Deployment

http://ulricqin.com/project/ops-updater/

