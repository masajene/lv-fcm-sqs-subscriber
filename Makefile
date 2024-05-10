.PHONY: build

docker-build:
	aws ecr get-login-password --region ap-northeast-1 --profile ambits | docker login --username AWS --password-stdin 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com
	docker build --build-arg VERSION="v1" -t lv-push-wrapper-fcm-sub-repo .
	docker tag lv-push-wrapper-fcm-sub-repo:latest 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com/lv-push-wrapper-fcm-sub-repo:latest
	docker push 703932351856.dkr.ecr.ap-northeast-1.amazonaws.com/lv-push-wrapper-fcm-sub-repo:latest
