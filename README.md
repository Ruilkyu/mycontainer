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

## 运行（增加volume功能）
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
## 运行(增加环境变量功能)
```
终端1：
./go_docker run -it --name envcontainer -e test1=123 -e test2=456 busybox sh
{"level":"info","msg":"createTty true","time":"2020-05-19T16:27:01+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-19T16:27:01+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-19T16:27:01+08:00"}
{"level":"info","msg":"Current location is /root/mnt/envcontainer","time":"2020-05-19T16:27:01+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-19T16:27:01+08:00"}
sh-4.2# env | grep test
test1=123
test2=456
终端2：
[root@yangzhou010010006012 srv]# ./go_docker ps
ID           NAME           PID         STATUS      COMMAND     CREATED
8015717058   envcontainer   49115       running     sh          2020-05-19 16:27:01
6593990918   hi             9329        running     sh          2020-05-17 13:32:53
[root@yangzhou010010006012 srv]# ./go_docker exec envcontainer sh
{"level":"info","msg":"container pid 49115","time":"2020-05-19T16:29:07+08:00"}
{"level":"info","msg":"command sh","time":"2020-05-19T16:29:07+08:00"}
sh-4.2# env | grep test
test1=123
test2=456
```
## 运行（增加创建网络功能）
```
./go_docker network create --driver bridge --subnet 10.200.0.1/24 testbridge1
查看网络：
ip a
6: testbridge1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UNKNOWN group default qlen 1000
    link/ether da:a5:ee:80:a5:18 brd ff:ff:ff:ff:ff:ff
    inet 10.200.0.1/24 brd 10.200.0.255 scope global testbridge1
       valid_lft forever preferred_lft forever
```
## 运行（增加容器互通功能）
```
终端1：
[root@yangzhou010010006012 srv]# ./go_docker run -it -net testbridge1 busybox sh
{"level":"info","msg":"createTty true","time":"2020-05-21T19:32:28+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-21T19:32:28+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-21T19:32:28+08:00"}
{"level":"info","msg":"Current location is /root/mnt/6252247727","time":"2020-05-21T19:32:28+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-21T19:32:28+08:00"}
sh-4.2# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
7: cif-62522@if8: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 72:c9:16:04:6f:74 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.200.0.2/24 brd 10.200.0.255 scope global cif-62522
       valid_lft forever preferred_lft forever
    inet6 fe80::70c9:16ff:fe04:6f74/64 scope link
       valid_lft forever preferred_lft forever
终端2：
[root@yangzhou010010006012 srv]# ./go_docker run -it -net testbridge1 busybox sh
{"level":"info","msg":"createTty true","time":"2020-05-21T19:32:55+08:00"}
{"level":"info","msg":"init come on","time":"2020-05-21T19:32:55+08:00"}
{"level":"info","msg":"command all is sh","time":"2020-05-21T19:32:55+08:00"}
{"level":"info","msg":"Current location is /root/mnt/4334394710","time":"2020-05-21T19:32:55+08:00"}
{"level":"info","msg":"find path /usr/bin/sh","time":"2020-05-21T19:32:55+08:00"}
sh-4.2# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
9: cif-43343@if10: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether fa:0e:7d:30:b4:04 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.200.0.3/24 brd 10.200.0.255 scope global cif-43343
       valid_lft forever preferred_lft forever
    inet6 fe80::f80e:7dff:fe30:b404/64 scope link
       valid_lft forever preferred_lft forever
终端1访问终端2：
sh-4.2# ping 10.200.0.3
PING 10.200.0.3 (10.200.0.3) 56(84) bytes of data.
64 bytes from 10.200.0.3: icmp_seq=1 ttl=64 time=0.125 ms
64 bytes from 10.200.0.3: icmp_seq=2 ttl=64 time=0.056 ms
64 bytes from 10.200.0.3: icmp_seq=3 ttl=64 time=0.056 ms
^C
--- 10.200.0.3 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 1999ms
rtt min/avg/max/mdev = 0.056/0.079/0.125/0.032 ms
```
## 运行（访问容器外网络）
```
sh-4.2# cat ./etc/resolv.conf
sh-4.2# echo "nameserver 114.114.114.114" > ./etc/resolv.conf
sh-4.2# cat ./etc/resolv.conf
nameserver 114.114.114.114
sh-4.2# ping www.baidu.com
PING www.wshifen.com (104.193.88.77) 56(84) bytes of data.
64 bytes from 104.193.88.77 (104.193.88.77): icmp_seq=1 ttl=40 time=146 ms
64 bytes from 104.193.88.77 (104.193.88.77): icmp_seq=2 ttl=40 time=146 ms
64 bytes from 104.193.88.77 (104.193.88.77): icmp_seq=3 ttl=40 time=146 ms
64 bytes from 104.193.88.77 (104.193.88.77): icmp_seq=4 ttl=40 time=146 ms
^C
--- www.wshifen.com ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 146.621/146.741/146.888/0.480 ms
```


