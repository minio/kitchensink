# kitchensink

Kitchensink is a command-line tool that creates a dataset on a specified endpoint and verifies the data through the 
md5sum of the objects. Use Make to build and install the program locally.


Commands
```
create                  creates a dataset in the endpoint on the specified bucket
verify                  verifies the data in the bucket by checking the md5sum
```

Examples 
```
./kitchensink create https://endpoint ACCESSKEY SECRETKEY BUCKETNAME
```
