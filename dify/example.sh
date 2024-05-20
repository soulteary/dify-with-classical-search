#!/bin/bash

curl 'http://localhost:7700/multi-search' \
  -H 'Authorization: Bearer soulteary' \
  -H 'Content-Type: application/json' \
  --data-raw '{"queries":[{"indexUid":"movies","q":"sky","limit":3,"offset":0}]}'