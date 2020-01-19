package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var count int
	path := "/host/proc"
	fileName := "/host/proc/sys/fs/file-max"

	parentDirs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("open %s err: %v", path, err)
	}
	for _, sonDir := range parentDirs {
		if sonDir.IsDir() && JudgeDir(sonDir.Name()) {
			_, err := os.Stat(path + "/"+sonDir.Name() + "/fd")
			if err == nil {
				fds, _ := ioutil.ReadDir(path + "/"+sonDir.Name() + "/fd")
				count+=len(fds)
			}
		}
	}

	max, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("open %v err :%v", err)
	}
	m,err:=MaxToInt(max)
	if err!=nil{
		log.Fatalf("Max strconv to int err: %v",err)
	}
	if count>(m*80/100){
		failureinfo:=fmt.Sprintf("current fd usage is %v and max is %v",count,m)
		exec.Command("/bin/bash", "-c", failureinfo)
		os.Exit(1)
	}
	successinfo:=fmt.Sprintf("node has no fd pressure")
	exec.Command("/bin/bash", "-c", successinfo)
	os.Exit(0)
}
//Whether dir name is `[0-9]*`.
func JudgeDir(name string) bool {
	r := regexp.MustCompile(`[0-9]+`)
	return r.MatchString(name)
}
//strconv max to int.
func MaxToInt(sourceMax []byte)(int,error)  {
	tempMax:= strings.Replace(string(sourceMax), "\n", "", -1)
	Max,err:=strconv.Atoi(tempMax)
	if err != nil {
		return 0,err
	}
	return Max,nil
}