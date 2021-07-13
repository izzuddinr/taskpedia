module taskpedia-worker

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)
