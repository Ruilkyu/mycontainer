# mycontainer
实现自己的容器引擎
## 环境
```
实验环境：CentOS7  
```
## 编译
```
不编译cgo:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
编译cgo:
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build .
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

## 运行（增加打包容器功能）
```
终端1：
./go_docker run -it sh
[root@yangzhou010010006012 srv]# ./go_docker run -it sh
{"level":"error","msg":"Untar dir /root/busybox/ error exit status 2","time":"2020-05-15T18:01:51+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-15T18:01:51+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-15T18:01:51+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-15T18:01:51+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-15T18:01:51+08:00"}
sh-4.2#
终端2：
[root@yangzhou010010006012 srv]# ls /root
busybox  merged  workLayer  writeLayer
[root@yangzhou010010006012 srv]# ./go_docker commit image
/root/image.tar[root@yangzhou010010006012 srv]# ls /root
busybox  image.tar  merged  workLayer  writeLayer
```
## 运行（增加ps功能）
```
终端1：
[root@yangzhou010010006012 srv]# ./go_docker run -it sh
{"level":"info","msg":"createTty true","time":"2020-05-16T21:15:21+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-16T21:15:21+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-16T21:15:21+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-16T21:15:21+08:00"}
sh-4.2#
终端2：
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME         PID         STATUS      COMMAND     CREATED
8329068587   8329068587   32493       running     sh          2020-05-16 21:15:21
[root@yangzhou010010006012 srv]# cat /var/run/mycontainer/8329068587/config.json
{"pid":"32493","id":"8329068587","name":"8329068587","command":"sh","createdTime":"2020-05-16 21:15:21","status":"running"}
```
## 运行（增加logs功能）
```
终端1：
root@yangzhou010010006012 srv]# ./go_docker run -it -l --name hi sh
{"level":"info","msg":"createTty true","time":"2020-05-17T13:32:53+08:00"}
{"level":"error","msg":"Mkdir dir /root/writeLayer/ error. mkdir /root/writeLayer/: file exists","time":"2020-05-17T13:32:53+08:00"}
{"level":"error","msg":"Mkdir dir /root/workLayer/ error. mkdir /root/workLayer/: file exists","time":"2020-05-17T13:32:53+08:00"}
{"level":"error","msg":"Mkdir dir /root/merged/ error. mkdir /root/merged/: file exists","time":"2020-05-17T13:32:53+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-17T13:32:53+08:00"}

终端2：
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
[root@yangzhou010010006012 srv]# ./go_docker logs hi
{"level":"info","msg":"init come on","time":"2020-05-17T13:32:53+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-17T13:32:53+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-17T13:32:53+08:00"}
```
## 运行（增加exec功能）
```
终端1:
[root@yangzhou010010006012 srv]# ./go_docker run -it --name mycontainer sh
{"level":"info","msg":"createTty true","time":"2020-05-18T13:41:58+08:00"}
{"level":"error","msg":"Mkdir dir /root/writeLayer/ error. mkdir /root/writeLayer/: file exists","time":"2020-05-18T13:41:58+08:00"}
{"level":"error","msg":"Mkdir dir /root/workLayer/ error. mkdir /root/workLayer/: file exists","time":"2020-05-18T13:41:58+08:00"}
{"level":"error","msg":"Mkdir dir /root/merged/ error. mkdir /root/merged/: file exists","time":"2020-05-18T13:41:58+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-18T13:41:58+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-18T13:41:58+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-18T13:41:58+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-18T13:41:58+08:00"}
sh-4.2#
终端2:
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME          PID         STATUS      COMMAND     CREATED
6593990918   hi            9329        running     sh          2020-05-17 13:32:53
3631446183   mycontainer   45413       running     sh          2020-05-18 13:41:58
[root@yangzhou010010006012 srv]# ./go_docker exec mycontainer sh
{"level":"info","msg":"container pid 45413","time":"2020-05-18T13:42:46+08:00"}
{"level":"info","msg":"command sh","time":"2020-05-18T13:42:46+08:00"}
sh-4.2#
```
## 运行（增加stop功能）
```
终端1：
./go_docker run -it --name mycontainer top
终端2：
[root@yangzhou010010006012 srv]# ps -ef | grep top
root      8478 44941  0 14:33 pts/1    00:00:00 ./go_docker run -it --name mycontainer top
root      8488  8478  0 14:33 pts/1    00:00:00 top
root      8509 44765  0 14:33 pts/0    00:00:00 grep --color=auto top
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME          PID         STATUS      COMMAND     CREATED
6593990918   hi            9329        running     sh          2020-05-17 13:32:53
6367279090   mycontainer   8488        running     top         2020-05-18 14:33:26
[root@yangzhou010010006012 srv]# ./go_docker stop mycontainer
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
```
## 运行（增加rm功能）
```
终端1：
[root@yangzhou010010006012 srv]# ./go_docker run -it --name container sh
{"level":"info","msg":"createTty true","time":"2020-05-18T16:08:25+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-18T16:08:25+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-18T16:08:25+08:00"}
{"level":"info","msg":"Current location is /root/merged","time":"2020-05-18T16:08:25+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-18T16:08:25+08:00"}
sh-4.2#
终端2：
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6176477697   container   30494       running     sh          2020-05-18 16:08:25
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
[root@yangzhou010010006012 srv]# ./go_docker rm container
{"level":"error","msg":"Couldn't remove running container","time":"2020-05-18T16:08:55+08:00"}
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6176477697   container   30494       running     sh          2020-05-18 16:08:25
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
[root@yangzhou010010006012 srv]# ./go_docker stop container
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6176477697   container               stopped     sh          2020-05-18 16:08:25
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
[root@yangzhou010010006012 srv]# ./go_docker rm container
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
6593990918   hi          9329        running     sh          2020-05-17 13:32:53
```
## 运行（增加容器制作镜像功能）
```
终端1：
./go_docker run -it --name container1 -v /root/from1:to1 busybox sh
{"level":"info","msg":"createTty true","time":"2020-05-19T14:05:52+08:00"}
{"level":"info","msg":"NewWorkSpace volume urls [\"/root/from1\" \"to1\"]","time":"2020-05-19T14:05:52+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-19T14:05:52+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-19T14:05:52+08:00"}
{"level":"info","msg":"Current location is /root/mnt/container1","time":"2020-05-19T14:05:52+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-19T14:05:52+08:00"}
sh-4.2#
终端2：
./go_docker run -it --name container2 -v /root/from2:to2 busybox sh
{"level":"info","msg":"createTty true","time":"2020-05-19T14:12:47+08:00"}
{"level":"info","msg":"NewWorkSpace volume urls [\"/root/from2\" \"to2\"]","time":"2020-05-19T14:12:47+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-19T14:12:47+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-19T14:12:47+08:00"}
{"level":"info","msg":"Current location is /root/mnt/container2","time":"2020-05-19T14:12:47+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-19T14:12:47+08:00"}
sh-4.2#
终端3：
[root@yangzhou010010006012 ~]# ls /root
busybox.tar  from1  from2  mnt  readLayer  workLayer  writeLayer
[root@yangzhou010010006012 ~]# ls /root/mnt/
container1  container2
[root@yangzhou010010006012 ~]# ls /root/mnt/container1
bin  dev  etc  home  proc  root  sys  tmp  to1  usr  var
[root@yangzhou010010006012 ~]# ls /root/mnt/container2
bin  dev  etc  home  proc  root  sys  tmp  to2  usr  var
[root@yangzhou010010006012 ~]# ls /root/readLayer/
busybox
[root@yangzhou010010006012 ~]# ls /root/workLayer/
container1  container2
[root@yangzhou010010006012 ~]# ls /root/writeLayer/
container1  container2
将容器制作镜像：
./go_docker commit container1 image1
[root@yangzhou010010006012 srv]# ls /root
busybox.tar  from1  from2  image1.tar  mnt  readLayer  workLayer  writeLayer
```
