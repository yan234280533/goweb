export IP=$1
echo "IP $IP"

for ((i=0;i<100;i++)); do sleep 0.01; curl -s "http://$IP:8000/cpu?cpu=3&count=400000000" > /dev/null ; done
curl http://$IP:8000/metrics

echo "end"
