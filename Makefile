build:
	./deployments/deploy.sh

unit-test:
	go test -v -cover ./...

integration-test:
	docker build .
	docker stop remi_app || true && docker rm remi_app || true
	docker-compose up -d
	go test -v -cover ./features
