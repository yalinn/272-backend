# Probee API
this project uses a [golang](https://go.dev/) framework, [gofiber](https://gofiber.io/)

Before you start:
  - Be sure you've installed [go programming language](https://go.dev/) to your device
  - Be sure you've installed [swagger framework](https://github.com/swaggo/swag) (You need to install even we're not using it yet)
  - Be sure you have a MongoDB cluster connection URL & Redis connection URL

example .env file:
```env
MONGO_URI="mongodb://localhost:27017/"
REDIS_URI="redis://127.0.0.1:6379"
# JWT settings:
JWT_SECRET_KEY="SUPER_SECRET_KEY"
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
PORT=":8080"
IMAP_S_HOST="-student-imap-server-domain-"
IMAP_T_HOST="-academic-imap-server-domain-"
IMAP_PORT=993
```
