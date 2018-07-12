#!/bin/bash
# get the tag name from the user
echo 'Make sure that all branches are the ones you would like to release.'
echo 'Enter the tag name for the release:'
read tagName
cd ../doppler-frontend
docker build . -t acstintern/doppler-frontend:$tagName
docker push acstintern/doppler-frontend:$tagName
git tag $tagName
git push origin $tagName
echo 'Finished tagging doppler-frontend'
cd ../doppler-api
docker build . -t acstintern/doppler-api:$tagName
docker push acstintern/doppler-api:$tagName
git tag $tagName
git push origin $tagName
echo 'Finished tagging doppler-api'
cd ../doppler-events
docker build . -t acstintern/doppler-events:$tagName
docker push acstintern/doppler-events:$tagName
git tag $tagName
git push origin $tagName
echo "Finished tagging doppler-events\nRelease Complete"