module infini-gateway

go 1.15

require (
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/nixgnehc/infini-framework v0.0.3
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
)


replace (
	"github.com/nixgnehc/infini-framework" => ../infini-framework
)