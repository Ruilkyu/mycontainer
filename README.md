# mycontainer
实现自己的容器引擎
## 编译
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

## 运行
```
[root@yangzhou010010001015 srv]# ./go_docker run -it /bin/sh
{"level":"info","msg":"init come on","time":"2020-05-05T17:32:48+08:00"}
{"level":"info","msg":"command /bin/sh","time":"2020-05-05T17:32:48+08:00"}
{"level":"info","msg":"command /bin/sh","time":"2020-05-05T17:32:48+08:00"}
sh-4.2# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 17:32 pts/8    00:00:00 /bin/sh
root         6     1  0 17:32 pts/8    00:00:00 ps -ef

```
