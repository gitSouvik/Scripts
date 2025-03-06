#!/bin/bash

# Check if there are enough arguments
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <filename1> <filename2> ... <file_extension>"
    exit 1
fi

# Get the last argument as the file extension
extension=${!#}

# Set all but the last argument as filenames
filenames=("${@:1:$(($# - 1))}")

# Create files with specified names
for name in "${filenames[@]}"; do
    touch "${name}.${extension}"
    echo "Created ${name}.${extension}"
done