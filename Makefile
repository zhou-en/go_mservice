install:
	which swager || GO11MODULE=off go get -u github.com/go-sweager/go-sweager/cmd/sweager

swagger:
	GO11MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
