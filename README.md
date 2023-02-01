# Build Instructions

These instructions assume that `minikube` is installed.

```
$ docker build --no-cache -f reader-dockerfile -t divjot/go-reader .

$ docker build --no-cache -f reader-dockerfile -t divjot/go-reader .

// Required to have minikube access locally built Docker images without needing
// to push them to something like Docker Hub.
$ eval $(minikube -p minikube docker-env)

$ minikube start

$ kubectl apply -f deployment.yml

$ kubectl get pods -A
```