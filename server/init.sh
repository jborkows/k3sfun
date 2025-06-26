#!/usr/bin/env bash

bash traefik/init.sh
bash metal/init.sh
bash certmanager/init.sh
bash pihole/init.sh
