NAME	=	12fact

all:
	make -C proto gen-proto
	go build -o ${NAME}

test:
	go test ./...

clean:
	make -C proto clean
	rm -f ${NAME}
