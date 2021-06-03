# fibonacci
fibonacci (docker)
Author: David Gnabasik
Date:   June 2, 2021.
Task:   Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.

Semi-Automated Installation:: See ./install_fibonacci.sh
    mkdir -p ~/github.com/dgnabasik && cd ~/github.com/dgnabasik 
    git clone https://github.com/dgnabasik/fibonacci 
    cd fibonacci && pwd
    docker-compose up --build
    "Do not forget to first disable any local instances of Postgres with: sudo systemctl stop postgresql@12-main.service"
    "Open another terminal and execute: cd ~/github.com/dgnabasik/fibonacci && ./migration.sh "
    "Execute the curl commands in README.md or open a browser to the web addresses in README.md."
    "Press ctrl-C to stop the server. Execute ./cleanDocker.sh to remove the docker containers."

Manual Client Installation::
 (a) mkdir ~/github.com/dgnabasik && cd ~/github.com/dgnabasik && git clone https://github.com/dgnabasik/fibonacci 
 (b) Run tests from a terminal prompt with: cd ~/github.com/dgnabasik/fibonacci && cat fib.env && source fib.env && export NODE_ENV=development && go test -v 
     if not running, start the local Postgres server with: sudo systemctl start postgresql@12-main.service
 (c) Run the docker container with: docker-compose up --build
 (d) Wait for the message of: "LOG:  database system is ready to accept connections."
 (d) Open a browser to enter API URLs or use the curl command.

The web API should expose operations to::
 (a) fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144): 
 (b) fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120). No, there are 11.
 (c) clear the data store. 

Non-Docker Development URL Examples::
 (a) http://localhost:5000/fib/clear     ==> ClearDataStore: true
 (b) http://localhost:5000/fib/10        ==> Fibonacci: 55
 (c) http://localhost:5000/fib/upper/120 ==> NumberMemoizedResults: 11

Docker URL Examples:: Or simply browse to the web address.
 (a) curl http://localhost:8080/fib/20          ==> Fibonacci: 6765
 (b) curl http://localhost:8080/fib/upper/500   ==> NumberMemoizedResults: 14
 (c) curl http://localhost:8080/fib/clear       ==> ClearDataStore: true

Ubuntu 18.04 Development Environment::
 (a) Docker (v20.10.6) & docker-compose (v1.26):   Install from https://docs.docker.com/engine/install/ubuntu/
 (b) Golang: v1.16.4    Install from https://golang.org/doc/install 
 (c) Postgres v12.6+    Install from https://www.postgresql.org/download/
 (d) golang-migrate     Install from https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
     migrate create -ext sql -dir migrations -seq create_fibonacci_table
 (e) Execute: docker-compose up --build
 (f) Execute: cd ~/github.com/dgnabasik/fibonacci && sudo ./migration.sh

Imported Go Packages::
 (a) go get github.com/jackc/pgx/v4/pgxpool
 (b) go get github.com/gin-contrib/cors
 (c) go get github.com/gin-gonic/gin
 (d) github.com/jackc/pgx/v4
 (e) github.com/jackc/pgx/v4/pgxpool

Environment Variables in fib.env:: 
 (a) NODE_ENV: production or development. Change NODE_ENV to 'production' to run in docker container.
 (b) Postgres datatbase variables within Docker container. Change <<PWD>> in FIB_DATABASE_URL to run local Postgres server.

Program Limitations::
 (a) math.MaxFloat64 = 1.798e+308 // 2**1023 * (2**53 - 1) / 2**52
 (b) math.MaxFloat32 = 3.4e+38  // 2**127 * (2**24 - 1) / 2**23

Postgres Database Tables:: See migrations/000001_create_fibonacci_table.up.sql
DROP TABLE IF EXISTS fibonacci;
CREATE TABLE IF NOT EXISTS fibonacci (
    id integer NOT NULL DEFAULT 0,
    fibvalue numeric(308,0),
   	CONSTRAINT fibonacci_pkey PRIMARY KEY (id)
)
WITH (OIDS=FALSE) TABLESPACE pg_default;

Git Commands::
echo "# fibonacci" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/dgnabasik/fibonacci.git
git push -u origin main
