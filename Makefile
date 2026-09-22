.PHONY: build-server build-client network run-server run-client clean

network:
docker network create rede-projeto || true

build-server:
docker build -t server-app -f Docker/Server/Dockerfile .

build-client:	
docker build -t client-app -f Docker/Client/Dockerfile .

build: build-server build-client

run-server: network
docker run -d --name servidor --network rede-projeto -p 6742:6742 server-app

run-client: network
docker run -it --rm --name cliente --network rede-projeto client-app

clean:
docker stop servidor || true
docker rm servidor || true
docker rm cliente || true
docker network rm rede-projeto || true