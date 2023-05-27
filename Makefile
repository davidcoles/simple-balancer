LIBBPF := $(PWD)/vc5

export CGO_CFLAGS  = -I$(LIBBPF)
export CGO_LDFLAGS = -L$(LIBBPF)/bpf

MAX_FLOWS ?= 100000

balancer: balancer.go vc5/kernel/bpf/bpf.o
	go build balancer.go

vc5/kernel/bpf/bpf.o: vc5
	cd vc5 && $(MAKE) kernel/bpf/bpf.o MAX_FLOWS=$(MAX_FLOWS)

vc5:
	git clone --branch v0.1.14 https://github.com/davidcoles/vc5.git

clean:
	rm -f balancer

distclean: clean
	rm -rf vc5
