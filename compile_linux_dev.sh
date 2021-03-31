
env GOOS=linux GOARCH=amd64 go build -o ottopoint-purchase

scp ottopoint-purchase rohmet@34.101.175.164:/home/build
