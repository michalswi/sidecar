## sidecar proxy


### \# single file with webapp and proxy
```
# local

PORT=8080 PPORT=5050 go run singleFile/webAndProxy.go
curl localhost:5050


# docker

docker build -t local/sidecar:0.0.1 .
docker run -it -p 5050:5050 -e PORT=8080 -e PPORT=5050 local/sidecar:0.0.1
curl localhost:5050
```

### \# separate webapp and proxy

```
# local

PORT=8080 go run web.go
PPORT=5050 APORT=8080 go run proxy.go
curl localhost:5050

# K8s - TODO
kubectl apply -f webProxyDeploy.yml
```

