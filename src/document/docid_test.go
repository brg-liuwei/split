package document

import (
	"math/rand"
	"os"
	"testing"
)

func TestDocId(t *testing.T) {
	DocInit()
	for i := 0; i != 25; i++ {
		id := AllocId()
		if id != int64(i) {
			t.Error("id: ", id, ", i: ", i)
		}
	}

	for i := 0; i != 25; i++ {
		if ActiveId(int64(i)) == false {
			t.Error("id: ", i, " disactive")
		}
	}

	delMap := make(map[int]bool)
	_ = rand.Int()
	for i := 0; i != 30; i++ {
		//if rand.Int()&0x1 != 0 {
		if i&0x1 == 1 {
			delMap[i] = true
		}
	}
	for k, _ := range delMap {
		DeleteId(int64(k))
	}

	for i := 0; i != 25; i++ {
		bl := ActiveId(int64(i))
		_, ok := delMap[i]
		if ok {
			if bl == true {
				t.Error("id ", i, " is deleted")
			}
		} else {
			if bl == false {
				t.Error("id ", i, " is not deleted")
			}
		}
	}

	for i := 25; i != 30; i++ {
		if ActiveId(int64(i)) {
			t.Error("id ", i, " is not Exist")
		}
	}
	err := os.Remove(delfile)
	if err != nil {
		t.Error(err)
	}
	err = os.RemoveAll(dbname)
	if err != nil {
		t.Error(err)
	}
}
