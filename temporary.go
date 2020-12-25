package tempfiles

import (
    "fmt"
    "sync"
    "time"

    "github.com/bozso/emath/rand"
    "github.com/bozso/gotoolbox/path"
)

const (
    InUse bool = true
    NotInUse bool = false
)

// A map that describes whether a specific file is in use or not
type Set map[*path.ValidFile]bool

// The main struct for managing temporary files in a single directory.
type Files struct {
    // Path to the root directory
    RootDir path.Dir
    // file set pointing to existing files
    files Set
    // random number generator for generating random file names
    rand.Rand
}

/*
Set up temporary file management for the specified directory
with the given randum number generator.
*/
func FromDir(rootDir path.Dir, rng rand.Rand) (t Files) {
    return Files{
        RootDir: rootDir,
        files: make(Set),
        Rand: rng,
    }
}

// Set up temporary file management with default parameters. Random
// number generator will be produced using the current unix time stamp.
func NewDefault() (f Files, err error) {
    src := rand.NewSource(time.Now().Unix())
    f, err = FromRand(rand.NoScale(src))
    return
}

// Set up temporary file management with the specified random number
// generator. Directory name will be randomly generated using the
// generator.
func FromRand(rng rand.Rand) (f Files, err error) {
    prefix := fmt.Sprintf("%d", rng.Int())

    f, err = New("", prefix, rng)
    return
}

func New(dir, prefix string, rng rand.Rand) (f Files, err error) {
    d, err := path.TempDirIn(dir, prefix)
    if err != nil {
        return
    }

    f = FromDir(d, rng)
    return
}

// Convert it to mutex guarded temporary file manager.
func (f Files) Mutexed() (m Mutexed) {
    return Mutexed{
        files: f,
    }
}

/*
Search for a valid file that is not in use managed by the receiver.
The second return argument marks whether a file that is not in use was
found.
*/
func (f *Files) Search() (vf *path.ValidFile, found bool) {
    for file, inUse := range f.files {
        if !inUse {
            vf, f.files[file], found = file, InUse, true
            break
        }
    }
    return
}

/*
Retreives a new temporary file to be used.
First it searches for a file that is not in use. If no such file is
found a new file will be created and registered in the receivers
fileset.
*/
func (f *Files) Get() (vf *path.ValidFile, err error) {
    vf, found := f.Search()
    if found {
        return
    }
    vf, err = f.NewFile()
    return
}

/*
Creates a new file to be used in the temporary file directory. Returns
error if file creation has failed.
*/
func (f *Files) NewFile() (vf *path.ValidFile, err error) {
    file := f.RootDir.Join(fmt.Sprintf("%d", f.Rand.Int()))

    _, err = file.Create()
    if err != nil {
        err = CreateFail{filePath: file.String(), err: err}
        return
    }

    vfile, err := file.ToValidFile()
    if err != nil {
        return
    }

    vf = &vfile

    return
}

/*
Signals to the receiver that the temporary file is no longer in use.
Should be used in conjunction with Get.

    var t = NewDefaultTempFiles()
    f, err := t.Get()
    if err != nil {
        // error handling
    }
    defer t.Put(f)
    // use f
*/
func (f *Files) Put(vf *path.ValidFile) {
    f.files[vf] = NotInUse
}

/*
Remove removes the temporary directory containing the temporary files.
*/
func (f *Files) Remove() (err error) {
    err = f.RootDir.Remove()
    return
}

// CreateFail describes file creation failure.
type CreateFail struct {
    filePath string
    err error
}

func (e CreateFail) Error() (s string) {
    s = fmt.Sprintf("failed to create temporary file '%s'", e.filePath)
    return
}

func (e CreateFail) Unwrap() (err error) {
    return e.err
}

// Concurrent safe TempFiles, guarded by mutex
type Mutexed struct {
    // The wrapped struct.
    files Files
    // mutex for protecting the locking the set
    mutex sync.Mutex
}

// Concurrent safe Get
func (m *Mutexed) Get() (vf *path.ValidFile, err error) {
    m.mutex.Lock()
    vf, err = m.files.Get()
    m.mutex.Unlock()
    return
}

// Concurrent safe Search
func (m *Mutexed) Search() (vf *path.ValidFile, found bool) {
    m.mutex.Lock()
    vf, found = m.files.Search()
    m.mutex.Unlock()
    return
}

// Concurrent safe NewFile
func (m *Mutexed) NewFile() (vf *path.ValidFile, err error) {
    m.mutex.Lock()
    vf, err = m.files.NewFile()
    m.mutex.Unlock()
    return
}

// Concurrent safe Put
func (m *Mutexed) Put(vf *path.ValidFile) {
    m.mutex.Lock()
    m.files.Put(vf)
    m.mutex.Unlock()
    return
}
