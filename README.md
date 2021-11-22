## sidecar proxy

### \# **single** container with webapp and proxy
```
cd singleFile/

# local

PORT=8080 PPORT=5050 go run singleFile/webAndProxy.go
curl localhost:5050

# docker

docker build -t local/sidecar:0.0.1 .
docker run -it -p 5050:5050 -e PORT=8080 -e PPORT=5050 local/sidecar:0.0.1
curl localhost:5050
```

### \# **separate** containers with webapp and proxy

#### **Run locally**
```
APPNAME=webapp make docker-build
APPNAME=proxy make docker-build

> run 'webapp' first then 'proxy' !

APPNAME=webapp make docker-run-webapp

APPNAME=proxy \
AIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' webapp) \
make docker-run-proxy

curl localhost:5050
```

#### **Run in K8s**
```
$ kubectl apply -f proxy-to-app.yml

$ kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
http-app-676f7f6b94-9zgfn   2/2     Running   0          29s

$ kubectl get svc
NAME           TYPE           CLUSTER-IP   EXTERNAL-IP      PORT(S)           AGE
http-app-svc   LoadBalancer   10.0.50.8    52.157.231.128   60000:31857/TCP   40s

$ curl -i 52.157.231.128:60000
HTTP/1.1 200 OK
Content-Length: 14
Content-Type: text/plain; charset=utf-8
Date: Mon, 22 Nov 2021 08:59:35 GMT
X-Content-Type-Options: nosniff

Cloud runner!

$ kubectl logs http-app-676f7f6b94-9zgfn http-sidecar
(...)
proxy 2021/11/22 08:59:35 proxy.go:52: request dump: &{GET / HTTP/1.1 1 1 map[Accept:[*/*] User-Agent:[curl/7.58.0]] {} <nil> 0 [] false 52.157.231.128:60000 map[] map[] <nil> map[] 10.244.1.1:42041 / <nil> <nil> <nil> 0xc00005ad00}
```
