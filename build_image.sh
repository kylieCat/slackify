#!/usr/bin/env bash
set -e

if [[ $# -lt 1 ]]; then
    echo "Usage: build_image.sh -r REPO_NAME -v IMAGE_VERSION [-p]" >&2
    exit 1
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/slackify main.go

while getopts "r:v:p" OPTION
do
    case $OPTION in
        r | --repo)
            repo=${OPTARG}
            ;;
        v | --version)
            new_tag=${OPTARG}
            echo "Using tag: app-$new_tag"
            docker build -t slackify:"$new_tag" .
            docker tag slackify:"$new_tag" $repo/slackify:"$new_tag"
            echo "Built image: $repo/slackify:$new_tag"
            ;;
        p | --push)
            echo "Pushing image to DockerHub: $repo/slackify:$new_tag"
            docker push $repo/slackify:$new_tag
            ;;
    esac
done
