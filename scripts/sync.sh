#!/bin/bash

# Array of hostnames
hostnames=(
#  "pi4.local:8080"
  "192.168.0.26:8085"
  "localhost:8085"
)

for hostname in "${hostnames[@]}"; do
  curl -X PUT http://$hostname/$1?name=$2
done
