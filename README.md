# Raft Consensus Algorithm and Kubernetes

Thanks for attending the talk at the Distributed Systems Meetup on February 27, 2019.

It may take more than 300 ms for you to try out the demos below :-)

-- Magnus

# Etcd demo

You can find the Etcd demo under [etcd-demo](etcd-demo).
Simply run `docker-compose up` to start the Etcd cluster with 3 servers.
Proceed by running the demo application under [etcd-demo/app](etcd-demo/app) using the `run.sh` script.

```sh
$ ./run.sh
+ export ETCD_ENDPOINTS=http://localhost:2001,http://localhost:2002,http://localhost:2003
+ ETCD_ENDPOINTS=http://localhost:2001,http://localhost:2002,http://localhost:2003
+ go run app.go
ERRO[0000] background color configuration not found in etcd
INFO[0000] current background color                      background_color="#ffffff"
INFO[0000] starting server on port 8080
INFO[0000] watching for background color changes in etcd
```

The demo application serves a webpage at [http://localhost:8080/](http://localhost:8080/).
The background color configuration parameter can be changed in Etcd using the [change-background-color.sh](etcd-demo/app/change-background-color.sh) script.

```
$ ./change-background-color.sh "#ccc"
+ etcdctl --endpoints=http://localhost:2001,http://localhost:2002,http://localhost:2003 put app-config/background-color '#ccc'
OK
```

The demo application listens for configuration changes and reacts immediately.

```
INFO[0287] background color changed                      new_value="#ccc" old_value=
```

# Leader elector demo

This demo consists of a Kubernetes [deployment manifest](leader-elector-demo/deployment.yaml) and a [demo application](leader-elector-demo/app/app.go).
The deployment manifest has two containers: [leader-elector](https://hub.docker.com/r/googlecontainer/leader-elector) and the demo application.

The podspec part of the deployment looks like this:

```yaml
spec:
  serviceAccountName: "operator"
  containers:
    - name: "leader-elector"
      image: googlecontainer/leader-elector:0.5
      args:
        - "--election=mboye-leader-election-demo"
        - "--http=localhost:4040"
        - "--election-namespace=mboye"
      ports:
        - containerPort: 4040
          protocol: TCP
    - name: "app"
      image: mboye/leader-election-app:v1
      ports:
        - containerPort: 8080
          protocol: TCP
```

The leader-elector container exposes leadership information on port 4040.
The demo application serves a webpage on port 8080, that informs whether the HTTP request was served by the leader pod or not.
