GO_SRC = conf.go files.go main.go buildinfo.go

# In non-debug builds remove symbol information (-s) and
# DWARF debug information (-w).
GO_STRIP = -ldflags="-s -w"
ifneq ($(DEBUG),)
	GO_STRIP =
endif

.PHONY: clean all buildinfo.go
all: file2gobyte nullserv
file2gobyte: file2gobyte.c
	$(CC) -O3 -Wall $< -o $@
buildinfo.go:
	@./mkbuildinfo.sh
nullserv: buildinfo.go
	go build -o $@ $(GO_STRIP) $(GO_SRC)
clean:
	@rm -f buildinfo.go file2gobyte nullserv
