build: build-protos build-php-sdk

build-protos:
	protoc --proto_path=protos --go_opt=paths=source_relative --go_out=plugins=grpc:protos protos/*.proto


build-php-sdk:
	cd sdks && bash build-php.sh