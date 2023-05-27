# A Simple Layer 4 Load-Balancer

## Backend servers

Set up two or more servers running a webserver on port 80.

## Client machine

This server will be used to send traffic to the load-balancer. You can
optionally run a routing daemon such as BIRD or Quagga; configure to
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

Build the binary: `make`

Start the load-balancer: `./balancer`

## Client ...

Back on the client machine, if you set up a routing daemon then you
should see the virtual IP address advertised and showing up in your
routing table (`ip r`).

If you're not using BGP then you can add a static route to the VIP via the load-balancer's address, eg:

`ip a add 192.168.101.1/32 via 10.11.12.13`

Try sending a request to the VIP:

`curl http://192.168.101.1/`


If that succeeds then send it some more traffic with ApacheBench:

`ab -n 100000 -c 100 http://192.168.101.1/`
