package trie

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestInsert(t *testing.T) {
	trie := T{}
	var k, v string

	trie.Insert([]byte("foo"), "a")
	trie.Insert([]byte("foo"), "b")
	trie.Insert([]byte("foos"), "c")
	trie.Insert([]byte("foor"), "f")
	trie.Insert([]byte("foe"), "d")
	trie.Insert([]byte("fo"), "e")

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
}

func TestWordList(t *testing.T) {

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		t.Fatal(err)
	}

	words := bytes.Split(w, []byte{'\n'})
	trie := T{}

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
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		b.Fatal(err)
	}

	words := bytes.Split(w, []byte{'\n'})
	b.N = len(words)

	b.StartTimer()

	trie := T{}

	for i, word := range words {
		trie.Insert(word, i)
	}
}

func BenchmarkLookup(b *testing.B) {
	b.StopTimer()

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		b.Fatal(err)
	}

	words := bytes.Split(w, []byte{'\n'})
	b.N = len(words)

	trie := T{}

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

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		b.Fatal(err)
	}

	words := bytes.Split(w, []byte{'\n'})
	b.N = len(words)

	trie := T{}

	for i, word := range words {
		trie.Insert(word, i)
	}

	b.StartTimer()

	for _, word := range words {
		trie.Remove(word)
	}
}
