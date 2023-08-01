
# Golang RESTful Apis
Project struct build in monolithic architecture but open for microservices architecture. The main functions:
- Open Api for authentication (login/register user)
- Open Api for User (Get/update user info)

## Technical and libraries
- Golang 1.19
- Go-Gin framework
- Gorm & Postgresql database
- Docker, docker-compose and integrated Github CI
- Some popular libraries (viper config, go-migrate, jwt-go, govalidator, logrus, samber-lodash, ...).

## How to use
Clone this repository to your local machine.
1. Ensure that you have Go environment
2. Postgresql database running (update the information config file app.yml)

For creating docker Postgresql, from root project run script:
```
docker-compose -f docker/docker-compose.local.yml up -d
```
For init tables, from root project run migration script:
```
ENV=local make migrate-up
```
### Run from source code
Go to root project and run:
```
cd services/question/ && go run main.go
```
### Run by docker
From root project run the script:
```
docker build -t sample-restful-api:lastest --build-arg APPCONFIG=app.yml -f ./services/question/Dockerfile .
docker run -d -p 3000:3000 sample-restful-api:lastest
```
## Sample request

Login/register api request:
```
curl --location 'http://localhost:3000/v1/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "nguyencongphuong@gmail.com",
    "password": "12345678"
}'
```

Response:

```json
{
    "data": {
        "id": "2d31eff9-b9c1-4422-b1b0-003acd2f58cb",
        "name": "",
        "email": "nguyencongphuong@gmail.com",
        "avatar": "",
        "access_token": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmQzMWVmZjktYjljMS00NDIyLWIxYjAtMDAzYWNkMmY1OGNiIiwiZW1haWwiOiJuZ3V5ZW5jb25ncGh1b25nQGdtYWlsLmNvbSIsImV4cCI6MTY5MzQ2MTAyMSwiaXNzIjoiU21hcnRLaWQtRWR1Y2F0aW9uIn0.YRDUxV5eapw-WqMuce9Vbp7EnOy_5N2faycV02n1Mac",
            "expired_at": "2023-08-31T12:50:21.978723+07:00"
        }
    },
    "error": null,
    "message": "Success",
    "success": true
}
```
## Questions / Feedbacks / Bugs
Feel free to reach out to me if you have any questions or feedback on how my code can be improved.
My email: nguyencongphuong@gmail.com

## TODO
- [x] REST APIs
- [x] Docker build
- [ ] Unit test
- [ ] Swagger documentation