package nsenter


/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
__attribute__((constructor)) void enter_namespace(void) {
	char *mycontainer_pid;
    // 从环境变量获取需要进入的pid
	mycontainer_pid = getenv("mycontainer_pid");
	if (mycontainer_pid) {
		//fprintf(stdout, "got mycontainer_pid=%s\n", mycontainer_pid);
	} else {
		//fprintf(stdout, "missing mycontainer_pid env skip nsenter");
		return;
	}
	char *mycontainer_cmd;
	// 从环境变量获取需要执行的命令
	mycontainer_cmd = getenv("mycontainer_cmd");
	if (mycontainer_cmd) {
		//fprintf(stdout, "got mycontainer_cmd=%s\n", mycontainer_cmd);
	} else {
		//fprintf(stdout, "missing mycontainer_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };
	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", mycontainer_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
    // 在进入的namespace中执行特定的命令
	int res = system(mycontainer_cmd);
	exit(0);
	return;
}
*/
import "C"
