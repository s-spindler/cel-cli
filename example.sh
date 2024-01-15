#!/usr/bin/env bash

set -euo pipefail

j='{ "potato": { "temp": "cool" }}'

if bin/cel-cli-linux -i "$j" -e "has(i.potato) && has(i.potato.temp) && i.potato.temp == 'cool'"
then
    echo noice
else
    echo pffffbbbt
fi
