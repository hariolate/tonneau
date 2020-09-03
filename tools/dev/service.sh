#!/usr/bin/env bash

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" >/dev/null 2>&1 && pwd)"

mkdir -p ${project_root}/tools/dev/bin
mkdir -p ${project_root}/tools/dev/db_data
mkdir -p ${project_root}/tools/dev/redis_data

docker-compose -f ${project_root}/tools/dev/docker-compose.yml  --project-directory ${project_root} "$@"