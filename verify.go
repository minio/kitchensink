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

//Processes arguments from command line
func mainVerify(ctx *cli.Context) error {
	argsLength := len(ctx.Args())
	if argsLength != 4 {
		cli.ShowCommandHelpAndExit(ctx, "verify", 1)
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
	Verify(URLstr, bucketname, options)
	return nil
}

//Verify will compare the md5sums of each of the object
func Verify(endpoint string, bucketname string, options minio.Options) {
	fmt.Println("VERIFYING DATASET...")

	s3Client, err := minio.New(endpoint, &options)
	if err != nil {
		log.Fatalln(err)
	}

	objectList := s3Client.ListObjects(context.Background(), bucketname, minio.ListObjectsOptions{
		WithMetadata: true,
		//Prefix:       "myprefix",
		Recursive: true,
	})

	var errCount int
	for object := range objectList {
		if object.Err != nil {
			log.Printf("ERR")
			fmt.Println(object.Err)
			return
		}

		data, err := s3Client.GetObject(context.Background(), bucketname, object.Key, minio.GetObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		hash := md5.New()
		_, err = io.Copy(hash, data)
		if err != nil {
			log.Fatalln(err)
		}

		md5val := hex.EncodeToString(hash.Sum(nil))
		metadata := object.UserMetadata["X-Amz-Meta-Content-Md5"]
		if md5val != metadata {
			log.Println("ERR: Object", object.Key, "hash does not match")
			errCount++
		}
	}

	if errCount == 0 {
		fmt.Println("Successfully Verified")
	} else {
		fmt.Println("Finished with ", errCount, "errors")
	}

}