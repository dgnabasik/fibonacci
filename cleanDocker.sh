#sudo /etc/init.d/couchdb stop
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -a -q)
docker volume rm $(docker volume ls -q)
# Do not delete {bridge, host, none}
docker network rm $(docker network ls -q | grep -v 6e41bf9d5fa3 | grep -v ab62e256644f | grep -v 64c52bcde6a0)
docker system prune -a
#sudo /etc/init.d/couchdb start
aplay /usr/share/sounds/speech-dispatcher/test.wav

