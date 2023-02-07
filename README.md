# Build Instructions

These instructions assume that `minikube` is installed.

```
$ minikube start

// Required to have minikube access locally built Docker images without needing
// to push them to something like Docker Hub.
$ eval $(minikube -p minikube docker-env)

$ docker build --no-cache -f Dockerfile --target build_client -t pipes-client .

$ docker build --no-cache -f Dockerfile --target build_server -t pipes-server .

// Verify that the images show up as expected.
$ docker images

$ kubectl apply -f deployment.yml

$ kubectl get pods -A

$ kubectl logs go-client-server --all-containers

// Stop the client/server pod.
$ kubectl delete -f deployment.yml
```

The `kubectl logs...` command should output something like this:
```
opening file
making aggregate request
file is opened
waiting to write
write succesful
waiting to write
write succesful
waiting to write
write succesful
waiting to write
write succesful
waiting to write
write succesful
cursor returned
ready to read
{"a": {"$numberInt":"1"}}
{"a": {"$numberInt":"1"}}
{"a": {"$numberInt":"1"}}
{"a": {"$numberInt":"1"}}
{"a": {"$numberInt":"1"}}
{"t":{"$date":"2023-02-04T06:34:07.639Z"},"s":"I",  "c":"CONTROL",  "id":5760901, "ctx":"main","msg":"Applied --setParameter options","attr":{"serverParameters":{"enableComputeMode":{"default":false,"value":true}}}}
```