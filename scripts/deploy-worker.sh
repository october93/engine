#!/usr/bin/env bash
set -euxo pipefail

ENV=${ENV:-development}
REPOSITORY=740130560720.dkr.ecr.us-west-2.amazonaws.com

go build -o worker/${TASK}/${TASK} ./cmd/${TASK}
cp ${ENV}.config.toml worker/${TASK}/config.toml

eval $(aws ecr get-login --no-include-email --region us-west-2)
docker build -t ${TASK}:${ENV} worker/${TASK}
docker tag ${TASK}:${ENV} ${REPOSITORY}/${TASK}:${ENV}
docker push ${REPOSITORY}/${TASK}:${ENV}

./scripts/ci/ecs-deploy.sh  -c ${TASK}-${ENV} -n ${TASK}-${ENV} -i ${REPOSITORY}/${TASK}:${ENV} -r us-west-2

rm worker/${TASK}/config.toml
