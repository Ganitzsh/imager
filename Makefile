NAME	=	12fact

all:
	go build -o ${NAME}

with-proto:
	make -C delivery/rpcv1/proto gen-proto
	go build -o ${NAME}

test:
	go test ./...

clean:
	make -C delivery/rpcv1/proto clean
	rm -f ${NAME}
