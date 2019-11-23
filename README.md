余田の本整理アプリ（GO, Nginx）

OverView

## Description
dockerを利用したgoの開発環境をgithubに上げます。  
【超絶参考】  
https://qiita.com/yoship1639/items/92405ab31779c8527c08

## Install
dockerをインストールして下さい。  
https://docs.docker.com/install/

## Usage

【起動方法】  
端末で
1. docker-compose up -d
2. localhost:80をブラウザで参照。何らかの時刻が出れば成功

【サーバーのログ確認方法】
1. docker-compose logs

【サーバへの入り方】
1. docker exec -it [containerID] /bin/bash
