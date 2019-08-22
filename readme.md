# kube-linker
> alpha-state, unfinished

### Run locally
export OPERATOR_NAME=kube-linker
operator-sdk up local

### Build
export IMAGE=vkfont/kube-linker:0.0.1
operator-sdk build $IMAGE
docker push $IMAGE

### Run tests
operator-sdk test local ./test

### Annotations
metadata:
  annotations:
    kube-linker/enabled: "true"
    kube-linker/description: Service description
    kube-linker/doc-url: https://documentation.ling
    
    kube-linker/name: optional name