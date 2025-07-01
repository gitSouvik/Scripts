#!/bin/bash

# Usage: marka <filename> — copies template.cpp into <filename>.cpp in the current directory

if [ $# -ne 1 ]; then
  echo "Usage: marka <filename>"
  exit 1
fi

TEMPLATE="$HOME/Scripts/template.cpp"
TARGET="$1.cpp"

cp "$TEMPLATE" "$TARGET" && echo "Done!"
