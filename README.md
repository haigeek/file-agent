# FileAgent
A tool read specific file by http server
## Quick Start
### RUN
```shell
./fileagent
```

when run success, you will see a http server running on port 8888
```
2024/08/11 22:49:38 http.go:86: listening http on 8888
```

### get file

when you get file, need to add token in header, the token key is `FA-Token`,you can customize it in config.yml

```shell
curl -H "FA-Token: mytoken" http://localhost:8888/read?file=test
```

```
hello fileAgent!
```

### update file
set file content in request body

```shell
curl -X POST \
     -H "FA-Token: mytoken"  \
     -d "This is new text." \
     http://127.0.0.1:8888/update?file=test
```

when get file, you will see new content
```
This is new text.
```

## Install

```shell
./fileagent --install
```


## Usage

when install success, you will see a config.yml file in current directory

### Config

```yaml
# http prot
port: 8888
# auth token
auth: mytoken
# file list
files:
  test: /tmp/test.txt
```

### Start

```shell
./fileagent --start
```

### Stop

```shell
./fileagent --stop
```

### UnInstall

```shell
./fileagent --uninstall
```





