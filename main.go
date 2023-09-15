package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MerkeTreeNode struct {
	name     string
	data     []byte
	children []MerkeTreeNode
}

type MerkeTree struct {
	root MerkeTreeNode
}

func main() {
    dirname := flag.String("dir", "./", "Directory to hash")

    flag.Parse()

	tree, err := NewMerkelTree(*dirname)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", tree.HashHex())
}

func (n *MerkeTreeNode) IsLeaf() bool {
	return len(n.children) <= 0
}

func (n *MerkeTreeNode) Hash() [32]byte {
    data := []byte(n.name)

    if n.IsLeaf() {
        data = append(data, n.data...)
        return sha256.Sum256(data)
    }

    for _, child := range n.children {
        hash := child.Hash()
        data = append(data, hash[:]...)
    }

    return sha256.Sum256(data) 
}

func (t *MerkeTree) Hash() [32]byte {
    return t.root.Hash()
}

func (t *MerkeTree) HashHex() string {
    hash := t.Hash()
    return hex.EncodeToString(hash[:])
}

func GetDirNodes(dirname string) ([]MerkeTreeNode, error) {
	entries, err := os.ReadDir(dirname)
	nodes := []MerkeTreeNode{}

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		filename := filepath.Join(dirname, entry.Name())

		if !entry.IsDir() {
			data, err := os.ReadFile(filename)

			if err != nil {
				return nil, err
			}

			nodes = append(nodes, MerkeTreeNode{name: filename, data: data})
			continue
		}

		children, err := GetDirNodes(filename)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, MerkeTreeNode{name: filename, children: children})
	}

	return nodes, nil
}

func NewMerkelTree(dirname string) (*MerkeTree, error) {
	children, err := GetDirNodes(dirname)
	tree := MerkeTree{}

	if err != nil {
		return nil, err
	}

	tree.root = MerkeTreeNode{}
	tree.root.children = children

	return &tree, nil
}
