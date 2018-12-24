# Build And Development
test:
	@go test -v -cover -covermode=atomic ./...

vendor-prepare:
	@go get -u github.com/golang/dep/cmd/dep

vendor:
	@dep ensure -v

engine:
	@go build -o engine app/main.go

run: 
	@docker-compose up -d

stop:
	@docker-compose down

docker:
	@docker build . -t bxcodec/tweetor:latest

docker-mock:
	@docker build . -t bxcodec/mock-tweetor -f DockerfileMockTweet

# Deployments
docker-push:
	@docker push bxcodec/tweetor

docker-push-mock:
	@docker push bxcodec/mock-tweetor

deploy-ingress:
	@kubectl apply -f app/deployments/ingress.yaml

deploy: docker docker-push
	@kubectl apply -f app/deployments/deployment.yaml

.PHONY: test engine vendor-prepare vendor run stop docker docker-push deploy-ingress deploy