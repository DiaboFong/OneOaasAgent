#!/bin/sh
#
# Startup script for Iagent
#
# chkconfig: - 85 18
# processname: iagent
# description: iagent
#
### BEGIN INIT INFO
# Provides: iagent
# Required-Start: $local_fs $remote_fs $network
# Required-Stop: $local_fs $remote_fs $network
# Default-Start: 2 3 4 5
# Default-Stop: 0 1 6
# Short-Description: start and stop iagent
### END INIT INFO

ulimit -n 30000

. /etc/rc.d/init.d/functions

RETVAL=0

start() {
    nohup ./iagent & > /dev/null 2>&1

    sleep 1 && clear
    pid=$(ps -ef |grep iagent|grep -v "grep"|awk '{print $2}')
    if [ "${pid}" == "" ];then
       echo "启动失败,请手动执行./iagent"
    else
        echo -e "------------欢迎使用iagent------------\n"
        echo -e "-----访问地址http://您的本机IP:$80-----\n"
    fi
}

stop() {
    echo -n $"Stopping $prog: "
    pid=$(ps -ef |grep iagentgrep -v "grep"|awk '{print $2}')
    if [ "${pid}" != "" ];then
        echo ${pid} |xargs kill -9
    fi
}


# See how we were called.
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo $"Usage: $prog {start|stop|restart}"
        RETVAL=2
esac

exit $RETVAL