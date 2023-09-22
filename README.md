# reakgo
Simple Framework to quickly build webapps in GoLang

# DevTool

## Overview

The **DevTool** is a convenient utility for generating Go source files, specifically model and controller files, in your project. This tool simplifies the process of creating these files, allowing you to focus on your code logic.

## Prerequisites

Before using the DevTool, make sure you have the following set up:

1. Go Programming Language installed on your system.
2. Execution permissions for the `filegenerator.sh` script.

## Usage

To use the DevTool, follow these steps:

1. Give execution permissions to the `filegenerator.sh` script by running the following command in your terminal:

    ```bash
    chmod +x filegenerator.sh
    ```

2. Run the `filegenerator.sh` script with the following command:

    ```bash
    ./filegenerator.sh <packagename> <filename>
    ```

    - `<packagename>`: The name of the package where you want to generate the file. You have two options:
      - Specify `models` to generate the file in the `models` package.
      - Specify `controllers` to generate the file in the `controllers` package.
      - Specify `both` to generate the file in both the `models` and `controllers` packages.

    - `<filename>`: The name of the Go source file you want to generate(dont add extension og the file).

3. After running the command, the DevTool will create the specified Go source file with the provided filename in the selected package(s).

## Examples

### Example 1: Generate a Model File

To generate a model file named `User` in the `models` package, run:

```bash
./filegenerator.sh models User
```