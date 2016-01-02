include Makefile.defs

all: deps
	@go build -o $(NAME)

clean:
	@rm -f $(NAME)

deps:
	@go get -v
