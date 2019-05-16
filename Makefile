
# In non-debug builds remove symbol information (-s) and
# DWARF debug information (-w).
GO_STRIP := -ldflags="-s -w"
ifneq ($(DEBUG),)
	GO_STRIP =
endif

.PHONY: clean all
all: file2gobyte nullsrv
file2gobyte: file2gobyte.c
	$(CC) -O3 -Wall $< -o $@
nullsrv: $(wildcard *.go)
	go build -o nullsrv $(shell ./go_ver.sh) $(GO_STRIP) $?
clean:
	@rm -f file2gobyte nullsrv
