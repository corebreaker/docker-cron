#! /bin/sh

cd $(dirname $0)

if [ -z $PROJECT_NAME ]; then
    if [ -z $PROJECT_SOURCE_PATH ]; then
        echo "No project name defined; the variable PROJECT_NAME or PROJECT_SOURCE_PATH must have a value"

        exit 1
    fi

    PROJECT_NAME=`basename $PROJECT_SOURCE_PATH`
fi

PROJECT_DIR="/projects/$PROJECT_NAME"
if [ ! -d "$PROJECT_DIR" ]; then
    echo "The directory $PROJECT_DIR does not exist; had you mounted the docker volume ?"

    exit 1
fi

if [ -z $CRON_FILE ]; then
    CRON_FILE=./docker-cron.yml
fi

. /etc/docker-config.env

start() {
    cd $PROJECT_DIR

    $CRONBIN $CRON_FILE
}

case "$1" in
    "") echo "Available commandes: start, shell";;
    "start") start;;
    "shell")
        shift
        sh "$@";;
    *)
        echo "Unknown command: $1"
        exit 1;;
esac
