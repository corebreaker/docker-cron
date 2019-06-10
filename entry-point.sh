#! /bin/sh

cd $(dirname $0)

if [ -z "$PROJECT_NAME" ]; then
    echo "No project name defined; the variable PROJECT_NAME must have a value"

    exit 1
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
    "start") start;;
    "shell") sh "$@";;
    "") echo "Available commandes: start, shell";;
    *)
        echo "Unknown command: $1"
        exit 1;;
esac
