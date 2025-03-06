#!/bin/bash

# Check if there are enough arguments
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <number_of_files> <filename1> <filename2> ... <file_extension>"
    exit 1
fi

# Get the number of files from the first argument
n=$1
shift # Remove the first argument (number of files) from the list

# Get the last argument as the file extension
extension=${!#}
set -- "${@:1:$(($# - 1))}" # Set all but the last argument as filenames

# Create files with specified names
for name in "$@"; do
    touch "${name}.${extension}"
done

echo "$n files created with the specified names and the extension '$extension'"