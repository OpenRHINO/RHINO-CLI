#!/bin/sh
set -o errexit
set -o nounset
set -o pipefail

# Look for executable files named $FUNC_NAME and check uniqueness
file_path=$(find ./ -type f -name "$FUNC_NAME" -executable)
if [ "$file_path" ]; then
    if [ "$(echo "$file_path" | wc -l)" -gt 1 ]; then
        echo "Found multiple executable files named '$FUNC_NAME'. Please check your Makefile!" >&2
        exit 1 
    fi
    mv "$file_path" "/app/$FUNC_NAME"
    echo "Loading app $FUNC_NAME"
else
# Exit and report an err when no $FUNC_NAME file is found
    echo "Cannot find file $FUNC_NAME!" >&2
    exit 1
fi

if [ ! -d "/shared_lib" ]; then
    mkdir "/shared_lib"
fi
cd "/shared_lib"
echo "The shared_lib dir created"

# Identify which libs need to be loaded
sharedlibs=$(ldd "/app/$FUNC_NAME" | grep -vE "ld-musl-x86_64|mpi" | awk '{print $3}' || true)
if [ "$sharedlibs" != "" ]; then
    echo "Shared libs found"
    echo "$sharedlibs" > path.txt
else
    echo "No shared lib need to be loaded"
    exit 0
fi

# Seek for the share libs
while read -r line
do
    if [ "$line" = "" ]; then continue; fi
    islink=$(ls -l "$line" | grep ' -> ' || true)
    if [ "$islink" != "" ]
    then
        dir=$(dirname "$line")
        base=$(basename "$line")
        file=$(ls -l "$line" | awk '{print $NF}')
        cp -p "${dir}/${file}" "/shared_lib/$base"
        echo "cp ${dir}/${file} /shared_lib/$base"
    else
        cp -p "$line" "/shared_lib/"
        echo "cp $line /shared_lib/"
    fi
done < path.txt
rm path.txt
echo "ldd analysis success!"