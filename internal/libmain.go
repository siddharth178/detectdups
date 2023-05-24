package internal

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	BUFSIZE = int64(4096)
)

type filegroups [][]string

func (fgs filegroups) PrintStats(title string) {
	fmt.Printf("--%s--\n", title)
	dupData := map[string]int64{}
	for _, fg := range fgs {
		fmt.Println(fg)

		info, err := os.Stat(fg[0])
		if err != nil {
			log.Warnln("skipping file:", fg[0], "err:", err)
		} else {
			dupData[fg[0]] = info.Size() * int64(len(fg)-1)
		}
	}
	dupBytes := int64(0)
	for _, size := range dupData {
		dupBytes += size
	}
	fmt.Printf("--%s-- #files with dups: %v\n", title, len(fgs))
	fmt.Printf("--%s-- dupBytes: %v\n", title, dupBytes)

}

func getFileBufSum(path string, pass int) (string, bool, error) {
	lastBuf := false
	buf := make([]byte, BUFSIZE)
	file, err := os.Open(path)
	if err != nil {
		return "", false, fmt.Errorf("open error. %w", err)
	}
	defer file.Close()

	_, err = file.Seek(int64(pass)*BUFSIZE, io.SeekStart)
	if err != nil {
		return "", false, fmt.Errorf("seek error. %w", err)
	}

	r := bufio.NewReader(file)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", false, fmt.Errorf("read error. %w", err)
	}
	if n < len(buf) {
		// read last buff
		lastBuf = true
	}

	h := sha256.New()
	h.Write(buf[:n])
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, lastBuf, nil
}

func buildFileGroupsBySum(files []string, pass int) (filegroups, filegroups) {
	fg := map[string][]string{}
	lastfg := map[string][]string{}
	for _, f := range files {
		sum, lastBuf, err := getFileBufSum(f, pass)
		if err != nil {
			log.Warnln("skipping file:", f, "err:", err)
			continue
		}
		if lastBuf {
			lastfg[sum] = append(lastfg[sum], f)
		} else {
			fg[sum] = append(fg[sum], f)
		}
	}
	fgs := filegroups{}
	dupgs := filegroups{}
	for _, files := range fg {
		if len(files) > 1 {
			fgs = append(fgs, files)
		}
	}
	for _, files := range lastfg {
		if len(files) > 1 {
			dupgs = append(dupgs, files)
		}
	}

	return fgs, dupgs
}

func processFileGroups(fgs filegroups, dupgs filegroups, pass int) (filegroups, filegroups) {
	newfgs := filegroups{}
	for _, fg := range fgs {
		tmpfgs, tmpdupgs := buildFileGroupsBySum(fg, pass)
		newfgs = append(newfgs, tmpfgs...)
		dupgs = append(dupgs, tmpdupgs...)
	}
	return newfgs, dupgs
}

func buildFileGroupsBySize(dirPath string) filegroups {
	bySize := map[int64][]string{}
	filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			bySize[info.Size()] = append(bySize[info.Size()], path)
		}
		return nil
	})

	fgs := filegroups{}
	for _, paths := range bySize {
		fgs = append(fgs, paths)
	}
	return fgs
}

func LibMain(dirPath string) {
	fgs := buildFileGroupsBySize(dirPath)
	pass := 0
	dupgs := filegroups{}
	for {
		fgs, dupgs = processFileGroups(fgs, dupgs, pass)
		if len(fgs) == 0 {
			fmt.Println("processing done")
			dupgs.PrintStats("final dupgs")
			break
		}
		pass++
	}
}
