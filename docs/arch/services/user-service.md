## User service - сервис для работы с пользователей

### Доступные методы

#### CreateUser - метод создания пользователя из auth-service
```
Request:
gRPC /api.UserService/CreateUser
{
    email: "<EMAIL>",
    password: "<PASSWORD>",
    nickname: "<NICKNAME>",
    avatarUrl: "<URL_TO_AVATAR>",
    createdAt: "<DATE>",
    updatedAt: "<DATE>",
    info: "<INFO>"
}
Response:
{
    id: "<ID>",
}
```

#### UpdateUser - метод обновления пользователя
```
Request:
PUT /api/user/:id
{
    info: "<INFO>",
    avatarUrl: "<URL_TO_AVATAR>"
}
Response:
{
    id: "<ID>",
}
```

#### GetUserByNikname - метод получения пользователя по никнейму
```
Request:
GET /api/user/:nickname
Response:
{
    id: "<ID>",
    email: "<EMAIL>",
    nickname: "<NICKNAME>",
    avatarUrl: "<URL_TO_AVATAR>",
    createdAt: "<DATE>",
    updatedAt: "<DATE>",
    info: "<INFO>"
}
```

#### Add Friend - метод добавления друга в друзья
```
Request:
POST /api/user/:id/friends/:friendId
Response:
{
    status: "OK"
}
```

#### Remove Friend - метод удаления друга из друзей
```
Request:
DELETE /api/user/:id/friends/:friendId
Response:
{
    status: "OK"
}
```

#### Get Friends - метод получения списка друзей пользователя
```
Request:
GET /api/user/:id/friends
Response:
{
    friendsList: [
        {
            id: "<ID>",
            nickname: "<NICKNAME>",
            avatarUrl: "<URL_TO_AVATAR>",
            approved: true|false
        }
    ]
}
```

#### Approve incoming Friend Request - метод подтверждения входящего запроса на дружбу
```
Request:
PUT /api/user/:id/friends/:friendId
Response:
{
    status: "OK"
}
```

#### Reject incoming Friend Request - метод отклонения входящего запроса на дружбу
```
Request:
DELETE 
```