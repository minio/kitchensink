# kitchensink

Kitchensink is a command-line tool that creates a set of nested objects of different sizes in a specified bucket and verifies the data by cross-checking md5sum of the objects before and after upload. This tool is helpful since it helps validate objects between different releases and versions. For example, if you started with an older version and upgrade to a newer version, this tool will verify if the hash for each object matches between the different versions. 

After pulling the code locally, use Make to install and build the program. 

## Commands
```
**create**        Creates a nested folder structure with objects of random prime number sizes in a pre-existing bucket

USAGE:
    kitchensink create [ARGUMEMTS] [FLAGS]

ARGUMENTS:
    endpoint
    access key
    secret key
    bucket name

FLAGS:
    --insecure       skips verification in transport when putting objects
    --help           displays help menu

EXAMPLE:
    kitchensink create https://play.min.io Q3AM3UQ867SPQQA43P2F zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG my-test-bucket 

```

```
**verify**        Verify gets each object and computes the hash. It then compares this hash with the metadata from create. 
USAGE:
    kitchensink verify [ARGUMEMTS] [FLAGS]

ARGUMENTS:
    endpoint
    access key
    secret key
    bucket name
    
FLAGS:
    --insecure       skips verification in transport options
    --help           displays help menu

EXAMPLE:
    kitchensink verify https://play.min.io Q3AM3UQ867SPQQA43P2F zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG my-test-bucket
```

```
**delete**        deletes all the objects in the specified bucket, cleans the bucket

USAGE:
    kitchensink delete [ARGUMEMTS] 

ARGUMENTS:
    endpoint
    access key
    secret key
    bucket name
    
EXAMPLE:
    kitchensink delete https://play.min.io Q3AM3UQ867SPQQA43P2F zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG my-test-bucket

```

  

