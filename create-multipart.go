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
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/minio/cli"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//processes inputs from command line
func mainCreateMulti(ctx *cli.Context) error {
	argsLength := len(ctx.Args())
	if argsLength != 4 {
		cli.ShowCommandHelpAndExit(ctx, "createM", 1)
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
	CreateMultipart(URLstr, bucketname, options)
	return nil
}

func createMultiObject() (obj *bytes.Buffer, md5val string, size int64) {
	//Creates an object of prime number size and gets the md5 hash of object
	hash := md5.New()
	object := bytes.NewBuffer(nil)
	writer := io.MultiWriter(hash, object)
	prime, _ := crand.Prime(crand.Reader, 19)
	log.Println("Multi FILE SIZE", prime)

	_, err := io.CopyN(writer, crand.Reader, prime.Int64())
	if err != nil {
		log.Fatalln(err)
	}
	md5val = hex.EncodeToString(hash.Sum(nil))

	return object, md5val, prime.Int64()
}

func putObjectMulti(object *bytes.Buffer, md5val string, size int64, bucketname string, objectname string) {
	//puts the passed in object with specified objectname and size
	object, md5, size := createMultiObject() //file sizes greater than partsize
	metadata := map[string]string{
		"content-md5": md5,
	}
	_, err := s3Client.PutObject(context.Background(), bucketname, objectname, object, -1,
		minio.PutObjectOptions{
			ContentType:  "application/octet-stream",
			UserMetadata: metadata,
			PartSize:     1024*1024*5 + 7, //random,prime but greater than 5mb, pick a random prime greater than 5mb within a range
		})
	if err != nil {
		log.Fatalln(err)
	}

}

// CreateMultipart populates a specified bucket with random files
func CreateMultipart(endpoint string, bucketname string, options minio.Options) {
	fmt.Println("CREATING MULTI-PART DATASET...")

	s3Client, _ = minio.New(endpoint, &options)

	//making folders with nested object
	for i := 0; i < 5; i++ {
		for j := 0; j < 1; j++ {

			//within outside folder file
			obj, md5sum, fileSize := createMultiObject()
			fileName := "folder00" + strconv.Itoa(i) + "/tests00@" + strconv.Itoa(j)
			putObjectMulti(obj, md5sum, fileSize, bucketname, fileName)
		}
		//nested folder files
		testObject, md5val, testSize := createMultiObject()
		fname := "folder00" + strconv.Itoa(i) + "/" + "folder0" + strconv.Itoa(i) + "/tests0" + strconv.Itoa(i)
		putObjectMulti(testObject, md5val, testSize, bucketname, fname)

	}
	hash := md5.New()
	object := bytes.NewBuffer(nil)
	writer := io.MultiWriter(hash, object)
	_, err := io.CopyN(writer, object, 0)
	if err != nil {
		log.Fatalln(err)
	}
	md5val := hex.EncodeToString(hash.Sum(nil))
	zeroByte, err := s3Client.PutObject(context.Background(), bucketname, "zero-bytes", object, 0,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			UserMetadata: map[string]string{
				"content-md5": md5val,
			},
		})
	_ = zeroByte
	//creates deeply nested object
	for num := 0; num < 1; num++ {
		object, md5, size := createMultiObject()
		name := "folder00" + strconv.Itoa(num) + "/folder0" + strconv.Itoa(num) + "/folder-test/folder-x/folder/sample!&" + strconv.Itoa(num)
		putObjectMulti(object, md5, size, bucketname, name)

	}
	log.Println("FINISHED SUCCESSFULLY")
}
