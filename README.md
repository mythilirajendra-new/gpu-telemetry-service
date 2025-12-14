[![Coverage](https://sonarcloud.io/)](https://sonarcloud.io/)

# GPU Telemetry Service (GTS)
Elastic GPU Telemetry Pipeline with Message Queue


## Dev Setup
### Get repo
```shell
$ cd $GOPATH/src/github.com
$ git clone https://github.com/mythilirajendra-new/gpu-telemetry-service.git
$ cd gpu-telemetry-service
```

### Download dependencies
```shell
$ make vendor
```

### Install oapi-codegen
```shell
$ go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.3
$ sudo cp ~/go/bin/oapi-codegen /usr/local/bin/oapi-codegen
$ sudo chmod +x /usr/local/bin/oapi-codegen
```

### Generate server stub from openAPI spec
```shell
$ make oapi-codegen
```

### GO format source
```shell
$ make fmt
```

### Build and Run gpu-telemetry-service from source
```shell
$ make build-local
$ bin/gpu-telemetry-service
```

### Check gpu-telemetry-service
#### status
```shell
$ curl -i http://localhost:8181/status
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 05 Jun 2022 15:27:17 GMT
Content-Length: 15

{"status":"OK"}
```
NOTE: status API (unauthorized) is used for liveliness check

### Run test cases
```shell
$ make test
$ make coverage
```
