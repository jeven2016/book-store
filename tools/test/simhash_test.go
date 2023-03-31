package test

import (
	"fmt"
	"github.com/go-dedup/simhash"
	"github.com/go-dedup/simhash/simhashCJK"
	"testing"
)

func TestChinese(t *testing.T) {
	var docs = [][]byte{
		[]byte("当山峰没有棱角的时候"),
		[]byte("当山谷没有棱角的时候"),
		[]byte("棱角的时候"),
		[]byte("你妈妈喊你回家吃饭哦，回家罗回家罗"),
		[]byte("你妈妈叫你回家吃饭啦，回家罗回家罗"),
	}

	hashes := make([]uint64, len(docs))
	sh := simhashCJK.NewSimhash()
	for i, d := range docs {
		fs := sh.NewWordFeatureSet(d)
		// fmt.Printf("%#v\n", fs)
		// actual := fs.GetFeatures()
		// fmt.Printf("%#v\n", actual)
		hashes[i] = sh.GetSimhash(fs)
		fmt.Printf("Simhash of '%s': %x\n", d, hashes[i])
	}

	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], simhash.Compare(hashes[0], hashes[1]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], simhash.Compare(hashes[0], hashes[2]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[3], simhash.Compare(hashes[0], hashes[3]))

}

func TestEnglish(t *testing.T) {
	var docs = [][]byte{
		[]byte("this is a test phrase"),
		[]byte("this is a test phrass"),
		[]byte("these are test phrases"),
		[]byte("foo bar"),
	}

	hashes := make([]uint64, len(docs))
	sh := simhash.NewSimhash()
	for i, d := range docs {
		hashes[i] = sh.GetSimhash(sh.NewWordFeatureSet(d))
		fmt.Printf("Simhash of '%s': %x\n", d, hashes[i])
	}

	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], simhash.Compare(hashes[0], hashes[1]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], simhash.Compare(hashes[0], hashes[2]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[3], simhash.Compare(hashes[0], hashes[3]))

	// Output:
	// Simhash of 'this is a test phrase': 8c3a5f7e9ecb3f35
	// Simhash of 'this is a test phrass': 8c3a5f7e9ecb3f21
	// Simhash of 'these are test phrases': ddfdbf7fbfaffb1d
	// Simhash of 'foo bar': d8dbe7186bad3db3
	// Comparison of `this is a test phrase` and `this is a test phrass`: 2
	// Comparison of `this is a test phrase` and `these are test phrases`: 22
	// Comparison of `this is a test phrase` and `foo bar`: 29
}
