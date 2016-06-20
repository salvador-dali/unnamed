### Tech used

 - [Go](https://golang.org/doc/install) Latest: 1.6.2 darwin/amd64
 - [PostgreSQL](https://www.postgresql.org/download/) 9.4.0 (update to latest after machine switch)
 - python 2.7.9 ([psycopg2](http://initd.org/psycopg/)). Used to setup a database and run tests 
 
### ENV variables and GoPath

Set up the following env variables (all DB variables are related to your psql database)

    export PROJ_DB_NAME=
    export PROJ_DB_USER=
    export PROJ_DB_HOST=
    export PROJ_DB_PWD=
    export PROJ_DB_PORT=
    export PROJ_HTTP_PORT=8080
    export PROJ_SECRET=
    export PROJ_JWT_EXP_DAYS=
    export PROJ_SALT_LEN_BYTE=
    
Set up GOPATH equal to a working directory of this repo.
    
    export GOPATH=

Run `go get` to install all Go dependencies.
