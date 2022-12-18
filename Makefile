build:
	./deployments/deploy.sh

unit-test:
	go test -v -cover ./...

integration-test:
	docker build .
	docker stop remi_app || true && docker rm remi_app || true
	docker-compose up -d
	go test -v -cover ./features

expose-all:
	kubectl port-forward service/remi-service 8080:8080