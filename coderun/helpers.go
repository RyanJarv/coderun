package coderun

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
)

func getImageName() (string, error) {
	var image, err = ioutil.ReadFile(".coderun/dockerimage")
	if err != nil {
		return "", err
	}
	return string(image), nil
}

func setImageName(image string) error {
	return ioutil.WriteFile(".coderun/dockerimage", []byte(image), 0644)
}

func cmd(c ...string) string {
	cmd := exec.Command(c[0], c[1:]...)
	log.Printf("%v", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("output: %s", out.String())
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s", err))
	}
	return out.String()
}

func newImageName() string {
	return fmt.Sprintf("coderun-%s", randString())
}

func randString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	log.Printf("randString: %s", string(b))
	return string(b)
}
