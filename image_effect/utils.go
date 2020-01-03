package image_effect

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pkg.deepin.io/lib/utils"
)

func getOutputFile(effect, filename string) (outputFile string) {
	outputDir := filepath.Join(cacheDir, effect)
	md5sum, _ := utils.SumStrMd5(filename)
	outputFile = filepath.Join(outputDir, md5sum+filepath.Ext(filename))
	return
}

func modTimeEqual(t1, t2 time.Time) bool {
	return t1.Unix() == t2.Unix() &&
		(t1.Nanosecond()/1000) == (t2.Nanosecond()/1000)
}

func setFileModTime(filename string, t time.Time) error {
	now := time.Now()
	return os.Chtimes(filename, now, t)
}

func runCmdRedirectStdOut(uid int, outputFile string, cmdline, envVars []string) error {
	args := append([]string{"-u", "#" + strconv.Itoa(uid)}, envVars...)
	args = append(args, cmdline...)
	cmd := exec.Command("sudo", args...)
	logger.Debugf("$ sudo %s > %q", strings.Join(args, " "), outputFile)
	cmd.Env = append(os.Environ(), envVars...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	fh, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		err = fh.Close()
		if err != nil {
			logger.Warning(err)
		}
	}()
	bufWriter := bufio.NewWriter(fh)
	var n int64
	n, err = io.Copy(bufWriter, stdout)
	logger.Debugf("copy %d bytes", n)
	if err != nil {
		return err
	}
	err = bufWriter.Flush()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if len(errBuf.Bytes()) > 0 {
		logger.Warningf("cmd stderr: %s", errBuf.Bytes())
	}
	if err != nil {
		return err
	}

	return nil

}
