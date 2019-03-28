# Lambda installation guide

First of all, [Golang installation](https://golang.org/dl) needs to be done.

**Development version:** Go 1.11.2 Linux / amd64

## Required Go librarys
[Beego](https://beego.me/docs/install/bee.md) as api framework

Cors plugin for Beego

    go get github.com/astaxie/beego/plugins/cors

Mysql driver for Go

    go get github.com/go-sql-driver/mysql

Redis driver for Go

    go get github.com/go-redis/redis

For scheduled jobs run

    go get github.com/robfig/cron

JWT token library

    go get github.com/juusechec/jwt-beego

> You need to change public key and private key at `/key` directory as name: `rsakey.pem`, `rsakey.pem.pub`

Create private key: `openssl genpkey -out rsakey.pem -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -pass pass:verysecret`

Create public key: `openssl rsa -in rsakey.pem -pubout > rsakey.pem.pub`

Lambda uses Mysql as a data base operations. [Installation link](https://www.mysql.com/downloads/)
You need to add file to `static/db.go`

    package static
    
    var OrmConnectionString = "root:{{verysecret}}@tcp(127.0.0.1:3306)/game?charset=utf8"

Lambda uses Redis as memory menagement [Installation link](https://redis.io/download)

# How to install Lambda

Enter the go path directory

    cd $GOPATH/src

Clone the repository

    git clone https://github.com/canerakdas/lambda.git

Migrate database

    bee migrate -driver=mysql -conn="root:@tcp(127.0.0.1:3306)/game"

RUN!

    bee run
