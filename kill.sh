ps -ef | grep ./ottopoint-product | grep -v grep | awk '{print $2}' | xargs kill