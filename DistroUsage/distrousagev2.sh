#!/bin/bash

pushd ~mirrorband/

YEST=$(date -d yesterday +%d/%b/%Y)

rm ./Yest.log

cat /var/log/nginx/access.log /var/log/nginx/access.log.1 | grep $YEST > ./Yest.log

./DU

popd