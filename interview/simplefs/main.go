package main

import (
	"errors"
	"fmt"
)

type Block struct {
	Data string
	Next *Block
}

type File struct {
	Head *Block
	Size int // number of blocks
}

type SimpleFS struct {
	BlockSize  int
	NumBlocks  int
	Head       *Block
	FreeBlocks int
	WriteHead  *Block
	Files      map[string]*File
}

func NewSimpleFS(numBlocks int) *SimpleFS {
	fs := &SimpleFS{
		BlockSize:  8,
		NumBlocks:  numBlocks,
		FreeBlocks: numBlocks,
		Files:      make(map[string]*File),
	}

	if numBlocks > 0 {
		fs.Head = &Block{}
		current := fs.Head
		for i := 1; i < numBlocks; i++ {
			current.Next = &Block{}
			current = current.Next
		}
	}

	fs.WriteHead = fs.Head
	return fs
}

// ---------- Milestone 1 ----------
func (fs *SimpleFS) Print() {
	curr := fs.Head
	out := ""
	for curr != nil {
		out += fmt.Sprintf("[%s]", curr.Data)
		if curr.Next != nil {
			out += " -> "
		}
		curr = curr.Next
	}
	fmt.Println(out)
}

// ---------- Milestone 2 + 3 ----------
func (fs *SimpleFS) Write(filename string, data string) error {
	if _, exists := fs.Files[filename]; exists {
		return errors.New("file already exists")
	}
	if len(data) == 0 {
		return errors.New("empty write")
	}
	if len(data)%fs.BlockSize != 0 {
		return errors.New("data length must be multiple of block size")
	}

	blocksNeeded := len(data) / fs.BlockSize
	if blocksNeeded > fs.FreeBlocks {
		return errors.New("not enough free blocks")
	}

	// Store file metadata
	fs.Files[filename] = &File{
		Head: fs.WriteHead,
		Size: blocksNeeded,
	}

	curr := fs.WriteHead
	for i := 0; i < len(data); i += fs.BlockSize {
		curr.Data = data[i : i+fs.BlockSize]
		curr = curr.Next
	}

	fs.WriteHead = curr
	fs.FreeBlocks -= blocksNeeded
	return nil
}

// ---------- Milestone 4 ----------
func (fs *SimpleFS) Read(filename string) (string, error) {
	file, ok := fs.Files[filename]
	if !ok {
		return "", errors.New("file does not exist")
	}

	curr := file.Head
	result := ""
	for i := 0; i < file.Size; i++ {
		result += curr.Data
		curr = curr.Next
	}
	return result, nil
}

// ---------- Milestone 5 ----------
func (fs *SimpleFS) Delete(filename string) error {
	file, ok := fs.Files[filename]
	if !ok {
		return errors.New("file does not exist")
	}

	curr := file.Head
	for i := 0; i < file.Size; i++ {
		curr.Data = ""
		curr = curr.Next
	}

	delete(fs.Files, filename)
	return nil
}

// ---------- Milestone 6 ----------
func (fs *SimpleFS) Defrag() {
	if fs.FreeBlocks == 0 {
		return
	}

	var usedHead, usedTail *Block
	var freeHead, freeTail *Block

	curr := fs.Head
	freeCount := 0

	for curr != nil {
		if curr.Data == "" {
			// empty block
			freeCount++
			if freeHead == nil {
				freeHead = curr
				freeTail = curr
			} else {
				freeTail.Next = curr
				freeTail = freeTail.Next
			}
		} else {
			// used block
			if usedHead == nil {
				usedHead = curr
				usedTail = curr
			} else {
				usedTail.Next = curr
				usedTail = usedTail.Next
			}
		}
		curr = curr.Next
	}

	// Combine lists
	if usedTail != nil {
		usedTail.Next = freeHead
	}
	if freeTail != nil {
		freeTail.Next = nil
	}

	// Update FS pointers
	if usedHead != nil {
		fs.Head = usedHead
	} else {
		fs.Head = freeHead
	}

	fs.WriteHead = freeHead
	fs.FreeBlocks = freeCount
}

// -------------------------------------------------------------

func main() {
	fs := NewSimpleFS(5)
	fs.Print()

	fs.Write("file1", "12345678abcdefgh")
	fs.Write("file2", "asdfasdf")
	fs.Print()

	data1, _ := fs.Read("file1")
	fmt.Println(data1)

	data2, _ := fs.Read("file2")
	fmt.Println(data2)

	fs.Delete("file1")
	fs.Print()

	fs.Defrag()
	fs.Print()

	fs.Write("file1", "12345678abcdefgh")
	fs.Print()

	fs.Write("file3", "ABCDEFGHABCDEFGH")
	fs.Print()
}
