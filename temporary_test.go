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
        vf, err := Create(&t)
        if err != nil {
            return err
        }

        if err = vf.MustExist(); err != nil {
            return err
        }

        t.Put(vf.ToFile())
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
        file, err := Create(&w)
        if err != nil {
            return err
        }
        fmt.Printf("%s\n", file)
        defer w.Put(file.ToFile())

        if fext := file.Ext(); ext != fext {
            return fmt.Errorf("expected file with extension '%s', got '%s'",
                ext, fext)
        }
    }
    return
}

func TestWithExtension(t *testing.T) {
    err := withExtension("png", 10)
    if err != nil {
        t.Fatalf("Error: %s\n", err)
    }
}
