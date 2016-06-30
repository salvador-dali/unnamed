### Tools used

 - [Go](https://golang.org/doc/install) Latest: 1.6.2 darwin/amd64
 - [PostgreSQL](https://www.postgresql.org/download/) 9.4.0 (update to latest after machine switch)
 - python 2.7.9 ([psycopg2](http://initd.org/psycopg/)). Used to setup a database and run tests 
 
### Install tools with [Homebrew](http://brew.sh)

    brew update && brew upgrade
    brew install go postgresql
 
### ENV variables and GoPath

Set up the following env variables (all DB variables are related to your psql database). Typical values are shown:

    export PROJ_DB_NAME=postgres
    export PROJ_DB_USER=$USER
    export PROJ_DB_HOST=localhost
    export PROJ_DB_PWD=
    export PROJ_DB_PORT=5432
    export PROJ_HTTP_PORT=8080
    export PROJ_SECRET=asd4q-ass21sflse41r123hsz
    export PROJ_JWT_EXP_DAYS=2
    export PROJ_SALT_LEN_BYTE=16
    export PROJ_MAILGUN_DOMAIN=sandbox4d69a15edfe64dfaa3680f1a19fa50fa.mailgun.org
    export PROJ_MAILGUN_PRIVATE=
    export PROJ_MAILGUN_PUBLIC=pubkey-703cb3b4f8ee0807862184954ff7c252
    
By default after psql installation your password is empty. In this project it is not possible to have
empty env variables, so you have to change it `ALTER USER "user_name" WITH PASSWORD 'new_password';`
    
Cd to this repo and set up GOPATH equal to $PWD (on macOS you also need to set GOBIN):

    export GOPATH=$PWD
    if [[ `uname` == 'Darwin' ]]; then
      mkdir -p bin 
      export GOBIN=$GOPATH/bin
    fi

Run `go get` to install all Go dependencies.
