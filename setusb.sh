#!/bin/sh

ip link set enxe04e7a9631bb netns ns1
ip netns exec ns1 ip addr add 192.168.11.192/16 dev enxe04e7a9631bb
ip netns exec ns1 ip link set enxe04e7a9631bb up
killall door
