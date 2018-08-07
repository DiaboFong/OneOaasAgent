#!/bin/bash

mkdir -p pkg/conf
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -ldflags "-s -w"
cp conf/app.conf pkg/conf/app.conf
cp conf/iagent.sh pkg/iagent.sh
cp -r rpm pkg/rpm
cp -r views static pkg
cp iagent pkg/iagent
mv pkg iagent_v1.0 && tar zcvf iagent_v1.0.tar.gz iagent_v1.0
rm -rf iagent_v1.0