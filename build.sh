echo "try kill process of server"
ps -aux|grep ./ocp-check-config|awk '{print $2}'|xargs -I {} kill -9 {}
echo "================================"

echo "check server process"
ps -aux|grep ./ocp-check-config
echo "================================"

echo "start build server: ocp-check-config"
rm -rf nohup.out
go build
nohup ./ocp-check-config &
echo "you can cat server log in nohup.out file"
cat nohup.out
echo "================================"

echo "start build client: easytool"
cd ./cmd/easytool/
go build
sudo mv easytool /usr/local/bin/easytool
easytool -h
easytool para get
echo "================================"
echo "you can use command: easytool -h for use details. "
