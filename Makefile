
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
version.go:
	@./go_ver.sh
nullsrv: version.go files.go main.go
	go build -o nullsrv $(GO_STRIP) $?
clean:
	@rm -f file2gobyte nullsrv version.go
