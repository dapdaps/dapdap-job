#!/usr/bin/env bash
export GOROOT=/root/go1.20.11
export GOPATH=/root/go
export PATH=GOROOT/bin:$PATH

basedir=$(cd "$(dirname "$0")"; pwd)
cd ${basedir}

function getpid() {
  pid=`ps -ef|grep bin/job | grep -v 'grep' | awk '{printf $2}'|cut -d/ -f1`
}

function build() {
    go build -o bin/job cmd/job/job.go
}

function start(){
    getpid;
    if [[ ! -z $pid ]];then
        echo "服务已运行中,pid:" $pid
        exit 1;
    fi
    cpu=`cat /proc/cpuinfo| grep "processor"| wc -l`
    ts=`expr $cpu \* 2`
    nohup ./bin/job -conf=config.toml  > web.out 2> web.err < /dev/null &
 }

function stop(){
    getpid;
    if [[ ! -z $pid ]];then
        echo "停止服务,pid:" $pid
    else
      echo "服务未启动,无需停止."
      exit 1;
    fi
    kill -9 $pid
}

function restart(){
    getpid;
    if [[ ! -z $pid ]];then
        echo "服务已运行中,pid:" $pid
            kill -9 $pid
        getpid;
        while [[ ! -z $pid ]]
        do
           sleep 1;
               getpid;
        done
    else
      echo "服务未启动,无需停止."
    fi
    start;

}

function status(){
    getpid;
    if [[ ! -z $pid ]];then
        echo "服务运行中,pid:" $pid
    else
        echo "服务未启动."
        exit 1;
    fi
}

function usage(){
    echo "$0 <start|stop|restart|status>"
}

case $1 in
    build)
      build;
      ;;
    start)
       start;
       ;;
    stop)
       stop;
       ;;
    restart)
       restart;
       ;;
    status)
        status;
        ;;
    *)
       usage;
       ;;
esac
