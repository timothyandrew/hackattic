#!/usr/bin/env bash

set -euo pipefail

# brew install pkcrack

TOKEN="$(cat ../token)"

curl https://www.gutenberg.org/cache/epub/50133/pg50133.txt > dunwich_horror.txt
zip -r plain.zip dunwich_horror.txt

ZIPURL="$(curl https://hackattic.com/challenges/brute_force_zip/problem?access_token=$TOKEN | jq -r '.zip_url')"
curl "$ZIPURL" > package.zip

pkcrack -C package.zip -c dunwich_horror.txt -P plain.zip -p dunwich_horror.txt -d out.zip -a
unzip -B out.zip

http POST "https://hackattic.com/challenges/brute_force_zip/solve\?access_token=$TOKEN" secret="$(cat secret.txt)"
