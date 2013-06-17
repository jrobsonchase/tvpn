package pki

import (
	"os/exec"
	"fmt"
)

func GenDH(ssl,outFile string,bits int) ([]byte, error) {
	cmd := exec.Command(ssl,"dhparam","-out",outFile,fmt.Sprintf("%d",bits))
	return cmd.CombinedOutput()
}
