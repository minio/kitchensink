# kitchensink

Kitchensink is a command-line tool that creates a set of nested objects of different sizes in a specified bucket and verifies the data by cross-checking md5sum of the objects before and after upload. This tool is helpful since it helps validate objects between different releases and versions. For example, if you started with an older version and upgrade to a newer version, this tool will verify if the hash for each object matches between the different versions. 

After pulling the code locally, use Make to install and build the program. 

## Commands
```
create        creates a nested folder structure with objects of random prime number sizes in the specified bucket
verify        verifies the data in the bucket by comparing the md5sum of the object from when it was uploaded to when it was retrieved. 
delete        deletes all the objects in the specified bucket
```
**Note:** All of the commands listed above take in the endpoint server along with the credentials (access key and secret key) for that server as arguments. In addition, the commands also need a bucketname where you want to create, verify, or delete. An example of the create command is shown below.   

## Example 
```
./kitchensink create https://endpoint ACCESSKEY SECRETKEY BUCKETNAME
```
