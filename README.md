# go-log-agent
A client service to read logs from a file and send it to a server.

## Install

To install you can use makefile or build using commands

```
make build
OR
go build main.go
```

To run, 

```
make run OR ./go-log-agent
```

To run tests
```
make test
OR
go test ./... -v
```

## Usage

In order to use this you have to create a config YAML file to configure the service and file to watch.

The config file looks like this, currently it supports only watching single file. Default file is in "config/config.yaml"
```
filepath: /var/log/syslog  // full path of file to watch
servicename: syslog         // Service name 
serverurl: https://webhook.site/c11e67a9-198d-4f2e-a130-6604aaaa471f   // Server URL where data will be sent.
```

Now, run program using the above file.
```
go run main.go -c new_config.yaml
```

## Working 

The application constantly watches a file given in config for any changes, each line of log is parsed and converted to a LogStruct{}.

All the logs present are transferred to a buffred channel (100 buffer). If the channel is full the channel is flushed to send data to server in a POST request with list of LogStruct{} as json.

If the the channel is not full a periodic service is run to fulsh the data to server every 5 seconds.


![alt text](https://raw.githubusercontent.com/neeraj9194/go-log-agent/main/docs/arch.png)





