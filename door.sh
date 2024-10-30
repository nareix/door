#!/bin/sh

sudo ip netns add ns1
sudo ./setusb.sh

while true; do
  sudo ip netns exec ns1 ./door
  sleep 1
done
