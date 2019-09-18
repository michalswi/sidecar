# sidecar proxy POC

### single file
```sh
$ docker build -t local/sidecar:0.0.1 .
$ docker run -it -p 5050:5050 -e PORT=8080 -e PPORT=5050 local/sidecar:0.0.1
```

### sidecar
```sh
# test
$ PORT=8080 go run web.go
$ PPORT=5050 APORT=8080 go run proxy.go
$ curl localhost:5050

# K8s - TODO
$ kubectl apply -f webProxyDeploy.yml
```

