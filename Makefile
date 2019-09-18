protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-kuma .

docker-run:
	(docker rm -f meshery-kuma) || true
	docker run --name meshery-kuma -d \
	-p 10000:10000 \
	-e DEBUG=true \
	layer5/meshery-kuma

run:
	DEBUG=true go run main.go