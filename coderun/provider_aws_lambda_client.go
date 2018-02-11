package coderun

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"archive/zip"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type CRLambda struct {
	lambda            *lambda.Lambda
	s3                *s3.S3
	s3uploader        *s3manager.Uploader
	bucket            string
	bucketKey         string
	iam               *iam.IAM
	iamRoleName       string
	iamRoleArn        string
	functionCode      *lambda.FunctionCode
	lambdaFunctionArn string
	runEnvironment    RunEnvironment
	sourceFiles       chan string
	lambdaName        string
	codeDir           string
}

func NewCRLambda(awsConfig *aws.Config) *CRLambda {
	d := &CRLambda{}
	sess := session.Must(session.NewSession(awsConfig))
	d.lambda = lambda.New(sess)
	d.iam = iam.New(sess)
	d.s3uploader = s3manager.NewUploader(sess)
	d.s3 = s3.New(sess)

	return d
}

func (d *CRLambda) Setup(r RunEnvironment) {
	d.lambdaName = r.Name
	d.codeDir = r.CodeDir
	//d.getOrCreateS3Bucket()
	d.zipAndUploadCode(r.CodeDir)
	d.createIamRole()
}

func (d *CRLambda) Deploy(r RunEnvironment, p awsLambdaProviderEnv) {
	d.createOrUpdateFunction()
}

func (d *CRLambda) Run(r RunEnvironment, p awsLambdaProviderEnv) {
	d.runFunction()
}

func (d *CRLambda) createIamRole() string {
	if role := d.getConfig("awsLambdaIAMRole"); role != "" {
		d.iamRoleArn = role
		return role
	}

	d.iamRoleName = fmt.Sprintf("CodeRunLambda-%s", RandString(5))
	resp, err := d.iam.CreateRole(&iam.CreateRoleInput{
		Description: aws.String("Create by Coderun for Lambda"),
		RoleName:    aws.String(d.iamRoleName),
		AssumeRolePolicyDocument: aws.String(
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
					}
				]
			}`)})

	if err != nil {
		log.Fatal(err)
	}

	d.iamRoleArn = *resp.Role.Arn
	d.setConfig("awsLambdaIAMRole", d.iamRoleArn)

	return d.iamRoleArn
}

func (d *CRLambda) updateLambdaFunction() {
	contents, err := ioutil.ReadFile(path.Join(".coderun", d.bucketKey))
	if err != nil {
		log.Fatal(err)
	}

	_, err = d.lambda.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(d.lambdaName),
		Publish:      aws.Bool(true),
		ZipFile:      contents,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (d *CRLambda) runFunction() *string {
	resp, err := d.lambda.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(d.lambdaName),
		LogType:      aws.String(lambda.LogTypeTail),
		Payload:      []byte("{}"),
	})
	if err != nil {
		log.Fatal(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(*resp.LogResult)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(decoded))
	return resp.LogResult
}

func (d *CRLambda) createOrUpdateFunction() string {
	d.lambdaFunctionArn = d.getConfig("awsLambdaFunctionArn")
	if d.lambdaFunctionArn != "" {
		d.updateLambdaFunction()
		return d.lambdaFunctionArn
	}
	contents, err := ioutil.ReadFile(path.Join(".coderun", d.bucketKey))
	if err != nil {
		log.Fatal(err)
	}

	createCode := &lambda.FunctionCode{
		ZipFile: contents,
	}

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String(d.lambdaName),
		Handler:      aws.String("test.lambda_handler"),
		Role:         aws.String(d.iamRoleArn),
		Runtime:      aws.String("python3.6"),
	}
	result, err := d.lambda.CreateFunction(createArgs)

	if err != nil {
		log.Fatal(err)
	}
	d.setConfig("awsLambdaFunctionArn", *result.FunctionArn)
	d.lambdaFunctionArn = *result.FunctionArn
	fmt.Println(result)
	return d.lambdaFunctionArn
}

func (d *CRLambda) zipper(output string, files chan string) {
	newfile, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	fmt.Printf("Zipping files for upload to S3:")
	for file := range files {
		log.Printf("zipper got file: %s", file)
		f, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("adding to zip")
		writer, err := zipWriter.Create(file)
		_, err = writer.Write(f)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("done copying")
	}
}

func (d *CRLambda) addFile(fullpath string, f os.FileInfo, err error) error {
	relativePath := strings.TrimPrefix(fullpath, path.Clean(d.codeDir)+"/")
	if f.IsDir() == true || strings.HasPrefix(relativePath, ".") {
		log.Printf("skipping: %s\n", relativePath)
		return nil
	}
	log.Printf("adding file: %s\n", relativePath)
	d.sourceFiles <- relativePath
	return nil
}

func (d *CRLambda) zipAndUploadCode(dir string) {
	// if key := d.getConfig("awsLambdaS3Key"); bucket != "" {
	// 	//check if zip matches and return
	// } else {
	// }

	d.bucketKey = fmt.Sprintf("coderun-%s.zip", RandString(20))
	d.setConfig("awsLambdaS3BucketKey", d.bucketKey)
	zipFile := path.Join(".coderun", d.bucketKey)
	//defer os.Remove(zipFile)

	d.sourceFiles = make(chan string)
	log.Printf("directory: %s", dir)
	go func() { filepath.Walk(dir, d.addFile); close(d.sourceFiles) }()
	d.zipper(zipFile, d.sourceFiles)
	log.Printf("done zipping")

	//r, err := os.Open(zipFile)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Printf("Uploading %s to %s/%s", zipFile, d.bucket, d.bucketKey)
	//log.Printf("bucket is %s", d.bucket)
	//result, err := d.s3uploader.Upload(&s3manager.UploadInput{
	//	Bucket: aws.String(d.bucket),
	//	Key:    aws.String(d.bucketKey),
	//	Body:   r,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	//return d.bucketKey
}

func (d *CRLambda) getConfig(key string) string {
	value, err := ioutil.ReadFile(fmt.Sprintf(".coderun/%s", key))
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal(err)
	}
	return string(value)
}

func (d *CRLambda) setConfig(key string, value string) {
	log.Printf("Setting config %s to %s\n", key, value)
	CreateCodeRunDir()
	err := ioutil.WriteFile(fmt.Sprintf(".coderun/%s", key), []byte(value), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *CRLambda) getOrCreateS3Bucket() {
	d.bucket = d.getConfig("awsLambdaS3Bucket")
	if d.bucket == "" {
		d.bucket = fmt.Sprintf("coderun-%s", RandString(20))
		log.Printf("creating bucket: %s\n", d.bucket)
		_, err := d.s3.CreateBucket(&s3.CreateBucketInput{
			ACL:    aws.String(s3.BucketCannedACLPrivate),
			Bucket: aws.String(d.bucket),
		})
		if err != nil {
			log.Fatal(err)
		}
		d.setConfig("awsLambdaS3Bucket", d.bucket)
	} else {
		log.Printf("found bucket config %s", d.bucket)
	}
	log.Printf("1----, %s", d.bucket)
}
