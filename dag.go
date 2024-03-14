package merkledag

import (
	"encoding/binary"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

// Add 将 Node 中的数据保存在 KVStore 中，并返回 Merkle Root
func Add(store KVStore, node Node, h hash.Hash) []byte {

	// TODO 将分片写入到KVStore中，并返回Merkle Root

	// 初始化一个 Merkle Tree 的叶子节点列表
	leafNodes := make([][]byte, 0)

	// 递归遍历 Node 对象，将每个节点的数据存储在 KVStore 中，并将其哈希值添加到 leafNodes 列表中
	var traverseNode func(node Node)
	traverseNode = func(node Node) {
		switch n := node.(type) {
		case File:
			// 将文件内容存储在 KVStore 中
			store.Put([]byte(n.Name()), n.Bytes())
			// 计算文件内容的哈希值并添加到 leafNodes 列表中
			h.Reset()
			h.Write(n.Bytes())
			hashValue := h.Sum(nil)
			leafNodes = append(leafNodes, hashValue)
		case Dir:
			// 遍历文件夹中的每个文件/文件夹
			it := n.It()
			for it.Next() {
				childNode := it.Node()
				traverseNode(childNode)
			}
		}
	}

	// 开始递归遍历 Node 对象
	traverseNode(node)

	// 构建 Merkle Tree
	merkleRoot := buildMerkleTree(leafNodes, h)

	return merkleRoot
}

// buildMerkleTree 构建 Merkle Tree 并返回 Merkle Root
func buildMerkleTree(leafNodes [][]byte, h hash.Hash) []byte {
	// 如果没有叶子节点，直接返回 nil
	if len(leafNodes) == 0 {
		return nil
	}

	// 按照哈希值的字典序排序叶子节点列表
	sortByteSlices(leafNodes)

	// 如果叶子节点数量为奇数，则复制最后一个节点以使其数量变为偶数
	if len(leafNodes)%2 != 0 {
		lastNode := leafNodes[len(leafNodes)-1]
		leafNodes = append(leafNodes, lastNode)
	}

	// 逐层构建 Merkle Tree
	for len(leafNodes) > 1 {
		// 每次循环处理一层
		var nextLevel [][]byte

		// 对每一对叶子节点计算其父节点的哈希值
		for i := 0; i < len(leafNodes); i += 2 {
			// 合并两个叶子节点的哈希值
			h.Reset()
			h.Write(leafNodes[i])
			h.Write(leafNodes[i+1])
			parentHash := h.Sum(nil)
			nextLevel = append(nextLevel, parentHash)
		}

		// 将下一层的节点列表作为当前层继续循环处理
		leafNodes = nextLevel
	}

	// 最终 Merkle Root 即为根节点的哈希值
	return leafNodes[0]
}

// sortByteSlices 对字节片段进行字典序排序
func sortByteSlices(slices [][]byte) {
	for i := 0; i < len(slices)-1; i++ {
		for j := 0; j < len(slices)-i-1; j++ {
			if binary.LittleEndian.Uint64(slices[j]) > binary.LittleEndian.Uint64(slices[j+1]) {
				slices[j], slices[j+1] = slices[j+1], slices[j]
			}
		}
	}
}
