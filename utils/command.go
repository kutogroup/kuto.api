package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

//CommandGetResult 获取命令行结果
func CommandGetResult(shellPath string, params ...string) string {
	cmd := exec.Command(shellPath, params...)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return ""
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return ""
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return ""
	}

	return string(bytes)
}

//CommandExecute 执行命令行
func CommandExecute(out io.Writer, shellPath string, params ...string) {
	cmd := exec.Command(shellPath, params...)
	cmd.Stdout = os.Stdout
	if out == nil {
		cmd.Stderr = os.Stdout
	} else {
		cmd.Stderr = out
	}

	cmd.Run()
}
