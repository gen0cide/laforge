package command

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codegangsta/cli"
	humanize "github.com/dustin/go-humanize"
	"github.com/gen0cide/laforge/competition"
	"github.com/olekukonko/tablewriter"
)

func CmdCDN(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdCDNLs(c *cli.Context) {
	comp, env := InitConfig()
	conn := GetS3Service(comp, env)
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(comp.S3Config.Bucket),
		MaxKeys: aws.Int64(9001),
	}

	result, err := conn.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(s3.ErrCodeNoSuchBucket, aerr.Error()))
			default:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(aerr.Error()))
			}
		} else {
			competition.LogFatal("S3 Error: " + fmt.Sprintln(err.Error()))
		}
		return
	}
	fileList := [][]string{}

	competition.Log("CDN File Listing")

	for _, obj := range result.Contents {
		fileList = append(fileList, []string{aws.StringValue(obj.Key), humanize.Bytes(uint64(aws.Int64Value(obj.Size))), obj.LastModified.String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Size", "Last Modified"})

	for _, v := range fileList {
		table.Append(v)
	}
	table.Render()

}

func CmdCDNLink(c *cli.Context) {
	objKey := c.Args().Get(0)
	if len(objKey) < 1 {
		competition.LogFatal("Please specify a filename as an argument to the link command.")
	}
	comp, env := InitConfig()
	conn := GetS3Service(comp, env)
	input := &s3.GetObjectInput{
		Bucket: aws.String(comp.S3Config.Bucket),
		Key:    aws.String(objKey),
	}

	_, err := conn.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(s3.ErrCodeNoSuchBucket, aerr.Error()))
			case s3.ErrCodeNoSuchKey:
				competition.LogFatal("S3 Error: File Does Not Exist")
			default:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(aerr.Error()))
			}
		} else {
			competition.LogFatal("S3 Error: " + fmt.Sprintln(err.Error()))
		}
		return
	}

	u, err := url.Parse(fmt.Sprintf("https://s3-%s.amazonaws.com/%s", comp.S3Config.Region, filepath.Join(comp.S3Config.Bucket, objKey)))
	if err != nil {
		competition.LogFatal("url error: " + err.Error())
	}
	competition.Log("Link for file: " + objKey)
	competition.LogPlain(u.String())

}

func CmdCDNUpload(c *cli.Context) {
	comp, env := InitConfig()
	conn := GetS3Service(comp, env)
	srcFile := c.Args().Get(0)
	if len(srcFile) < 1 {
		competition.LogFatal("Please specify a local file as the first argument.")
	}
	objKey := c.Args().Get(1)
	if len(objKey) < 1 {
		competition.LogFatal("Please specify a remote filename as the second argument.")
	}

	file, err := os.Open(srcFile)
	if err != nil {
		competition.LogFatal(fmt.Sprintf("err opening file: %s", err))
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(comp.S3Config.Bucket),
		Key:           aws.String(objKey),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String(s3.ObjectCannedACLPublicRead),
	}

	_, err = conn.PutObject(params)
	if err != nil {
		competition.LogFatal("S3 Error: " + fmt.Sprintf("bad response: %s", err))
	}
	competition.Log(fmt.Sprintf("File successfully uploaded:"))
	competition.Log(fmt.Sprintf("%s => %s", srcFile, objKey))
	u, err := url.Parse(fmt.Sprintf("https://s3-%s.amazonaws.com/%s", comp.S3Config.Region, filepath.Join(comp.S3Config.Bucket, objKey)))
	if err != nil {
		competition.LogFatal("url error: " + err.Error())
	}
	competition.Log("-- URL: " + u.String())
}

func CmdCDNRm(c *cli.Context) {
	objKey := c.Args().Get(0)
	if len(objKey) < 1 {
		competition.LogFatal("Please specify a filename as an argument to the rm command.")
	}
	comp, env := InitConfig()
	conn := GetS3Service(comp, env)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(comp.S3Config.Bucket),
		Key:    aws.String(objKey),
	}

	_, err := conn.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(s3.ErrCodeNoSuchBucket, aerr.Error()))
			case s3.ErrCodeNoSuchKey:
				competition.LogFatal("S3 Error: File Does Not Exist")
			default:
				competition.LogFatal("S3 Error: " + fmt.Sprintln(aerr.Error()))
			}
		} else {
			competition.LogFatal("S3 Error: " + fmt.Sprintln(err.Error()))
		}
		return
	}
	competition.Log(fmt.Sprintf("File Deleted: %s", objKey))
}

func GetS3Service(c *competition.Competition, e *competition.Environment) *s3.S3 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(e.AWSConfig.Region),
			Credentials: credentials.NewStaticCredentials(c.AWSCred.APIKey, c.AWSCred.APISecret, ""),
		},
	}))

	return s3.New(sess, &aws.Config{
		Region: aws.String(c.S3Config.Region),
	})
}
