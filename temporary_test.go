package tempfiles

import (
    "testing"
    "fmt"
)

const nTries = 1000

func testTempFiles(nTries int) (err error) {
    t, err := NewDefault()
    if err != nil {
        return
    }

    for ii := 0; ii < nTries; ii++ {
        vf, err := t.Get()
        if err != nil {
            return err
        }

        if err = vf.MustExist(); err != nil {
            return err
        }

        t.Put(vf)
    }
    return t.Remove()
}


func TestTempFiles(t *testing.T) {
    if err := testTempFiles(nTries); err != nil {
        t.Fatalf("Error: %s\n", err)
    }
}

func withExtension(ext string, nTries int) (err error) {
    f, err := NewDefault()
    if err != nil {
        return
    }
    w := WithExtension{Like: &f, extension: ext}

    for ii := 0; ii < nTries; ii++ {
        file, err := w.Get()
        if err != nil {
            return err
        }

        if fext := file.Ext(); ext != fext {
            return fmt.Errorf("expected file with extension '%s', got '%s'",
                ext, fext)
        }
    }
    return
}

func TestWithExtension(t *testing.T) {
    err := withExtension("png", 1000000)
    if err != nil {
        t.Fatalf("Error: %s\n", err)
    }
}
