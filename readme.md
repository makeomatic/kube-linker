### Run locally
export OPERATOR_NAME=kube-linker
operator-sdk up local --namespace=default

### Build
export IMAGE=...
operator-sdk build $IMAGE
docker push $IMAGE

### Run tests
https://github.com/operator-framework/operator-sdk-samples/tree/master/memcached-operator/test/e2e

operator-sdk test local ./test
