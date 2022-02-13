check_install:
	which swagger || GO11MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO11MODULE=off swagger generate spec -o ./swagger.yaml --scan-models