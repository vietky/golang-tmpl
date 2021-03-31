WSM_DEFAULT_BRANCH=$1
docker run --rm -it -v `pwd`:/app -v ~/.ssh:/root/.ssh -e "WSM_FETCH_ALL_FIRST=false" -e "WSM_DEFAULT_BRANCH=master" docker.chotot.org/wsm:0.0.3 sh -c "touch /app/tag.txt && wsm get-version > /app/tag.txt"
DOCKER_VERSION=$(cat ./tag.txt)
# echo $DOCKER_VERSION
read -n 2 -p "Next tagging version is '${DOCKER_VERSION}'. Proceed (y/n)? " choice

if [ "$choice" = "y" ]; then
    git tag $DOCKER_VERSION | exit 1
    docker build -t docker.chotot.org/cxsvc-backend:$DOCKER_VERSION .
    docker push docker.chotot.org/cxsvc-backend:$DOCKER_VERSION
    git push origin HEAD:$DOCKER_VERSION
    rm ./tag.txt
fi