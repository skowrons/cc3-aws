#!/bin/bash

init () {
    echo "initializing copilot app with the following name: $1"
    copilot app init $1 || exit 1
    echo ""

    echo "+--------------------+"
    echo "| deploy api service |"
    echo "+--------------------+"
    copilot init -n api -t "Load Balanced Web Service" -d ./api/Dockerfile --deploy 
    echo ""

    echo "+--------------------+"
    echo "| deploy cat service |"
    echo "+--------------------+"
    copilot init -n cat -t "Backend Service" -d ./cat/Dockerfile --deploy
    echo ""

    echo "everything is set up - showing services"
    copilot svc ls

    echo "you can reach the api under the following url:" 
    copilot svc show -a $1 -n api | grep http | sed 's/^.*http/http/'
}

public () {
    if [ -z "$var" ]
    then
        echo "you can reach the api under the following url:"
        copilot svc show -n api | grep http | sed 's/^.*http/http/'
    else
        echo "you can reach the api under the following url:"
        copilot svc show -a $1 -n api | grep http | sed 's/^.*http/http/'
    fi
}

printHelp () {
    echo "Please provide a valid parameter."
    echo ""
    echo "Use one of the following:"
    echo ""
    echo "init NAME - initialize the app with the provided name and deploy it to the testing environment"
    echo -e "\te.x. ./setup.sh init catapp"
    echo ""
    echo "public NAME - show public reachable api for the given name, if no name is provided use standard"
    echo -e "\te.c. ./setup.sh public catapp"
}

test () {
    echo "start sleep 10"
    sleep 10
    echo "finished sleep"
}

case "$1" in
    "init")
        init $2;;
    "public")
        public $2;;
    "help")
        printHelp
        exit 0;;
    *)
        printHelp
        exit 1;;
esac