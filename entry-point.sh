#! /bin/sh

cd `dirname $0`

if [ -z "$PROJECT_NAME" ]; then
    echo "No project name defined; the variable PROJECT_NAME must have a value"

    exit 1
fi

PROJECT_DIR="/projects/$PROJECT_NAME"
if [ ! -d "$PROJECT_DIR" ]; then
    echo "The directory $PROJECT_DIR does not exist; had you mounted the docker volume ?"

    exit 1
fi

sleep 10000
