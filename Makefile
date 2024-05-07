.PHONY: build

docker-build-2:
	aws ecr get-login-password --region ap-northeast-1 --profile ambits | docker login --username AWS --password-stdin 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com
	docker build --build-arg VERSION="v1" -t dev-lv-push-wrapper-fcm-queue-subscriber .
	docker tag dev-lv-push-wrapper-fcm-queue-subscriber:latest 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com/dev-lv-push-wrapper-fcm-queue-subscriber:latest
	docker push 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com/dev-lv-push-wrapper-fcm-queue-subscriber:latest
