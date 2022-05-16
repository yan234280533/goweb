export IP=$1
echo "IP $IP"

for ((j=0;;j++))
do
  echo "loop: $j," `date "+%Y-%m-%d %H:%M:%S"`
  curl http://$IP:8080/restart
  time for ((i=0;i<200;i++)); do sleep 0.4; curl -s "http://$IP:8000/cpu?cpu=1&count=400000000" > /dev/null ; done
  curl http://$IP:8000/metrics| grep http_response_time
  curl http://$IP:8080/stop
  echo "end loop $j"
done