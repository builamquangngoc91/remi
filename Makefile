build:
	docker build -t remi-app .
	docker tag remi-app quangngoc430/remi:v1.0.2
	docker push quangngoc430/remi:v1.0.2
	kubectl rollout restart deployment remi

unit-test:
	go test -v -cover ./...

integration-test:
	docker build .
	docker stop remi_app || true && docker rm remi_app || true
	docker-compose up -d
	go test -v -cover ./features
	