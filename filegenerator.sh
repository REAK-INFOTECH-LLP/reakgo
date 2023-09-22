#!/bin/bash

# Check if two arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: ./generate.sh <packagename> <filename>"
    exit 1
fi

# Assign the arguments to variables
packagename="$1"
filename="$2"

# Compile the generator program
cd generator
go build
cd ..

# Run the generator executable with the provided package name and filename as arguments
./generator/generator "$packagename" "$filename"
