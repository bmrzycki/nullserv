GO_SRC := $(wildcard *.go) version.go

# In non-debug builds remove symbol information (-s) and
# DWARF debug information (-w).
GO_STRIP := -ldflags="-s -w"
ifneq ($(DEBUG),)
	GO_STRIP =
endif

.PHONY: clean all
all: file2gobyte nullserv
file2gobyte: file2gobyte.c
	$(CC) -O3 -Wall $< -o $@
version.go:
	@./go_ver.sh
nullserv: $(GO_SRC)
	go build -o $@ $(GO_STRIP) $?
clean:
	@rm -f file2gobyte nullserv version.go
