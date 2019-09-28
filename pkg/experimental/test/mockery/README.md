# [mockery](https://github.com/vektra/mockery)

## Installation

```bash
$ go get github.com/vektra/mockery/.../
$ mockery -version
1.0.0

$ mockery -name DB -dir ./pkg/test/mockery/service/ -output ./pkg/test/mockery/service/mocks
Generating mock for: DB in file: pkg/test/mockery/service/mocks/DB.go

### see how to use the generated file in pkg/test/mockery/service/greeter_test.go

```
