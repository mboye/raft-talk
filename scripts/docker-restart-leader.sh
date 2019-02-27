ETCD_ENDPOINTS=http://localhost:2001,http://localhost:2002,http://localhost:2003
leader_no="$(etcdctl --endpoints="$ETCD_ENDPOINTS" endpoint status -w table | grep true | grep -E -o '200(.)' | cut -c4-)"
echo "Restarting node ${leader_no}"
docker restart "raft-talk_node-${leader_no}_1"
