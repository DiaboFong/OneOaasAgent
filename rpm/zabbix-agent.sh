#!/bin/bash

#ip  addr|grep "\binet\b"|grep -v "127"|awk '{print $2}'|awk -F"/" '{print $1}'|xargs
Server_IP="127.0.0.1"
Zabbix_RPM="zabbix-agent-3.0.15-1.el7.centos.zbx.x86_64.rpm|6fda7582c43c79808bb9c3b0d5016723,zabbix-agent-3.0.15-1.el6.zbx.x86_64.rpm|a1f1c803030b1600728a36e93cf07e51"

function down_cmd() {
	#which curl >/dev/null || c1=1 
	which wget >/dev/null || c2=2
	if [ "${c2}" == "2" ]; then
		yum install wget -y

	fi
}

function down_url() {
	VER=$(cat /etc/redhat-release | awk -F "." '{print $1}' | awk '{print $NF}')
	RPM=$(echo ${Zabbix_RPM} | sed "s/,/\n/g" | grep "el${VER}" | head -n 1)
	for IP in $(echo ${Server_IP} );
	do
	    ping -c 2 ${IP} >/dev/null
        if [ "$?" == "0" ];then
		    Server_IP=${IP}
			break 
	    fi
	done
	URL=http://${Server_IP}/download/${RPM}
	echo $URL
}

function install() {
	mkdir -p /tmp/.zabbix
	cd /tmp/.zabbix
	c=0
	for i in 1 2 3; do
		RPM_PKG=$(echo $(down_url) | awk -F"download/" '{print $NF}' | awk -F"|" '{print $1}')
		RPM_URL=$(down_url | awk -F"|" '{print $1}')
		if [ -f "${RPM_PKG}" ]; then
			rm -rf ./${RPM_PKG}
		fi
		wget ${RPM_URL} --connect-timeout=3
		RPM_MD5_ORIG=$(echo $(down_url) | awk -F"download/" '{print $NF}' | awk -F"|" '{print $2}')
		RPM_MD5_DOWN=$(md5sum ${RPM_PKG} | awk '{print $1}')
		if [ "${RPM_MD5_DOWN}" != "${RPM_MD5_ORIG}" ]; then
			echo "下载的文件MD5与脚本中文件的MD5不一致"
			c=$(($c + 1))
		else
			break
		fi
	done

	if [ "${c}" == "3" ]; then
		echo "下载的文件MD5与脚本中文件的MD5不一致,中断安装"
		exit 1
	fi
	remove
	yum localinstall ${RPM_PKG} -y
}

function remove() {
	rpm -qa | grep zabbix | grep agent | xargs rpm -e
}

down_cmd

case "$1" in
install)
	install
	;;
remove)
	remove
	;;
*)
	echo "Usage: $0 {install |remove }"
	;;
esac
