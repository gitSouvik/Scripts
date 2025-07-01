#!/bin/bash

if [ $# -eq 1 ]; then
  file="$1"
  seed_arg=()  # No seed passed → gen.cpp handles random seed
elif [ $# -eq 2 ]; then
  seed_arg=("$1")
  file="$2"
else
  echo "Usage: gen <file> OR gen <seed> <file>"
  exit 1
fi

g++-14 -O2 -std=c++20 gen.cpp -o gen || { echo "gen.cpp failed to compile"; exit 1; }

# Run and capture input
input=$(./gen "${seed_arg[@]}")  # ← THIS is the correct line
echo -e "-----------------------\nInput:"
printf "%s\n" "$input"

# Run and show output
printf "%s\n" "$input" | /Users/woofwoof/Scripts/runsamples-debug.sh "$file"