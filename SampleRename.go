// prog to rename a batch of sample files 
// into the format 'md5.ext', 
// while 'md5' is the md5 sum of the file 
// and 'ext' beeing the filetype extension.

package main

import("path/filepath")
import("encoding/hex")
import("crypto/md5")
import("io/ioutil")
import("bufio")
import("log")
import("os")
import("io")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func pathValidator(ipath string) (opath string) {
	//log.Println("pathValidator ipath: ", ipath)
	opath, err := filepath.Abs(ipath)
	check(err)
	//log.Println("pathValidator opath: ", opath)
	return opath
}

func readFile(fname string) (data []byte) {
		fname = pathValidator(fname)
		f, err := os.Open(fname)
		check(err)
		defer f.Close()
		rdr := bufio.NewReader(f)
		data, err = ioutil.ReadAll(rdr)
		check(err)
		return data
}

func genMd5(data []byte) (md5str string) {
		md5 := md5.New()
		mw := io.MultiWriter(md5)
		mw.Write(data)
		md5sum := md5.Sum(nil)
		md5str = hex.EncodeToString(md5sum)
		return md5str
}

func genOutputFilename(ifile string) (ofile string) {
		v, err := filepath.Abs(ifile)
		check(err)
		data := readFile(v)
		filename := genMd5(data)
		ext := filepath.Ext(v)
		dir := filepath.Dir(v)
		ofile = filepath.Join(dir, filename)
		ofile = ofile+ext
		//log.Println("OUTFILE: ", ofile)
		return ofile
}

func main() {
	args := os.Args[1:]
	m := make(map[string]string)
	for _, v := range args {
		filename := genOutputFilename(v)
		m[filename] = v
		err := os.Rename(v, filename)
		check(err)
		basepath, err := os.Getwd()
		check(err)
		relpath, err := filepath.Rel(basepath, filename)
		check(err)
		log.Printf("renamed '%s' -> '%s'\n", v, relpath)
	}
	//log.Println("map:", m)
}
