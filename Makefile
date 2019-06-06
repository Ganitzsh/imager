NAME	=	12fact

all:
	make -C delivery/rpcv1/proto gen-proto
	go build -o ${NAME}

test:
	go test ./...

clean:
	make -C delivery/rpcv1/proto clean
	rm -f ${NAME}
