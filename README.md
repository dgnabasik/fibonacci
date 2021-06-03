# fibonacci
fibonacci (docker)
Author: David Gnabasik
Date:   June 2, 2021.

Bonus points:
    √ Go tests.
    √ Performance data.
    Use dockertest.
    Include a Makefile.

Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.
The maximum possible value returned is math.MaxFloat64 = 1.798e+308 // 2**1023 * (2**53 - 1) / 2**52
But cutoff happens at math.MaxFloat32 = 3.4e+38  // 2**127 * (2**24 - 1) / 2**23

Non-Docker Program Installation Requirements::
 (a) O/S: Ubuntu Linux 18.04.
 (b) Golang: v1.16.4    Install from https://golang.org/doc/install 
 (c) Postgres v12.6+    Install from https://www.postgresql.org/download/
 (d) mkdir ~/myprogams && cd ~/myprogams && git clone https://github.com/dgnabasik/fibonacci  -OR- go get github.com/dgnabasik/fibonacci/...  (gets all dependencies)
 (e) Run tests from a terminal prompt with: cd ~/myprogams/dgnabasik/fibonacci && go test -v fibonacci
 (f) Run the docker container with: <<<<
 (g) Browse to http://localhost/fibonacci/ to interact with the web page.

The web API should expose operations to::
 (a) fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144): 
 (b) fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120), 
 (c) clear the data store. 

Imported Go Packages::
 (a) go get github.com/jackc/pgx/v4/pgxpool
 (b) go get github.com/gin-contrib/cors
 (c) go get github.com/gin-gonic/contrib/static
 (d) go get github.com/gin-gonic/gin


Postgres Database Tables::
DROP TABLE IF EXISTS public.fibonacci;
CREATE TABLE public.fibonacci (
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

