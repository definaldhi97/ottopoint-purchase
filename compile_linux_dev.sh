
env GOOS=linux GOARCH=amd64 go build -o ottopoint-purchase

scp -i ~/.ssh/LightsailDefaultKey-ap-southeast-1-new.pem ottopoint-purchase ubuntu@13.228.25.85:/home/ubuntu
