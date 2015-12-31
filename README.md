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
```
400 uniq images downloaded (1229 images in total)
2757 uniq layers (5640 layers in total)
```
####Format
`$name $tag $CVECountOftheImage [$CVECounts of each layer, ascending order]`
####Explaination
`clojure latest 49 [24 24 27 31 31 31 31 31 31 31 49 49 49 49 49 49 49 49 49 49]`
The 'clojure:latest' image has 49 CVE bugs.
The first layer has 24 bugs.
The second has 24 bugs either, but all the 24 bugs are introduced by the first layer. 
(I do not output the 'layer diff' here).
The third layer has 27 bugs which possibily means it introduces 3 new CVE bugs.
And the etc...
####Conclusion
1. Nearly all the images has CVE issues
   Pay attention: their 'official' images!
2. More layer more CVE
   It means: no 'defensive' layer.


###TODO
1. Continue to make a better 'image-mining' work
2. Clair is slow 
   Even the query takes time. (30ms ~120ms per query)

