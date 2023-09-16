#!/bin/bash

# `air` is aliased so this lets the script find it
shopt -s expand_aliases

tailwindcss -i css/input.css -o css/output.css --watch &
air