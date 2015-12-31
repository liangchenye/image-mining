##CVE data mining on images.
1. Aim to find the CVE status report on images.
2. Aim to make 'Clair' better by pushing lots of insert/query.

###Codes
```
./collect
    Way to get images
./collect/imagesource.go
    Provides an interface to download images for different container hub/repo
./collect/imagesource/docker-official.go
    The first implementation of imagesource interface, using repos in 'docker-library/official-image'

./libs
    Libs used by image-mining
./libs/image.go
    Docker image in reality. Provide ways to pull/save/scan a docker image by its name.
```
###Output
#### Docker-offical-result
./docker-official-result
####Format
`$name $tag $CVECountOftheImage [$CVECounts of each layer, ascending order]`
####Example
