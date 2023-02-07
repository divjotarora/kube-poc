# Build Instructions

These instructions assume that `minikube` is installed.

```
$ minikube start

// Required to have minikube access locally built Docker images without needing
// to push them to something like Docker Hub.
$ eval $(minikube -p minikube docker-env)

// Included for completeness. I had issues installing Go in the Ubuntu container
// due to ARM weirdness, so we can manually build the Go binary locally and the
// Dockerfile will copy it over into the container. The pipes-client-linux-arm64
// is already included in the Git repo, so this command isn't strictly necessary
// to run the container.
$ GOOS=linux GOARCH=arm64 go build -o pipes-client-linux-arm64 ./main.go

$ docker build --no-cache -f Dockerfile -t named-pipes-poc .

// Verify that the images show up as expected.
$ docker images

$ kubectl apply -f deployment.yml

$ kubectl get pods -A

$ kubectl logs go-named-pipes --all-containers

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