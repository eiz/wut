set -e
export PATH=/Users/eiz/code/go_src/go/bin/go:$PATH
export GOOS=linux GOARCH=amd64
GOARCH=386 go tool asm -I ~/code/go_src/go/src/runtime wut.s
go tool asm -I ~/code/go_src/go/src/runtime wut64.s
go tool compile -+ -wb=false main.go
go tool pack c main.a main.o wut.o wut64.o
go tool link -f -v -multiboot -T 0x101000 -E _load_kernel main.a
