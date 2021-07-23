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
	"log"

	"github.com/minio/cli"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//Processes inputs from command line
func mainDelete(ctx *cli.Context) error {
	argsLength := len(ctx.Args())
	if argsLength != 4 {
		cli.ShowCommandHelpAndExit(ctx, "delete", 1)
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
	Delete(URLstr, bucketname, options)
	return nil
}

//Delete removes all the objects in the specified bucket
func Delete(endpoint string, bucketname string, options minio.Options) {

	log.Println("CLEANING BUCKET...")

	s3Client, err := minio.New(endpoint, &options)
	if err != nil {
		log.Fatalln(err)
	}

	objectList := s3Client.ListObjects(context.Background(), bucketname, minio.ListObjectsOptions{
		WithMetadata: true,
		Recursive:    true,
	})

	for object := range objectList {
		if object.Err != nil {
			log.Println("Error occured:", object.Err)
			return
		}
		Rerr := s3Client.RemoveObject(context.Background(), bucketname, object.Key, minio.RemoveObjectOptions{})
		if Rerr != nil {
			log.Fatalln(err)
		}
	}

	log.Println("BUCKET CLEANED SUCCESSFULLY...")
}
