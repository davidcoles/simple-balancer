# A Simple Layer 4 Load-Balancer

This is a simple example of building a Layer 4 load-balancer with the
[VC5](https://github.com/davidcoles/vc5) library.

As the Go module requires an embedded object file (the compiled
XDP/eBPF code) currently it can't be simply imported, so I have added a
"replace" clause in the go.mod file to use a copy of the repository
which is automatically checked out by the Makefile. There's probably
a better way of doing this, and I'm open to suggestions!

The balancer in this example does no health checking of the services
on the backends. Make sure you have added the VIP to them or it won't
work!

There's also no statistics, little logging or any status indication,
so take a look at the more full-featured balancer code in main VC5
repository.

If you're interested I can add more complete examples - let me know!

## Backend servers

Set up two or more servers running a webserver on port 80. Add the VIP
that you're going to use to the loopback device so that the balanced
traffic is handled correctly.

`ip a add 192.168.100.1/32 dev lo`

Check that this works locally on each backend server:

`curl http://192.168.100.1/`


## Client machine

This server will be used to send traffic to the load-balancer. You can
optionally run a routing daemon such as BIRD or Quagga; configure it to
accept sessions/prefixes from the load-balancer's IP address.

## Load-Balancer

On the machine that you wish to use as the load-balancer, install the
necessary development environment to compile the code. Eg., on Ubuntu
20.04:
	
`apt-get install git build-essential libelf-dev clang libc6-dev libc6-dev-i386 llvm golang-1.16 libyaml-perl libjson-perl ethtool`
	
`ln -s /usr/lib/go-1.16/bin/go /usr/local/bin/go`
  
Check out this repository and edit the balancer.go source code,
setting the variables at the top of the file appropriately. If you are
using BGP then remember to set the autonomous system number correcly
and list client machine (our your router) under BGP peers. Check the
name of your network interface and specify this in "interfaces".

Build the binary:

`make`

Start the load-balancer. We need to raise the limit on memory which can be locked so that it is not paged out to disk:

`ulimit -l unlimited; ./balancer`

## Client ...

Back on the client machine, if you set up a routing daemon then you
should see the virtual IP address advertised and showing up in your
routing table (`ip r`).

If you're not using BGP then you can add a static route to the VIP via the load-balancer's address, eg:

`ip a add 192.168.100.1/32 via 10.11.12.13`

Try sending a request to the VIP:

`curl http://192.168.100.1/`


If that succeeds then send it some more traffic with ApacheBench:

`ab -n 100000 -c 100 http://192.168.100.1/`

You'll get a few minutes of load-balancing goodness before the
automatic kill-switch stops handling code in the kernel, and then the
balancer exits shortly afterwards.

If you're feeling brave you can remove the killswitch and put a `for`
around the sleep stetement, but make sure that all the parameters
you've put in are correct; this operates at layer 2 and there are
opportunities to introduce loops if you're not careful.
