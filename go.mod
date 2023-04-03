module github.com/jialunzhai/crimemap

go 1.20

replace (
	gopkg.in/jcmturner/rpc.v1 => ./vendor/gopkg.in/jcmturner/rpc.v1
	gopkg.in/jcmturner/gokrb5.v6 => ./vendor/gopkg.in/jcmturner/gokrb5.v6
)

require (
	github.com/trinodb/trino-go-client v0.310.0
	golang.org/x/sync v0.1.0
	google.golang.org/grpc v1.54.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
)
