push: docker-build deploy

docker-build:
	docker exec rgbmatrix-build-arm64-1 make build
	chmod +x .bin/led

build:
	go build -tags with_cgo -o .bin/led cmd/main.go

deploy:
	ssh piz sudo systemctl stop led
	scp -C .bin/led tnolle@192.168.0.26:/home/tnolle/led/led
	ssh piz sudo systemctl start led