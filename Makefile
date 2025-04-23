push: docker-build deploy

docker-build:
	docker exec rgbmatrix-build-arm64-1 make build
	chmod +x .bin/led

build:
	go build -tags with_cgo -o .bin/led cmd/main.go

deploy:
	ssh pi4 sudo systemctl stop led
	scp -C .bin/led pi4:/home/tnolle/led/led
	ssh pi4 sudo systemctl start led