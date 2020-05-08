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

## 运行（增加路径寻找功能）
```
[root@YZ01-Prometheus-Grafana1-19 tmp]# ./go_docker run -it sh
{"level":"info","msg":"command all is sh","time":"2020-05-07T17:42:29+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-07T17:42:29+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-07T17:42:29+08:00"}
sh-4.2# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 17:42 pts/1    00:00:00 sh
root         6     1  0 17:42 pts/1    00:00:00 ps -ef
```

## 运行（增加volume）
```
[root@yangzhou010010006017 ~]# ./go_docker run -it -v /root/volume:/containerVolume sh
{"level":"info","msg":"[\"/root/volume\" \"/containerVolume\"]","time":"2020-05-08T20:11:51+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-08T20:11:51+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-08T20:11:51+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-08T20:11:51+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-08T20:11:51+08:00"}
sh-4.2# ls
bin  containerVolume  dev  etc	home  proc  root  sys  tmp  usr  var
sh-4.2# pwd
/root/merged
```

```
sh-4.2# ls
bin  containerVolume  dev  etc	home  proc  root  sys  tmp  usr  var
sh-4.2# cd containerVolume/
sh-4.2# ls
sh-4.2# touch hahaha
sh-4.2# ls -al /root/volume/writeLayer/
总用量 0
drwxr-xr-x 2 root root 20 5月   8 20:12 .
drwxr-xr-x 5 root root 58 5月   8 20:11 ..
-rw-r--r-- 1 root root  0 5月   8 20:12 hahaha
```
