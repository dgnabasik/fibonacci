echo dgnabasik
#git pull https://github.com/dgnabasik/fibonacci
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.0"
git push -u origin main
