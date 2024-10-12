## Auth service - сервис авторизация/аутентификации и регистрации

### Доступные методы

#### Register - метод регистрации пользователя с передачей email/password
Request:
```
POST /api/auth/register
{
    "email": "<EMAIL>",
    "password": "<PASSWORD>"
}
Response:
{
    "token": "<KEY>",
    
}
```


#### Login - метод входа в систему с передачей email/password

#### CheckAuth - проверка авторизации пользователя
