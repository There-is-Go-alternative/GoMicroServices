module github.com/There-is-Go-alternative/GoMicroServices/chats

go 1.16

require github.com/google/uuid v1.3.0

require (
	firebase.google.com/go v3.13.0+incompatible
	github.com/There-is-Go-alternative/GoMicroServices/account v0.0.0-20211004145423-166650af87b6
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/api v0.58.0
)

require github.com/gin-gonic/gin v1.7.4

require (
	cloud.google.com/go/storage v1.18.0 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/rs/zerolog v1.25.0
	golang.org/x/net v0.0.0-20211011170408-caeb26a5c8c0 // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/sys v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20211012143446-e1d23e1da178 // indirect
	google.golang.org/grpc v1.41.0 // indirect
)
