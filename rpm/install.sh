#!/bin/bash

BN=$(basename $(pwd))
if [ "${BN}" != "rpm" ];then
   cd rpm
fi
IP=$(ip  addr|grep "\binet\b"|grep -v "127"|awk '{print $2}'|awk -F"/" '{print $1}'|xargs)
RPM=$(md5sum $(ls -l *.rpm|sort|egrep "(el6|el7)"|head -n 2|awk '{print $NF}')|awk '{print $2"|"$1}'|xargs|sed "s/ /,/g")
sed -i -e "s/Server_IP=\"127.0.0.1\"/Server_IP=\"${IP}\"/g" -e "s/Zabbix_RPM=.*/Zabbix_RPM=\"${RPM}\"/g" zabbix-agent.sh

