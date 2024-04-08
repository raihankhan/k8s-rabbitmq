# k8s-rabbitmq

RabbitMQ official github repo:  

Go client Library: https://github.com/rabbitmq/amqp091-go

Version Supports: https://www.rabbitmq.com/versions.html

Server Doc: https://www.rabbitmq.com/admin-guide.html

YT: 
- [RabbitMQ on Kubernetes Deep Dive - David Ansari, VMware](https://www.youtube.com/watch?v=GxdyQSUEj5U&t=1357s)
- [RabbitMQ on Kubernetes for beginners - The Devops Guy](https://www.youtube.com/watch?v=_lpDfMkxccc)

### Publish and Consume messages

Build and push publisher image to your own registry:
```bash
make build REGISTRY=<your docker registry>
make push  REGISTRY=<your docker registry>
```
or, use image pushed to my dockerhub public repository - `raihankhanraka/rabbitmq-publisher:v1.0.0`. Image Tag might vary at the time of your deployment.

You can also push the image to your local [Kind](https://kind.sigs.k8s.io/) cluster very easily.
```bash
make push-to-kind
```

Finally, deploy the publisher pod. The following command with create a k8s Deployment with your provided ENV variables.

```bash
make deploy-publisher
```

Expose the service to your localhost with a specified port

```bash
$ kubectl port-forward svc/rabbitmq-publisher 9090:80
Forwarding from 127.0.0.1:9090 -> 80
Forwarding from [::1]:9090 -> 80
```

Publish data to rabbitmq:

```bash
curl -X POST --user 'guest:guest' "http://localhost:9090/publish/hello"
```