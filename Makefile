LIBBPF     := $(PWD)/vc5
LIBBPF_LIB := $(PWD)/vc5/bpf

export CGO_CFLAGS = -I$(LIBBPF)
export CGO_LDFLAGS = -L$(LIBBPF_LIB)

balancer: balancer.go vc5/kernel/bpf/bpf.o
	go build balancer.go

vc5/kernel/bpf/bpf.o: vc5
	cd vc5 && $(MAKE) kernel/bpf/bpf.o

vc5:
	git clone --branch v0.1.14 git@github.com:davidcoles/vc5.git

clean:
	rm -f balancer

distclean: clean
	rm -rf vc5
