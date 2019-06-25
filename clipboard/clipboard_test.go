package clipboard

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test createWindow, getSelectionOwner
func Test_getSelectionOwner(t *testing.T) {
	conn, err := x.NewConn()
	if err != nil {
		t.Skip("failed to connect x")
	}
	defer func() {
		conn.Close()
	}()

	win, err := createWindow(conn)
	require.Nil(t, nil, "%v", err)
	assert.NotZero(t, win)

	sel, err := conn.GetAtom("CLIPBOARD_test")
	require.Nil(t, err, "%v", err)
	assert.NotZero(t, sel)

	err = x.SetSelectionOwnerChecked(conn, win, sel, x.CurrentTime).Check(conn)
	require.Nil(t, err, "%v", err)

	owner, err := getSelectionOwner(conn, sel)
	assert.Equal(t, win, owner)
	assert.Nil(t, err, "%v", err)
}

func Test_getBytesMd5sum(t *testing.T) {
	sum := getBytesMd5sum([]byte("hello world"))
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", sum)
}

func Test_emptyDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "dde-daemon-clipboard-test")
	if err != nil {
		assert.FailNow(t, "failed to create temp dir: %v", err)
	}
	t.Log("dir:", dir)
	err = ioutil.WriteFile(filepath.Join(dir, "1"), []byte("abc"), 0644)
	assert.Nil(t, err, "%v", err)

	err = os.Mkdir(filepath.Join(dir, "d1"), 0755)
	assert.Nil(t, err, "%v", err)

	err = ioutil.WriteFile(filepath.Join(dir, "d1/1"), []byte("abc"), 0644)
	assert.Nil(t, err, "%v", err)

	err = emptyDir(dir)
	assert.Nil(t, err, "%v", err)

	err = os.Remove(dir)
	assert.Nil(t, err, "%v", err)
}
