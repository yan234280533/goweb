export IP=$1
echo "IP $IP"

curl http://$IP:8080/restart

for ((i=0;i<100;i++)); do sleep 0.01; curl -s "http://$IP:8000/cpu?cpu=3&count=400000000" > /dev/null ; done

curl http://$IP:8000/metrics| grep http_response_time

curl http://$IP:8080/stop

echo "end"
