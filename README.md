# json rpc loader

This is a tool for Eclipse Che json rpc permomance tests



### Installing


```
docker pull ksmster/json-rpc-loader
```


## Running the tests
```
docker run -e JSONRPCLOADER_CHE_HOST=che-eclipse-che.192.168.64.67.nip ksmster/json-rpc-loader 
```
Output
```
This application is configured via the environment. The following environment
variables can be used:

KEY                           TYPE             DEFAULT    REQUIRED    DESCRIPTION
JSONRPCLOADER_CHE_HOST        String                      true        Che Server host
JSONRPCLOADER_CHE_TOKEN       String                                  User token for multi-user che
JSONRPCLOADER_CLIENT          Integer          10                     Number of clients used to send messages
JSONRPCLOADER_WS_TIMEOUT      Duration         10s                    Websocket connection timeout 
JSONRPCLOADER_SECURE          True or False    false                  Whatever secure websocket aka wss connection should be used
JSONRPCLOADER_MULTIPLEXING    True or False    false                  Whatever use single websocket connection by each client to send request
Configuration is set to:
CheHost: che-eclipse-che.192.168.64.67.nip.io
Token: 
Threads: 10
Timeout: 10s
Secure: false
Multiplexing: false
Iteration at 0/5s 
Iteration at 384/5s 
Iteration at 202/5s 
Iteration at 223/5s 
Iteration at 179/5s 
Iteration at 158/5s 
Iteration at 241/5s 
```


## License

This project is licensed under the EPL License - see the [LICENSE.md](LICENSE.md) file for details
