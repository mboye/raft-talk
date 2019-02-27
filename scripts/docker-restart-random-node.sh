docker ps | grep 'raft-talk_node' | awk '{print $1}' | sort -R | head -n1 | xargs docker restart
