#!/bin/bash

# Ensure that a file name is provided
if [ -z "$1" ]; then
  echo "Usage: ./script.sh <file_name>"
  exit 1
fi

# Variables
file_name=$1
CPP_FILE="${file_name}.cpp"
EXECUTABLE="${file_name}"

# Compile the file
echo "[Debug Mode] Compiling $CPP_FILE with g++20."
g++-14 -O2 -std=c++20 -DWOOF_ -o "$EXECUTABLE" "$CPP_FILE"

# Check if compilation was successful
if [ $? -ne 0 ]; then
  echo -e "\033[1;31m   Compilation failed.\033[0m\n"
  exit 1
fi

# # Loop through all test cases (files named as file_name-*.in)
# for INPUT_FILE in ${file_name}-*.in; do
#   # Extract test number using basename and sed
#   TEST_NUMBER=$(basename "$INPUT_FILE" | sed -E 's/.*-([0-9]+)\.in/\1/')

#   # Define the expected output file and output file for this test case
#   EXPECTED_FILE="${file_name}-${TEST_NUMBER}.out"
#   OUTPUT_FILE="${file_name}-${TEST_NUMBER}-output.out"                                                                                  # ***** ADDED NEW *****

#   echo -e "Running ${INPUT_FILE} (Test Case $TEST_NUMBER):\n-------------------\nOutput:"

#   # Run the executable and measure time and memory
#   { /usr/bin/time -p ./"$EXECUTABLE" < "$INPUT_FILE" > "$OUTPUT_FILE" ; } 2> time_output.txt
#   # { /usr/bin/time -p ./"$EXECUTABLE" < "$INPUT_FILE" > "$OUTPUT_FILE" 2> "$ERROR_FILE"; } 2> time_output.txt
  
#   # Check if the program crashed                                                                                                       # ***** ADDED NEW *****        
#   if [ $? -ne 0 ]; then
#     echo -e "\033[1;31m>>> Test Case $TEST_NUMBER Failed (Segmentation Fault)!\033[0m\n"
#     continue
#   fi

#   ELAPSED_TIME=$(grep real time_output.txt | awk '{print $2}')  # Extract elapsed time
#   MEMORY_USAGE=$(ps -o rss= -p $$)  # Get memory usage

#   # Display output and expected output
#   cat "$OUTPUT_FILE"
#   echo -e "-------------------\nExpected:"
  
#   if [ -f "$EXPECTED_FILE" ]; then
#     cat "$EXPECTED_FILE"
#   else
#     echo "Expected output file ${EXPECTED_FILE} not found!"
#   fi
  
#   # Show memory and time
#   echo -e "-------------------\nMemory: $MEMORY_USAGE KB\nTime: ${ELAPSED_TIME} s"



# cout << debug file starts here !
# Loop through all test cases (files named as file_name-*.in)
for INPUT_FILE in ${file_name}-*.in; do
  # Extract test number using basename and sed
  TEST_NUMBER=$(basename "$INPUT_FILE" | sed -E 's/.*-([0-9]+)\.in/\1/')

  # Define output files
  EXPECTED_FILE="${file_name}-${TEST_NUMBER}.out"
  OUTPUT_FILE="${file_name}-${TEST_NUMBER}-output.out"
  DEBUG_FILE="${file_name}-${TEST_NUMBER}-debug.txt"  # stderr → Debug

  echo -e "Running ${INPUT_FILE} (Test Case $TEST_NUMBER):\n-------------------\nOutput:"

  # Run the executable, capturing stdout and stderr separately
  ./"$EXECUTABLE" < "$INPUT_FILE" > "$OUTPUT_FILE" 2> "$DEBUG_FILE"

  # Check if the program crashed                                                                                                       # ***** ADDED NEW *****        
  if [ $? -ne 0 ]; then
    echo -e "\033[1;31m>>> Test Case $TEST_NUMBER Failed (Segmentation Fault)!\033[0m\n"
    continue
  fi

  # Display output and expected output
  cat "$OUTPUT_FILE"
  echo -e "-------------------\nExpected:"

  if [ -f "$EXPECTED_FILE" ]; then
    cat "$EXPECTED_FILE"
  else
    echo "Expected output file ${EXPECTED_FILE} not found!"
  fi

  # Display debug output
  if [ -s "$DEBUG_FILE" ]; then
    echo -e "-------------------\nDebug:"
    cat "$DEBUG_FILE"
  fi
  echo -e "\n-------------------"
# cout << debug file ends here !



  # Compare output with expected output if expected file exists
  if [ -f "$EXPECTED_FILE" ]; then
    DIFF=$(diff -q "$OUTPUT_FILE" "$EXPECTED_FILE")

    if [ -z "$DIFF" ]; then
      echo -e "\033[38;5;82m>>> Test Case $TEST_NUMBER Passed!\033[0m\n"
    else
      echo -e "\033[1;31m>>> Test Case $TEST_NUMBER Failed!\033[0m\n"
    fi
  else
    echo -e "\033[1;31m>>> Test Case $TEST_NUMBER Failed (Missing expected output)!\033[0m\n"
  fi
done

# # Clean up
# rm time_output.txt