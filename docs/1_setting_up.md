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
    export PROJ_SECRET=secret
    export PROJ_JWT_EXP_DAYS=365
    export PROJ_SALT_LEN_BYTE=16
    
Cd to this repo and set up GOPATH equal to $PWD (on macOS you also need to set GOBIN):

    export GOPATH=$PWD
    if [[ `uname` == 'Darwin' ]]; then
      mkdir -p bin 
      export GOBIN=$GOPATH/bin
    fi

Run `go get` to install all Go dependencies.
