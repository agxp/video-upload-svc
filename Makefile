SHELL=/bin/bash
dev: build-local deploy-local
production: build deploy

build:
	rm -f ./bin/*
	protoc --proto_path=${GOPATH}/src --micro_out=. --go_out=. -I. proto/upload.proto
	go get
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/video_upload -installsuffix cgo .
	docker build -t agxp/router .

build-local:
	rm -f ./bin/*
	protoc --proto_path=${GOPATH}/src --micro_out=. --go_out=. -I. proto/upload.proto
	go get
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/video_upload -installsuffix cgo .
	@eval $$(minikube docker-env) ;\
	docker build -t video_upload .

run:
	docker run --net="host" \
		-p 50051 \
		-e DB_HOST=localhost \
		-e DB_PASS=password \
		-e DB_USER=postgres \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		-e MINIO_URL=minio-0 \
		-e MINIO_ACCESS_KEY=minio \
		-e MINIO_SECRET_KEY=minio123 \
		video-upload

deploy:
	docker push agxp/video-upload-svc
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yaml
	kubectl apply -f ./deployments/deployment.yaml

deploy-local:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yaml
	kubectl apply -f ./deployments/deployment.yaml