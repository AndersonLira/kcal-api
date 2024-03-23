#!/bin/bash

set -e

echo "starting build"
go build *.go
echo "build ok"

m=db
f="kcal-api"

arr=( $f )


echo "moving $m to $f"
mv $m $f
echo "moved"

remote="13.38.72.113"
user="ec2-user"
pemKey=~/.ssh/ALU-AWS-Lighsail-eu-west-3.pem

for var in "${arr[@]}"
do
    echo "copying $var to $remote"
    scp -i $pemKey $var $user@$remote:
    echo "copied"
done



set +e
echo "stoping service on $remote"
ssh -i $pemKey $user@$remote "pkill $f"
echo "stopped"

set -e
apps="apps"

for var in "${arr[@]}"
do
    echo "copying to $app/$var on $remote"
    ssh -i $pemKey $user@$remote "mv $var $apps/"
    echo "copied"
done

echo "starting service on $remote"
ssh -i $pemKey $user@$remote "export KCAL_API_DB_FOLDER=/home/$user/apps/data/;cd $apps/; ./$f >> /dev/null" &
echo "started"

echo "Done"

