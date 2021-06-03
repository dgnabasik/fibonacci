mkdir -p ~/github.com/dgnabasik && cd ~/github.com/dgnabasik 
git clone https://github.com/dgnabasik/fibonacci  
docker-compose up --build
echo "Open another terminal and execute: cd ~/github.com/dgnabasik/fibonacci && ./migration.sh "
echo "Execute the curl commands in README.md or open a browser to the web addresses in README.md."
echo "Press ctrl-C to stop the server. Execute ./cleanDocker.sh to remove the docker containers."
