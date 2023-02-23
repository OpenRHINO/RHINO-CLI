#!/bin/sh
# set -o errexit
# set +o nounset
# set -o pipefail

file_path=$(find ./ -name $FUNC_NAME)
if [ "$file_path" ]; then
    mv $file_path /app/$FUNC_NAME
    echo "loding app $FUNC_NAME"
else
    echo "cannot find file $FUNC_NAME!"
    exit 1
fi

if [ ! -d "/shared_lib" ]; then
    mkdir /shared_lib
fi
cd /shared_lib
echo "dir created"

ldd /app/$FUNC_NAME | grep -vE "ld-musl-x86_64|mpi" | awk '{print $3}' > path.txt
echo "shared lib found"

while read -r line
do  
    islink=$(ls -l $line | grep -e " -> ")
    if [[ "$islink" != "" ]]
    then
        dir=$(dirname $line)
        base=$(basename $line)
        file=$(ls -l $line | awk '{print $NF}')
        cp ${dir}"/"${file} /shared_lib/$base
        echo cp ${dir}"/"${file} /shared_lib/$base
    else
        cp $line /shared_lib/
        echo cp $line /shared_lib/
    fi
done < path.txt
rm path.txt
echo "ldd analysis success!"      