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
	crand "crypto/rand"
	"math/rand"
	"strconv"

	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/minio/cli"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3Client *minio.Client

//Processes inputs from command line
func mainCreate(ctx *cli.Context) error {
	argsLength := len(ctx.Args())
	if argsLength != 4 {
		cli.ShowCommandHelpAndExit(ctx, "create", 1)
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
func createObject() (obj *bytes.Buffer, md5val string, size int64) {
	hash := md5.New()
	object := bytes.NewBuffer(nil)
	writer := io.MultiWriter(hash, object)
	bits := rand.Intn(24) + 3
	//Creates the prime number file size
	prime, _ := crand.Prime(crand.Reader, bits)

	_, err := io.CopyN(writer, crand.Reader, prime.Int64())
	if err != nil {
		log.Fatalln(err)
	}
	md5val = hex.EncodeToString(hash.Sum(nil))

	return object, md5val, prime.Int64()
}

//Puts the passed in object with specified objectname and size
func putObject(object *bytes.Buffer, md5val string, size int64, bucketname string, objectname string) {

	object, md5, size := createObject()
	metadata := map[string]string{
		"content-md5": md5,
	}
	_, err := s3Client.PutObject(context.Background(), bucketname, objectname, object, size,
		minio.PutObjectOptions{
			ContentType:  "application/octet-stream",
			UserMetadata: metadata,
		})
	if err != nil {
		log.Fatalln(err)
	}

}

// Create populates a specified bucket with random files in a nested directory structure
func Create(endpoint string, bucketname string, options minio.Options) {

	fmt.Println("CREATING DATASET...")

	s3Client, _ = minio.New(endpoint, &options)

	//making folders with nested object
	for i := 0; i < 5; i++ {
		for j := 0; j < 1; j++ {

			//within outside folder file
			obj, md5sum, fileSize := createObject()
			fileName := "folder00" + strconv.Itoa(i) + "/tests00@" + strconv.Itoa(j)
			putObject(obj, md5sum, fileSize, bucketname, fileName)
		}
		//nested folder files
		testObject, md5val, testSize := createObject()
		fname := "folder00" + strconv.Itoa(i) + "/" + "folder0" + strconv.Itoa(i) + "/tests0" + strconv.Itoa(i)
		putObject(testObject, md5val, testSize, bucketname, fname)

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
		object, md5, size := createObject()
		name := "folder00" + strconv.Itoa(num) + "/folder0" + strconv.Itoa(num) + "/folder-test/folder-x/folder/sample!&" + strconv.Itoa(num)
		putObject(object, md5, size, bucketname, name)

	}
	log.Println("FINISHED SUCCESSFULLY")
}
