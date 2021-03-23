
env GOOS=linux GOARCH=amd64 go build -o ottopoint-purchase

# scp ottopoint-purchase rohmet@34.101.119.111:/home/build

scp ottopoint-purchase ichsan@10.10.43.64:/home/ichsan
