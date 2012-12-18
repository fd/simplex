package trie

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	trie := New()
	var k, v string

	trie.Insert([]byte("foo"), "a")
	trie.Insert([]byte("foo"), "b")
	trie.Insert([]byte("foos"), "c")
	trie.Insert([]byte("foor"), "f")
	trie.Insert([]byte("foe"), "d")
	trie.Insert([]byte("fo"), "e")
	trie.Insert([]byte("f"), "g")

	fmt.Printf("trie: %v\n", trie)

	k, v = "foo", "b"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foos", "c"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foor", "f"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foe", "d"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "fo", "e"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "f", "g"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("Was supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("Was supposed to find a string")
	} else if s != v {
		t.Errorf("Was supposed to find `%s` instead of %+v", v, val)
	}
}

func TestWordList(t *testing.T) {

	words := words()
	trie := New()

	for i, word := range words {
		trie.Insert(word, i)
	}

	for i, word := range words {
		if val, f := trie.Lookup(word); !f {
			t.Fatalf("Was supposed to find `%s`", word)
		} else if s, ok := val.(int); !ok {
			t.Fatalf("Was supposed to find a int but was a %T", val)
		} else if s != i {
			t.Fatalf("Was supposed to find `%d` instead of %+v", i, val)
		}
	}

	m := trie.ConsumedMemory()
	fmt.Printf("mem: %f (%f / N)\n", float64(m)/1024/1024, float64(m)/float64(len(words)))
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	words := words()
	b.N = len(words)

	b.StartTimer()

	trie := New()

	for i, word := range words {
		trie.Insert(word, i)
	}
}

func BenchmarkLookup(b *testing.B) {
	b.StopTimer()

	words := words()
	b.N = len(words)

	trie := New()

	for i, word := range words {
		trie.Insert(word, i)
	}

	b.StartTimer()

	for _, word := range words {
		trie.Lookup(word)
	}
}

func BenchmarkRemove(b *testing.B) {
	b.StopTimer()

	words := words()
	b.N = len(words)

	trie := New()

	for i, word := range words {
		trie.Insert(word, i)
	}

	b.StartTimer()

	for _, word := range words {
		trie.Remove(word)
	}
	b.StopTimer()

	time.Sleep(30 * time.Second)
}

var v_words [][]byte

func words() [][]byte {
	if len(v_words) > 0 {
		return v_words
	}

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		panic(err)
	}

	v_words := bytes.Split(w, []byte{' '})

	return v_words
}
