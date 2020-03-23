 kill -9 $(lsof -i TCP:8006 | grep LISTEN | awk '{print $2}')


