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
)

type CRLambda struct {
	lambda            *lambda.Lambda
	iam               *iam.IAM
	iamRoleName       string
	iamRoleArn        string
	lambdaFunctionArn string
	sourceFiles       chan string
	lambdaName        string
	codeDir           string
	zipFile           string
}

func NewCRLambda(awsConfig *aws.Config) *CRLambda {
	d := &CRLambda{}
	sess := session.Must(session.NewSession(awsConfig))
	d.lambda = lambda.New(sess)
	d.iam = iam.New(sess)

	return d
}

func (d *CRLambda) Setup(r RunEnvironment) {
	d.lambdaName = r.Name
	d.codeDir = r.CodeDir

	d.zipAndUploadCode(r.CodeDir)
	d.createIamRole()
}

func (d *CRLambda) Deploy(r RunEnvironment, p awsLambdaProviderEnv) {
	d.lambdaFunctionArn = d.getConfig("awsLambdaFunctionArn")

	if d.lambdaFunctionArn != "" {
		d.updateLambdaFunction()
	} else {
		arn := d.createLambdaFunction()
		d.setConfig("awsLambdaFunctionArn", arn)
		d.lambdaFunctionArn = arn
	}
	os.Remove(d.zipFile)
}

func (d *CRLambda) Run(r RunEnvironment, p awsLambdaProviderEnv) {
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
	fmt.Print(string(decoded))
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

func (d *CRLambda) createLambdaFunction() string {
	contents, err := ioutil.ReadFile(d.zipFile)
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
	resp, err := d.lambda.CreateFunction(createArgs)

	if err != nil {
		log.Fatal(err)
	}
	return *resp.FunctionArn
}

func (d *CRLambda) updateLambdaFunction() {
	contents, err := ioutil.ReadFile(d.zipFile)
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

func (d *CRLambda) zipper(output string, files chan string) {
	Logger.info.Printf("Creating zip file at %s", output)
	newfile, err := os.Create(output)
	defer newfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	for file := range files {
		Logger.debug.Printf("Zipper got file: %s", file)
		f, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		writer, err := zipWriter.Create(file)
		_, err = writer.Write(f)
		Logger.debug.Printf("Zipper wrote file: %s", file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *CRLambda) addFile(fullpath string, f os.FileInfo, err error) error {
	relativePath := strings.TrimPrefix(fullpath, path.Clean(d.codeDir)+"/")

	if f.IsDir() == true || strings.HasPrefix(relativePath, ".") {
		Logger.debug.Printf("Skipping: %s\n", relativePath)
		return nil
	}

	Logger.info.Printf("Found file: %s\n", relativePath)
	d.sourceFiles <- relativePath
	return nil
}

func (d *CRLambda) zipAndUploadCode(dir string) {
	d.zipFile = path.Join(".coderun", "lambda-"+RandString(20)+".zip")
	Logger.debug.Printf("Zipping directory: %s", dir)

	d.sourceFiles = make(chan string)
	go func() { filepath.Walk(dir, d.addFile); close(d.sourceFiles) }()
	d.zipper(d.zipFile, d.sourceFiles)

	Logger.info.Printf("Done zipping")
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
	Logger.info.Printf("Setting config %s to %s\n", key, value)

	err := os.Mkdir(".coderun", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf(".coderun/%s", key), []byte(value), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
