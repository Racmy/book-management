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

### 起動方法  
端末で
1. docker-compose up -d
2. localhost:80をブラウザで参照。何らかの時刻が出れば成功

### サーバーのログ確認方法
* 全て出力する場合  
1. docker-compose logs
* 最新の結果だけ出したい場合、毎回docker-compose logsをやらなくて済む
2. docker logs -f <コンテナID>
