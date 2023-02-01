# Build Instructions

These instructions assume that `minikube` is installed.

```
$ minikube start

// Required to have minikube access locally built Docker images without needing
// to push them to something like Docker Hub.
$ eval $(minikube -p minikube docker-env)

$ docker build --no-cache -f Dockerfile --target build_reader -t go-reader .

$ docker build --no-cache -f Dockerfile --target build_writer -t go-writer .

$ kubectl apply -f deployment.yml

$ kubectl get pods -A

$ kubectl logs go-client-server --all-containers

// Should see something like:
2023/02/01 20:38:24 accepted conn, reading
2023/02/01 20:38:24 read msg from conn: hello

2023/02/01 20:38:24 writing

// Stop the client/server pod.
$ kubectl delete -f deployment.yml
```