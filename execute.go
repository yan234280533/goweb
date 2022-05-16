package main

import (
	"os/exec"

	"k8s.io/klog/v2"
)

func Start(filename string) error {
	cmd := exec.Command("/bin/bash", "-c", filename)
	err := cmd.Start()
	if err != nil {
		klog.Errorf("failed to start command start:%v", err)
		return err
	}
	klog.Infof("command start succeed")
	return nil
}

func Restart(serviceName string, filename string) error {
	err := Stop(serviceName)
	if err != nil {
		//continue
	}

	return Start(filename)
}

func Stop(serviceName string) error {
	cmd := exec.Command("pkill", serviceName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("combined out:\n%s\n", string(out))
		return err
	}

	klog.Infof("combined out:\n%s\n", string(out))
	return nil
}
