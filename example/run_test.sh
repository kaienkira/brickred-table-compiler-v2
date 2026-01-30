#!/bin/bash

set -o pipefail

script_name=`basename "$0"`
script_abs_name=`readlink -f "$0"`
script_path=`dirname "$script_abs_name"`

cd "$script_path"/../compiler && make build
if [ $? -ne 0 ]; then exit 1; fi

# create test dir
test_dir="$script_path"/test
mkdir -p "$test_dir"
if [ $? -ne 0 ]; then exit 1; fi
cd "$test_dir"
if [ $? -ne 0 ]; then exit 1; fi

# copy files
cp "$script_path"/../compiler/bin/brickred-table-compiler .
if [ $? -ne 0 ]; then exit 1; fi
cp "$script_path"/table.xml .
if [ $? -ne 0 ]; then exit 1; fi
cp "$script_path"/main.cc .
if [ $? -ne 0 ]; then exit 1; fi

copy_table() {
    cp "$script_path"/copy.csv .
    if [ $? -ne 0 ]; then exit 1; fi
    cp "$script_path"/effect.csv .
    if [ $? -ne 0 ]; then exit 1; fi
    cp "$script_path"/item.csv .
    if [ $? -ne 0 ]; then exit 1; fi
    cp "$script_path"/matchmaking.csv .
    if [ $? -ne 0 ]; then exit 1; fi
    cp "$script_path"/npc.csv .
    if [ $? -ne 0 ]; then exit 1; fi
    cp "$script_path"/skill_level.csv .
    if [ $? -ne 0 ]; then exit 1; fi
}

# cpp test
copy_table
./brickred-table-compiler -f table.xml -l cpp -r server
if [ $? -ne 0 ]; then exit 1; fi
g++ -I "$script_path"/../cpp/src \
    -o "cpp_test" \
    main.cc \
    resource_item.cc \
    tbl_copy.cc \
    tbl_item.cc \
    tbl_matchmaking.cc \
    tbl_npc.cc \
    tbl_skill_level.cc \
    "$script_path"/../cpp/src/brickred/table/column_spliter.cc \
    "$script_path"/../cpp/src/brickred/table/line_reader.cc \
    "$script_path"/../cpp/src/brickred/table/util.cc
if [ $? -ne 0 ]; then exit 1; fi
