### Run locally
export OPERATOR_NAME=kube-linker
operator-sdk up local

### Build
export IMAGE=...
operator-sdk build $IMAGE
docker push $IMAGE

### Run tests
https://github.com/operator-framework/operator-sdk-samples/tree/master/memcached-operator/test/e2e

operator-sdk test local ./test


### TODO
[x] create ingress fetcher proof-of-concept
[x] handle valid ingresses: skip by labels, extract description, place to storage
[x] create web-server to display ingresses
[ ] made nice-looking UI
[ ] add some tests
[ ] create deployment
[ ] deploy to streamlayer
[?] add virtualservices