//go:build mage
// +build mage

package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
)

var Default = Proto

// regenerate protobuf
func Proto() error {
	twirpProtoFiles := []string{
		"livekitx_recording.proto",
		"livekitx_room.proto",
	}
	grpcProtoFiles := []string{
		"livekitx_internal.proto",
		"livekitx_model.proto",
		"livekitx_rtc.proto",
		"livekitx_webhook.proto",
	}

	fmt.Println("generating protobuf")
	target := "livekitx"
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}
	protocGoPath, err := getToolPath("protoc-gen-go")
	if err != nil {
		return err
	}
	twirpPath, err := getToolPath("protoc-gen-twirp")
	if err != nil {
		return err
	}
	protocGrpcGoPath, err := getToolPath("protoc-gen-go-grpc")
	if err != nil {
		return err
	}
	fmt.Println("generating twirp protobuf")
	args := append([]string{
		"--go_out", target,
		"--twirp_out", target,
		"--go_opt=paths=source_relative",
		"--twirp_opt=paths=source_relative",
		"--plugin=go=" + protocGoPath,
		"--plugin=twirp=" + twirpPath,
		// "-I=C:\\Users\\4658\\go\\bin",
		"-I=/home/go/go/bin",
		// "-I=C:\\Users\\4658\\protoc-3.19.4-win64\\include\\google\\",
		"-I=/home/go/go/src/google",
		"-I=.",
		// "-I=C:\\Users\\4658\\go\\src",
		"-I=/home/go/go/src",
	}, twirpProtoFiles...)
	cmd := exec.Command(protoc, args...)
	connectStd(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("generating grpc protobuf")
	args = append([]string{
		"--go_out", target,
		"--go-grpc_out", target,
		"--go_opt=paths=source_relative",
		"--go-grpc_opt=paths=source_relative",
		"--plugin=go=" + protocGoPath,
		"--plugin=go-grpc=" + protocGrpcGoPath,
		// "-I=C:\\Users\\4658\\go\\bin",
		"-I=/home/go/go/bin",
		// "-I=C:\\Users\\4658\\protoc-3.19.4-win64\\include\\google\\",
		"-I=/home/go/go/src/google",
		"-I=.",
		// "-I=C:\\Users\\4658\\go\\src",
		"-I=/home/go/go/src",
	}, grpcProtoFiles...)
	cmd = exec.Command(protoc, args...)
	connectStd(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// run tests
func Test() error {
	cmd := exec.Command("go", "test", "./...")
	connectStd(cmd)
	return cmd.Run()
}

// helpers

func getToolPath(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	// check under gopath
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	p := filepath.Join(gopath, "bin", name)
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	return p, nil
}

func connectStd(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}
