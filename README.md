# go_hasura_tutorial

## setting and local development

- docker for mac setting
- hasura cli
## 以下のコマンドでローカル環境を作成する

```
$ git clone https://github.com/tmk616window/go_hasura_tutorial.git
$ cd go_hasura_tutorial
$ docker network create internal-api
$ docker-compose build
$ docker-compose up -d
$ cd hasura
$ hasura console
```

## データベース設計

https://user-images.githubusercontent.com/62201552/175231831-cc8a32d8-a175-45b5-b8b2-1feac395a2a6.png
