# msa-messenger

[Архитектура приложения](docs/arch/README.md)

## Docker
Общий Dockerfile находится в директории build/Dockerfile. Он на вход принимает один аргумент
```
MAIN_PATH
```
который указывает на путь к main файлу для модуля.

docker-compose.yml описывает запуск всех модулей, передавая значение аргумента и выставляю наружу свой порт. Все внутренние порты одинаковые 8080.

## Makefile

### Build all
```
make build
```

### Run all
```
make run-all
```

## Kubernetes
-1. Сборка образов
```
make build-all
```
0. Запуск и настройка minikube
```
➜  msa-messenger git:(main) ✗ minikuber start
➜  msa-messenger git:(main) ✗ eval $(minikube docker-env)
➜  msa-messenger git:(main) ✗ minikube image load msa-messanger-user:latest
➜  msa-messenger git:(main) ✗ minikube image load msa-messanger-auth:latest
➜  msa-messenger git:(main) ✗ minikube image load msa-messanger-messaging:latest
➜  msa-messenger git:(main) ✗ minikube image ls --format table
```
1. Создать namespace для приложения
```
➜  msa-messenger git:(main) ✗ kubectl apply -f deploy/kuber/namespace.yaml 
➜  msa-messenger git:(main) ✗ kubectl config set-context --current --namespace=msa-messagenger-namespace
```
2. Создать deployment для приложения
```
➜  msa-messenger git:(main) ✗ kubectl apply -f ./deploy/kuber/deployments.yaml

deployment.apps/msa-messanger-user-deployment created
deployment.apps/msa-messanger-auth-deployment created
deployment.apps/msa-messanger-messaging-deployment created

➜  msa-messenger git:(main) ✗ kubectl get pods
NAME                                                  READY   STATUS    RESTARTS   AGE
msa-messanger-auth-deployment-7f8dd95fc8-b9xqv        1/1     Running   0          3m56s
msa-messanger-messaging-deployment-5958d8fd66-f59pc   1/1     Running   0          3m56s
msa-messanger-user-deployment-574746c775-hdnj2        1/1     Running   0          3m56s
```
3. Создать service для приложения
```
➜  msa-messenger git:(main) ✗ kubectl apply -f deploy/kuber/services.yaml
service/msa-messanger-user-service created
service/msa-messanger-auth-service created
service/msa-messanger-messaging-service created

➜  msa-messenger git:(main) ✗ kubectl get svc
NAME                              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
msa-messanger-auth-service        ClusterIP   10.109.206.34    <none>        8080/TCP   7s
msa-messanger-messaging-service   ClusterIP   10.104.172.69    <none>        8080/TCP   7s
msa-messanger-user-service        ClusterIP   10.104.178.207   <none>        8080/TCP   7s

➜  msa-messenger git:(main) ✗ kubectl port-forward service/msa-messanger-user-service 8090:8080
Forwarding from 127.0.0.1:8090 -> 8080
Forwarding from [::1]:8090 -> 8080
Handling connection for 8090
Handling connection for 8090
Handling connection for 8090
```

4. Добавить Ingress
```
➜  msa-messenger git:(main) ✗ minikube addons enable ingress

➜  msa-messenger git:(main) ✗ kubectl get pods -n ingress-nginx
NAME                                        READY   STATUS      RESTARTS      AGE
ingress-nginx-admission-create-fg7kt        0/1     Completed   0             181d
ingress-nginx-admission-patch-7vckq         0/1     Completed   1             181d
ingress-nginx-controller-7799c6795f-5vbng   1/1     Running     6 (27m ago)   181d

➜  msa-messenger git:(main) ✗ kubectl apply -f deploy/kuber/ingress.yaml 
ingress.networking.k8s.io/msa-messenger-ingress configured
```

5. Добавить в /etc/host адрес minikube
```
➜  msa-messenger git:(main) ✗ minikube ip  
192.168.49.2
```

```
➜  msa-messenger git:(main) ✗ vi /etc/hosts
192.168.49.2    app.test.com
```

6. Выполнить запросы к сервисам
```
➜  msa-messenger git:(main) ✗ curl http://app.test.com/user/ready
{"status":"user module ready"}% 
➜  msa-messenger git:(main) ✗ curl http://app.test.com/auth/ready
{"status":"auth module ready"}
➜  msa-messenger git:(main) ✗ curl http://app.test.com/messaging/ready
{"status":"messaging module ready"}
```