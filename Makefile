.PHONY: clean all
all: file2gobyte nullsrv
file2gobyte: file2gobyte.c
	$(CC) -O3 -Wall $< -o $@
nullsrv: $(wildcard *.go)
	go build -o nullsrv $(shell ./go_ver.sh) $?
clean:
	@rm -f file2gobyte nullsrv
