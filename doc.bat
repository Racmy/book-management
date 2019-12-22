for /f "usebackq" %x in (`docker ps -aq`) do docker stop %x
for /f "usebackq" %x in (`docker ps -aq`) do docker rm %x
for /f "usebackq" %x in (`docker images -q`) do docker rmi %x
