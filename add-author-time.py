#!/usr/bin/env python3

import datetime
import sys
import os

def prepend_cpp_header(filename):
    cwd = os.getcwd() 
    file_name = filename;
    filename = os.path.join(cwd, filename) 

    if not os.path.exists(filename):
        print(f"Error: {filename} not found.")
        sys.exit(1)

    with open(filename, 'r') as f:
        original_code = f.read()

    header = f"""/**
 *    author:  Calypsoo  
 *    created: {datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")}
**/"""

    with open(filename, 'w') as f:
        f.write(header + '\n\n' + original_code)

    print(f"Tagged {file_name}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: tag <filename>")
        sys.exit(1)

    prepend_cpp_header(sys.argv[1])