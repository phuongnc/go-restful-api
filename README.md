
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
Clone this repository to your local machine. Ensure that you have Go environment and Postgresql database running (update the information config file app.yml)

### Run from source code
 Go to root project and run:
```
go run services/question/main.go
```
### Run by docker
From root project run the script:
```
docker build -t backend-service:lastest --build-arg APPCONFIG=app.yml -f ./services/question/Dockerfile .
docker run -d -p 3000:3000 backend-service:lastest
```

## Questions / Feedbacks / Bugs
Feel free to reach out to me if you have any questions or feedback on how my code can be improved.
My email: nguyencongphuong@gmail.com

## TODO
- [x] REST APIs
- [x] Docker build
- [ ] Unit test
- [ ] Swagger documentation