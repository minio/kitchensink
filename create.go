// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"context"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"strconv"

	"github.com/minio/cli"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3Client *minio.Client

const jsonFilename = "kitchensink"

//var etags []ETagData
var etags = map[string]string{}

//Processes inputs from command line
func mainCreate(ctx *cli.Context) error {
	argsLength := len(ctx.Args())
	if argsLength != 4 {
		cli.ShowCommandHelpAndExit(ctx, "create", 1)
		return nil
	}
	endpoint := ctx.Args().Get(0)
	secure, URLstr, transport := validateEndpoint(ctx, endpoint)
	access := ctx.Args().Get(1)
	secret := ctx.Args().Get(2)
	bucketname := ctx.Args().Get(3)

	options := minio.Options{
		Creds:     credentials.NewStaticV4(access, secret, ""),
		Secure:    secure,
		Transport: transport,
	}
	Create(URLstr, bucketname, options)
	return nil
}

//Creates an object of prime number size and gets the md5 hash of object
func createObject(isZero bool) (obj *bytes.Buffer, md5val string, size int64) {
	hash := md5.New()
	object := bytes.NewBuffer(nil)
	writer := io.MultiWriter(hash, object)

	fileSize := int64(0)
	//creates 0 byte object
	if isZero == true {
		_, err := io.CopyN(writer, object, 0)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		//Creates the random prime number file size

		bits := rand.Intn(23) + 2
		prime, _ := crand.Prime(crand.Reader, bits)
		fileSize = prime.Int64()

		_, err := io.CopyN(writer, crand.Reader, fileSize)
		if err != nil {
			log.Fatalln(err)
		}
	}

	md5val = hex.EncodeToString(hash.Sum(nil))

	return object, md5val, fileSize
}

//Puts the passed in object with specified objectname and size
func putObject(object *bytes.Buffer, md5val string, size int64, bucketname string, objectname string) {

	//object, md5, size := createObject()
	metadata := map[string]string{
		"content-md5": md5val,
	}
	uploadInfo, err := s3Client.PutObject(context.Background(), bucketname, objectname, object, size,
		minio.PutObjectOptions{
			ContentType:  "application/octet-stream",
			PartSize:     1024*1024*5 + 7, //random prime but greater than 5mb
			UserMetadata: metadata,
		})
	if err != nil {
		log.Fatalln(err)
	}

	etags[uploadInfo.Key] = uploadInfo.ETag
}

// Create populates a specified bucket with random files in a nested directory structure
func Create(endpoint string, bucketname string, options minio.Options) {

	var err error
	s3Client, err = minio.New(endpoint, &options)
	if err != nil {
		log.Fatalln(err)
	}

	//checks if bucket exists otherwise makes the bucket
	if found, _ := s3Client.BucketExists(context.Background(), bucketname); !found {
		s3Client.MakeBucket(context.Background(), bucketname, minio.MakeBucketOptions{})
	}

	log.Println("Creating Dataset")

	//making folders with nested object
	for i := 0; i < 3; i++ {

		//within outside folder file
		obj, md5sum, fileSize := createObject(false)
		fileName := "folder00" + strconv.Itoa(i) + "/tests00@" + strconv.Itoa(i)
		putObject(obj, md5sum, fileSize, bucketname, fileName)

		//nested folder files
		testObject, md5val, testSize := createObject(false)
		fname := "folder00" + strconv.Itoa(i) + "/" + "folder0" + strconv.Itoa(i) + "/tests0" + strconv.Itoa(i)
		putObject(testObject, md5val, testSize, bucketname, fname)

	}

	obj, md5sum, fileSize := createObject(true)
	putObject(obj, md5sum, fileSize, bucketname, "zero-bytes")

	//creates deeply nested object
	nestedObject, md5hash, size := createObject(false)
	name := "folder000/folder00/folder-test/folder-x/folder/sample!&"
	putObject(nestedObject, md5hash, size, bucketname, name)

	dataBytes, err := json.Marshal(etags)
	if err != nil {
		log.Fatalln("error:", err)
	}

	object := bytes.NewBuffer(dataBytes)

	//hashes the json file
	hash := md5.New()
	written, _ := io.Copy(hash, bytes.NewReader(dataBytes))
	md5val := hex.EncodeToString(hash.Sum(nil))
	putObject(object, md5val, written, bucketname, jsonFilename)

	log.Println("Finished Successfully")
}
