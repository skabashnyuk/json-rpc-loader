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

KEY                            TYPE             DEFAULT    REQUIRED    DESCRIPTION
JSONRPCLOADER_CHE_HOST         String                      true        Che Server host
JSONRPCLOADER_CHE_TOKEN        String                                  User token for multi-user che
JSONRPCLOADER_MAJOR_THREADS    Integer          10                     Number of clients used to send message to major websocket endpoint
JSONRPCLOADER_MINOR_THREADS    Integer          10                     Number of clients used to send message to minor websocket endpoint
JSONRPCLOADER_WS_TIMEOUT       Duration         10s                    Websocket connection timeout
JSONRPCLOADER_SECURE           True or False    false                  Whatever secure websocket aka wss connection should be used
Configuration is set to:
CheHost: che-eclipse-che.192.168.64.67.nip.io
Token:
MajorThreads: 10
MinorThreads: 10
Timeout: 10s
Secure: false
Major rate 0/s  Minor rate 0/s
Major rate 23515/s  Minor rate 22678/s
```


## License

This project is licensed under the EPL License - see the [LICENSE.md](LICENSE.md) file for details
