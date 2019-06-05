PROTO_DIR = .

PROTO_FILES = \
	${PROTO_DIR}/main.proto

all:
	docker run -v "$PWD:/tmp" -w "/tmp"

gen-proto:
	protoc  -I ${PROTO_DIR} --go_out=plugins=grpc:${PROTO_DIR} ${PROTO_FILES}

clean:
	rm -f *.pb.go
