# json rpc loader

This is a tool for Eclipse Che json rpc permomance tests



### Installing


```
docker pull ksmster/json-rpc-loader
```


## Running the tests
```
docker run ksmster/json-rpc-loader \
       -cheurl  ws://che-eclipse-che.192.168.64.12.nip.io/api/websocket \
       -token=t1  \
       -mnum=10 \
       -tnum=1
```
Output
```
[12:19:42]sj:json-rpc-loader[master]#: docker run ksmster/json-rpc-loader \
>        -cheurl  ws://che-eclipse-che.192.168.64.12.nip.io/api/websocket \
>        -token=t1  \
>        -mnum=10 \
>        -tnum=1
thum: 1
mnum: 10
cheurl: ws://che-eclipse-che.192.168.64.12.nip.io/api/websocket
cheToken: t1
2018/12/16 10:19:44 Messaget from thread 0 number 0
2018/12/16 10:19:44 Messaget from thread 0 number 1
2018/12/16 10:19:44 Messaget from thread 0 number 2
2018/12/16 10:19:44 Messaget from thread 0 number 3
2018/12/16 10:19:44 Messaget from thread 0 number 4
2018/12/16 10:19:44 Messaget from thread 0 number 5
2018/12/16 10:19:44 Messaget from thread 0 number 6
2018/12/16 10:19:44 Messaget from thread 0 number 7
2018/12/16 10:19:44 Messaget from thread 0 number 8
2018/12/16 10:19:44 Messaget from thread 0 number 9
```


## License

This project is licensed under the EPL License - see the [LICENSE.md](LICENSE.md) file for details
