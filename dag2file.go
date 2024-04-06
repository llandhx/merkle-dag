// package merkledag

// // Hash to file
// func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
// 	// 根据hash和path， 返回对应的文件, hash对应的类型是tree
// 	return nil
// }

package merkledag

import (
	"encoding/json"
)

// Decode 从 JSON 数据中解码对象
func (obj *Object) Decode(data []byte) error {
	return json.Unmarshal(data, obj)
}

// getObject 从 KVStore 中获取哈希对应的对象
func getObject(store KVStore, hash []byte, _ HashPool) (*Object, error) {
	// 根据哈希值从 KVStore 中获取数据
	data, err := store.Get(hash)
	if err != nil {
		// 处理获取数据失败的情况
		return nil, err
	}

	// 解码数据为 Object 对象
	var obj Object
	if err := obj.Decode(data); err != nil {
		// 处理解码数据失败的情况
		return nil, err
	}

	return &obj, nil
}

// Hash2File 根据哈希值和路径返回对应的文件内容
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据哈希值从 KVStore 中获取对象
	obj, err := getObject(store, hash, hp)
	if err != nil {
		// 处理获取对象失败的情况
		// 这里可以根据实际需求返回错误信息或默认值
		return nil
	}

	// 如果对象为文件，则直接返回文件内容
	if len(path) == 0 || path == "/" {
		// 如果路径为空或根路径，则返回对象的数据字段
		return obj.Data
	}

	// 如果对象为文件夹，则根据路径逐级查找对应的文件或子文件夹
	// 对于文件夹，我们需要遍历其链接，找到与路径匹配的对象
	for _, link := range obj.Links {
		if link.Name == path {
			// 如果找到了与路径匹配的链接，递归调用 Hash2File 函数获取对应的文件内容
			return Hash2File(store, link.Hash, "", hp)
		}
	}

	// 如果在文件夹中未找到与路径匹配的对象，则返回 nil（或其他默认值）
	return nil
}
