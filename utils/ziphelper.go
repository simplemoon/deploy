package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CheckFunc = func(p string) bool

// 压缩文件
func Compress(source, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		p := filepath.Join(source, path)
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		return compress(f, info, "", w)
	})
	if err != nil {
		return err
	}
	return nil
}

// 压缩对应的文件
func compress(file *os.File, info os.FileInfo, prefix string, zw *zip.Writer) error {
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, fi, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func DeCompress(zipFile, dest string, fn CheckFunc) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	cfList := make([]io.Closer, 0)
	cfList = append(cfList, reader)
	defer func() {
		for _, c := range cfList {
			c.Close()
		}
	}()

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		// 检查一下是否需要拷贝啊
		if fn != nil && !fn(file.Name) {
			continue
		}
		cfList = append(cfList, rc)
		filename := filepath.Join(dest, file.Name)
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		cfList = append(cfList, w)
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
