# fibonacci
fibonacci (docker)
Author: David Gnabasik
Date:   June 2, 2021.
Task:   Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.

Automated Installation:: 
 (a) ./install_fibonacci.sh

Manual Client Installation::
 (a) mkdir ~/github.com && cd ~/github.com && git clone https://github.com/dgnabasik/fibonacci  -OR- go get github.com/dgnabasik/fibonacci/...  (gets all dependencies)
 (b) Run tests from a terminal prompt with: cd ~/github.com/dgnabasik/fibonacci && go test -v 
 (c) Run the docker container with: docker-compose up --build
 (d) Wait for the message of: "LOG:  database system is ready to accept connections."
 (d) Open a browser to enter API URLs.

The web API should expose operations to::
 (a) fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144): 
 (b) fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120). No, there are 11.
 (c) clear the data store. 

Non-Docker Development URL Examples::
 (a) http://localhost:5000/fib/clear     ==> true
 (b) http://localhost:5000/fib/10        ==> 55
 (c) http://localhost:5000/fib/upper/120 ==> 11 (not 12!)

Docker URL Examples::
    curl http://localhost:8080/fib/upper/120

 (a) http://localhost:5000/fib/clear     ==> true
 (b) http://localhost:5000/fib/10        ==> 55
 (c) http://localhost:5000/fib/upper/120 ==> 11 (not 12!)

Ubuntu 18.04 Development Environemnt::
 (a) Docker (v20.10.6) & docker-compose (v1.26):   Install from https://docs.docker.com/engine/install/ubuntu/
     docker-compose up --build
 (b) Golang: v1.16.4    Install from https://golang.org/doc/install 
 (c) Postgres v12.6+    Install from https://www.postgresql.org/download/
 (h) golang-migrate     Install from https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
    migrate create -ext sql -dir migrations -seq create_items_table

Imported Go Packages::
 (a) go get github.com/jackc/pgx/v4/pgxpool
 (b) go get github.com/gin-contrib/cors
 (c) go get github.com/gin-gonic/contrib/static
 (d) go get github.com/gin-gonic/gin

Environment Variables in fib.env:: 
 (a) NODE_ENV: production or development.
 (b) Postgres datatbase variables within Docker container.
 
Program Limitations::
 (a) math.MaxFloat64 = 1.798e+308 // 2**1023 * (2**53 - 1) / 2**52
 (b) math.MaxFloat32 = 3.4e+38  // 2**127 * (2**24 - 1) / 2**23

Postgres Database Tables:: See migrations/000001_create_items_table.up.sql
DROP TABLE IF EXISTS public.fibonacci;
CREATE TABLE IF NOT EXISTS public.fibonacci (
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

