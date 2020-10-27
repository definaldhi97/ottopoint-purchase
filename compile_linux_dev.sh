
env GOOS=linux GOARCH=amd64 go build -o ottopoint-purchase

# scp -i ~/.ssh/LightsailDefaultKey-ap-southeast-1-new.pem ottopoint-purchase ubuntu@13.228.25.85:/home/ubuntu
scp ottopoint-purchase ichsan@10.10.43.64:/home/ichsan
