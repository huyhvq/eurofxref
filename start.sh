#!/bin/bash

help()
{
   echo ""
   echo "Args:"
   echo -e "\tmigrate - Run migration"
   exit 1 # Exit script after printing help
}

case $1 in
  help)
    help
    ;;
  migrate)
    docker-compose exec api /ws/start.sh/start.sh migrate
    ;;
  *)
    docker-compose up
    ;;
esac
