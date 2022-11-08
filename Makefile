build:
	docker build -t remi-app .
	docker tag remi-app quangngoc430/remi:v1.0.2
	docker push quangngoc430/remi:v1.0.2
	kubectl rollout restart deployment remi