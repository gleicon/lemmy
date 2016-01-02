include src/Makefile.defs
# set this to the proper PATH on macosx + homebrew.
PKG_CONFIG_PATH=/usr/local/Cellar/sqlite/3.8.8.3/lib/pkgconfig/

all: server

deps:
	make -C src deps

server:
	make -C src
	@mkdir -p $(BINARY_DIR)
	@cp src/$(NAME) $(BINARY_DIR)/$(NAME)

clean:
	make -C src clean
	@rm -f $(BINARY_DIR)/$(NAME)

run:
	make -C src run

test:
	make -C src test

