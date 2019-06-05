PROTO_DIR = .

PROTO_FILES = \
	${PROTO_DIR}/12fact.proto

all:
	docker run -v "$PWD:/tmp" -w "/tmp"

gen-proto:
	protoc  -I ${PROTO_DIR} --go_out=plugins=grpc:${PROTO_DIR} ${PROTO_FILES}
