package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

func unquoteString(str string) string {
	if (strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`)) ||
		(strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`)) {
		return str[1 : len(str)-1]
	}
	return str
}

func execAndWait(timeout int, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		logError(err.Error()) // TODO
		return err
	}

	// wait for process finish
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			logError(err.Error()) // TODO
			return err
		}
		<-done
		logInfo("time out and process was killed") // TODO
	case err := <-done:
		logInfo("process output: %s", stdout.String())
		if err != nil {
			logError("process error output: %s", stderr.String())
			logError("process done with error = %v", err) // TODO
			return err
		}
	}
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func unTarGz(archiveFile string, destDir string, prefix string) error {
	destDir = path.Clean(destDir) + string(os.PathSeparator)

	// open the archive file
	fr, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer fr.Close()

	// create a gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	// create a tar reader
	tr := tar.NewReader(gr)

	// loop files
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return err
		}

		if !strings.HasPrefix(hdr.Name, prefix) {
			continue
		}

		fi := hdr.FileInfo()
		destFullPath := destDir + hdr.Name
		logInfo("UnTarGzing file: " + hdr.Name) // TODO

		if hdr.Typeflag == tar.TypeDir {
			// create dir
			os.MkdirAll(destFullPath, fi.Mode().Perm())
			os.Chmod(destFullPath, fi.Mode().Perm())
		} else {
			// create the parent dir for file
			os.MkdirAll(path.Dir(destFullPath), fi.Mode().Perm())

			// write data to file
			fw, err := os.Create(destFullPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
			fw.Close()

			os.Chmod(destFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

// find if a file in archive and return its path
func findFileInTarGz(archiveFile string, targetFile string) (string, error) {
	// open the archive file
	fr, err := os.Open(archiveFile)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	// create a gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	// create a tar reader
	tr := tar.NewReader(gr)

	// loop files
	targetPath := ""
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return "", err
		}

		if hdr.Typeflag != tar.TypeDir && strings.HasSuffix(hdr.Name, targetFile) {
			targetPath = hdr.Name
		}
	}
	return targetPath, nil
}
