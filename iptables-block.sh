#!/bin/bash

#conn_count=$(sysctl --values net.netfilter.nf_conntrack_count)
#
#sysctl --write net.netfilter.nf_conntrack_max=${conn_count}
#
#sysctl --write net.netfilter.nf_conntrack_buckets=$((${conn_count}/4))

sudo iptables --flush

#sudo iptables --new-chain RATE-LIMIT

sudo iptables --append INPUT --protocol tcp --syn --jump RATE-LIMIT

sudo iptables --append RATE-LIMIT \
    --match hashlimit \
    --hashlimit-mode srcip \
    --hashlimit-upto 1/sec \
    --hashlimit-burst 5 \
    --hashlimit-name conn_rate_limit \
    --jump ACCEPT

sudo iptables --append RATE-LIMIT --jump DROP