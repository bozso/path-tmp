package tempfiles

import (
    "github.com/bozso/gotoolbox/path"
)

type Like interface {
    Get() (path.File, error)
    Put(path.File)
    Remove() error
}

func Create(l Like) (vf path.ValidFile, err error) {
    file, err := l.Get()
    if err != nil {
        return
    }

    tmpFile, err := file.Create()
    if err != nil {
        err = CreateFail{filePath: file.String(), err: err}
        return
    }

    vf, err = file.ToValidFile()
    if err != nil {
        return
    }
    err = tmpFile.Close()
    return
}


type WithExtension struct {
    Like
    extension string
}

func (w WithExtension) Get() (f path.File, err error) {
    f, err = w.Like.Get()
    if err != nil {
        return
    }

    ext := w.extension
    // if we already have the proper extension just return
    if ext == f.Ext() {
        return
    }

    withExt := f.AddExt(ext)

    err = withExt.Touch()
    if err != nil {
        return
    }

    f = withExt.ToFile()
    return
}

type KeepAlive struct {
    Like
}

func (_ KeepAlive) Remove() (err error) {
    return nil
}
