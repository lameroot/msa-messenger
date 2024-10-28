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
➜  minikube ssh -- docker rmi -f msa-messenger-db:latest
➜  msa-messenger git:(main) ✗ minikube image load msa-messenger-user:latest
➜  msa-messenger git:(main) ✗ minikube image load msa-messenger-auth:latest
➜  msa-messenger git:(main) ✗ minikube image load msa-messenger-messaging:latest
➜  msa-messenger git:(main) ✗ minikube image ls --format table
```
1. Создать namespace для приложения
```
➜  msa-messenger git:(main) ✗ kubectl apply -f deploy/kuber/namespace.yaml 
➜  msa-messenger git:(main) ✗ kubectl config set-context --current --namespace=msa-messagenger-namespace
```
2. Создаем конфиги, store и deployment для БД
```
kubectl apply -f deploy/kuber/db/db-config.yaml
kubectl apply -f deploy/kuber/db/db-deployment.yaml 
kubectl apply -f deploy/kuber/db/db-service.yaml 

➜ kubectl get pods
NAME                                           READY   STATUS    RESTARTS   AGE
msa-messenger-db-deployment-5798b7f448-h2kll   1/1     Running   0          40s
➜ kubectl exec -it msa-messenger-db-deployment-5798b7f448-52pjx  -- psql -h localhost -U messenger --password  -p 5432 messenger 
Password: messenger
psql (12.4 (Debian 12.4-1.pgdg100+1))
Type "help" for help.

messenger=# select * from users;
 id | email | password_hash | nickname | created_at | updated_at | avatar_url | info 
----+-------+---------------+----------+------------+------------+------------+------
(0 rows)

messenger=# 
```
2. Создать config, deployment, service для приложения
```
➜ kubectl apply -f deploy/kuber/config.yaml 
➜ kubectl apply -f deploy/kuber/deployments.yaml

deployment.apps/msa-messenger-user-deployment created
deployment.apps/msa-messenger-auth-deployment created
deployment.apps/msa-messenger-messaging-deployment created

➜ kubectl get pods                                                
NAME                                                  READY   STATUS    RESTARTS   AGE
msa-messenger-auth-deployment-65cb98986f-2x4hc        1/1     Running   0          15m
msa-messenger-db-deployment-5798b7f448-h2kll          1/1     Running   0          20m
msa-messenger-messaging-deployment-75fcdb6d7c-dstmh   1/1     Running   0          45s
msa-messenger-user-deployment-78c9b6b4c6-kndv6        1/1     Running   0          45s
```
3. Создать service для приложения
```
➜ kubectl apply -f deploy/kuber/services.yaml
service/msa-messenger-user-service created
service/msa-messenger-auth-service created
service/msa-messenger-messaging-service created

➜ kubectl get svc
NAME                              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
msa-messenger-auth-service        ClusterIP   10.109.206.34    <none>        8080/TCP   7s
msa-messenger-messaging-service   ClusterIP   10.104.172.69    <none>        8080/TCP   7s
msa-messenger-user-service        ClusterIP   10.104.178.207   <none>        8080/TCP   7s
```

Пример (выполнять не требуется, показывает как пробросить порты)
```
➜ kubectl port-forward service/msa-messenger-user-service 8090:8080
Forwarding from 127.0.0.1:8090 -> 8080
Forwarding from [::1]:8090 -> 8080
Handling connection for 8090
Handling connection for 8090
Handling connection for 8090
```


4. Добавить Ingress
```
➜ minikube addons enable ingress

➜ kubectl get pods -n ingress-nginx
NAME                                        READY   STATUS      RESTARTS      AGE
ingress-nginx-admission-create-fg7kt        0/1     Completed   0             181d
ingress-nginx-admission-patch-7vckq         0/1     Completed   1             181d
ingress-nginx-controller-7799c6795f-5vbng   1/1     Running     6 (27m ago)   181d

➜ kubectl apply -f deploy/kuber/ingress.yaml 
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

7. Регистрация пользователя test1
```
curl --location --request POST 'http://app.test.com/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email":"test1@asd.ru",
    "password":"1234567",
    "nickname" : "test1"
}'

{"user":{"id":"4de84508-c0b3-4a6f-83b1-59241aa05fbd","email":"test1@asd.ru","nickname":"test1","avatar_url":"","info":""}}
```

Регистрация пользователя test2
```
curl --location --request POST 'http://app.test.com/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email":"test2@asd.ru",
    "password":"1234567",
    "nickname" : "test2"
}'
{"user":{"id":"29130e05-adf2-4d3f-844c-23026a6d52de","email":"test2@asd.ru","nickname":"test2","avatar_url":"","info":""}}
```

8. Login пользователя test1
```
curl --location --request POST 'http://app.test.com/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email" : "test1@asd.ru",
    "password" : "1234567"
}'
{"user":{"id":"4de84508-c0b3-4a6f-83b1-59241aa05fbd","email":"test1@asd.ru","nickname":"test1","avatar_url":"","info":""},"token":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3MzA0MDY4NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.FSgpVefKHNFTElTZAHXgU0GFD0MuGxdhOAfUNt55ZS4"}}
```

9. Добавить в друзья test1 test2 (id = 29130e05-adf2-4d3f-844c-23026a6d52de)
```
curl --location --request POST 'http://app.test.com/user/friends' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM' \
--header 'Content-Type: application/json' \
--data-raw '{
    "friend_id" : "29130e05-adf2-4d3f-844c-23026a6d52de"
}'

{"friendship_status":"pending"}
```

10. Список друзей test1 (id = 4de84508-c0b3-4a6f-83b1-59241aa05fbd)
```
curl --location --request GET 'http://app.test.com/user/friends' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM'

[{"friend_id":"29130e05-adf2-4d3f-844c-23026a6d52de","friendship_status":"pending"}]
```

11. Отправить сообщение от test1 к test2 (id = 29130e05-adf2-4d3f-844c-23026a6d52de)
```
curl --location --request POST 'http://app.test.com/messaging/send' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM' \
--header 'Content-Type: application/json' \
--data-raw '{
    "friend_id" : "29130e05-adf2-4d3f-844c-23026a6d52de",
    "message" : "this is test message1",
    "sent_time" : 1729328392
}'

{"success":true,"delivery_time":1729802370}

curl --location --request POST 'http://app.test.com/messaging/send' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM' \
--header 'Content-Type: application/json' \
--data-raw '{
    "friend_id" : "29130e05-adf2-4d3f-844c-23026a6d52de",
    "message" : "this is test message2",
    "sent_time" : 1729328392
}'

{"success":true,"delivery_time":1729802401}
```

12. Получить все сообщения от пользователя test2 (29130e05-adf2-4d3f-844c-23026a6d52de)
```
curl --location --request GET 'http://app.test.com/messaging/messages' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGRlODQ1MDgtYzBiMy00YTZmLTgzYjEtNTkyNDFhYTA1ZmJkIiwiZW1haWwiOiJ0ZXN0MUBhc2QucnUiLCJzdWIiOiI0ZGU4NDUwOC1jMGIzLTRhNmYtODNiMS01OTI0MWFhMDVmYmQiLCJleHAiOjE3Mjk4MDI5NjcsIm5iZiI6MTcyOTgwMjA2NywiaWF0IjoxNzI5ODAyMDY3fQ.qGBLhbPahmQH3Bhgy_tPLUjAkM8FvVCI0hvLuunGTLM' \
--header 'Content-Type: application/json' \
--data-raw '{
    "friend_id" : "4c467512-1a2c-4921-bfaa-ebe6b56cb4b0",
    "count_messages" : 10
}'

[
  {
    "sender_id": "4de84508-c0b3-4a6f-83b1-59241aa05fbd",
    "receiver_id": "29130e05-adf2-4d3f-844c-23026a6d52de",
    "message": "this is test message2",
    "read": false,
    "sent_time": 1729802401
  },
  {
    "sender_id": "4de84508-c0b3-4a6f-83b1-59241aa05fbd",
    "receiver_id": "29130e05-adf2-4d3f-844c-23026a6d52de",
    "message": "this is test message1",
    "read": false,
    "sent_time": 1729802370
  }
]
```

## Testing
Mockery
```
go install github.com/vektra/mockery/v2@latest
```

### Generate mocks
```
go generate ./...
```

### Тесты
1. Важно указать в начале теста t.Parallel()
2. Создаем type args struct {} , которая будет описывать переменные для запроса метода, который тестируем
3. tests := []struct{
    name string //название теста
    args args //аргументы
    want *тип ответа //какой ответ будет возвращаться 
    assertErr assert.ErrorAssertionFunc //какая будет ошибка возвращаться
    mock func(t *testing.T) //вернуть объект с мокированными сущностями
}
### Запускаем тесты

Minimock
```
go install github.com/gojuno/minimock/v3/cmd/minimock@latest
```
