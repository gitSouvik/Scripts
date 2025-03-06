#!/bin/bash

# # Check if filename is provided
# if [ -z "$1" ]; then
#     echo "Usage: ./rrun.sh <filename> [test_number]"
#     exit 1
# fi

# A="$1"
# N="$2"

# # Compile the C++ file
# g++-14 -O2 -std=c++20 -DWOOF_ "$1.cpp" -o "$1"
# if [ $? -ne 0 ]; then
#     echo -e "\033[1;31m>>> Compilation failed.\033[0m\n"
#     exit 1
# fi

# # Create a temporary file to store the output and error
# temp_output=$(mktemp)
# temp_error=$(mktemp)

# # Run the executable and redirect the output & error to the temp file
# ./"$1" > "$temp_output" 2> "$temp_error"

# # Display the output in the terminal
# echo -e "-----------------------\nOutput:\n-----------------------"
# cat "$temp_output"

# # Check if there were any errors
# if [ -s "$temp_error" ]; then
#     # echo -e "-----------------------\n\033[1;31mDebug:\033[0m\n-----------------------"
#     echo -e "-----------------------\nDebug:\n-----------------------"
#     cat "$temp_error"
# fi

# # Remove the temporary output and error files
# rm "$temp_output" "$temp_error"

# Check if filename is provided
if [ -z "$1" ]; then
    echo "Usage: ./rrun.sh <filename> [test_number]"
    exit 1
fi

A="$1"
N="$2"

# Compile the C++ file
g++-14 -O2 -std=c++20 -DWOOF_ "$A.cpp" -o "$A"
if [ $? -ne 0 ]; then
    echo -e "\033[1;31m>>> Compilation failed.\033[0m\n"
    exit 1
fi

# If no test case is provided (default case when running `rrun A`)
if [ -z "$N" ]; then
    # Run normally with custom input
    temp_output=$(mktemp)
    temp_error=$(mktemp)

    # Run the executable and redirect the output & error to the temp file
    ./"$A" > "$temp_output" 2> "$temp_error"

    # Display the output in the terminal
    echo -e "-----------------------\nOutput:"
    cat "$temp_output"

    # Check if there were any errors
    if [ -s "$temp_error" ]; then
        echo -e "-----------------------\nDebug:"
        cat "$temp_error"
        echo -e "\n-----------------------"
    fi

    # Remove the temporary output and error files
    rm "$temp_output" "$temp_error"
    exit 0
fi

# If a specific test case number is provided, run that test case
if [ "$N" != "0" ]; then
    INPUT_FILE="${A}-${N}.in"
    OUTPUT_FILE="${A}-${N}-output.out"

    if [ ! -f "$INPUT_FILE" ]; then
        echo -e "\033[1;31m>>> Test file missing: $INPUT_FILE not found.\033[0m\n"
        exit 1
    fi

    # echo -e "Running (Test Case $N):"
    # echo -e "\033[38;5;82m>>> (Test Case $i):\033[0m"

    # Run the executable and store output
    ./"$A" < "$INPUT_FILE" > "$OUTPUT_FILE" 2> temp_debug.log

    cat "$INPUT_FILE"

    # Display the output
    echo -e "-----------------------\nOutput:"
    cat "$OUTPUT_FILE"

    # Display debug/error info if any
    if [ -s temp_debug.log ]; then
        echo -e "-----------------------\nDebug:"
        cat temp_debug.log
        echo -e "\n-----------------------"
    fi

    # Clean up temporary files
    rm "$OUTPUT_FILE" temp_debug.log
    exit 0
fi

i=1

# If N is 0, run all available test cases
if [ "$N" == "0" ]; then    
    for INPUT_FILE in ${A}-*.in; do
        # Extract test number using basename and sed
        TEST_NUMBER=$(basename "$INPUT_FILE" | sed -E 's/.*-([0-9]+)\.in/\1/')

        # Define output file for this test case
        OUTPUT_FILE="${A}-${TEST_NUMBER}-output.out"  # Store actual output separately

        # echo -e "Running (Test Case $N):"
        # echo -e "\033[38;5;82m>>> (Test Case $i):\033[0m"
        echo -e "Running "$1".in \033[38;5;82m(Test Case $i)\033[0m:"

        # Run the executable and store output
        ./"$A" < "$INPUT_FILE" > "$OUTPUT_FILE" 2> temp_debug.log

        # cat "$INPUT_FILE"

        # Display the output
        echo -e "-----------------------\nOutput:"
        cat "$OUTPUT_FILE"

        # Display debug/error info if any
        if [ -s temp_debug.log ]; then
            echo -e "-----------------------\nDebug:"
            cat temp_debug.log
            echo -e "\n-----------------------"
        fi

        # Clean up temporary files
        rm "$OUTPUT_FILE" temp_debug.log
        ((i++))
    done
    exit 0
fi

# # Run normally with custom input
# temp_output=$(mktemp)
# temp_error=$(mktemp)

# # Run the executable and redirect the output & error to the temp file
# ./"$1" > "$temp_output" 2> "$temp_error"

# # Display the output in the terminal
# echo -e "-----------------------\nOutput:"
# cat "$temp_output"

# # Check if there were any errors
# if [ -s "$temp_error" ]; then
#     echo -e "-----------------------\nDebug:"
#     cat "$temp_error"
#     # echo -e "\n-----------------------"
# fi

# # Remove the temporary output and error files
# rm "$temp_output" "$temp_error"
