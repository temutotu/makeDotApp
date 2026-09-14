#!/bin/bash
export AWS_ACCESS_KEY_ID=AKIAVMDSHRQVROSQBDKB
export AWS_SECRET_ACCESS_KEY=$(cat ~/.aws/secret_access_key)
export AWS_DEFAULT_REGION=ap-northeast-1

aws s3 cp s3://pokopia-dotmaker/releases/release.tar.gz .

if [ -f 'release']; then
        rm -rf release
fi
mkdir release
tar -zxvf release.tar.gz -C release

rsync -av release/ /usr/local/makeDotApp/release/