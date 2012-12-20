package trie

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"testing"
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
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foos", "c"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foor", "f"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "foe", "d"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "fo", "e"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	k, v = "f", "g"
	if val, f := trie.Lookup([]byte(k)); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(string); !ok {
		t.Errorf("is supposed to find a string")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}
}

func Test_Case_1A(t *testing.T) {
	trie := New()
	k := []byte("")
	v := 1

	val, f := trie.Insert(k, v)
	if !f {
		t.Errorf("is supposed to insert (`%s`, %v)", k, v)
	}
	if val != nil {
		t.Errorf("no value was set for `%s` so nil is expected instead of %v", k, val)
	}

	val, f = trie.Insert(k, 2)
	if !f {
		t.Errorf("is supposed to insert (`%s`, %v)", k, 2)
	}
	if a, ok := val.(int); ok && a != v {
		t.Errorf("1 was set for `%s` so 1 is expected instead of %v", k, val)
	}
}

func Test_Case_1B(t *testing.T) {
	trie := New()
	k := []byte("")
	v := 1

	trie.Insert(k, v)

	if val, f := trie.Lookup(k); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(int); !ok {
		t.Errorf("is supposed to find an int")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}
}

func Test_Case_1C(t *testing.T) {
	trie := New()
	k := []byte("")
	v := 1

	trie.Insert(k, v)

	if val, f := trie.Remove(k); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(int); !ok {
		t.Errorf("is supposed to find an int")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
	}

	if val, f := trie.Lookup(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_2(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)

	trie_str := `{
  K: ,
  V: <nil>,
  { K: foo, V: 1 }
}`

	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}

	k = []byte("foobar")
	v = 2

	trie.Insert(k, v)

	trie_str = `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: 1,
    { K: bar, V: 2 }
  }
}`

	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_3A(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	if val, f := trie.Lookup(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}

	trie.Insert(k, v)
	k = []byte("foobar")

	if val, f := trie.Lookup(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_3B(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	if val, f := trie.Remove(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}

	trie.Insert(k, v)
	k = []byte("foobar")

	if val, f := trie.Remove(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_4(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)

	k = []byte("fo")
	v = 2
	trie_str := `{
  K: ,
  V: <nil>,
  {
    K: fo,
    V: 2,
    { K: o, V: 1 }
  }
}`
	trie.Insert(k, v)

	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_5A(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)
	k = []byte("fo")

	if val, f := trie.Lookup(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_5B(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)
	k = []byte("fo")

	if val, f := trie.Remove(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_6(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)

	k = []byte("foe")
	v = 2
	trie_str := `{
  K: ,
  V: <nil>,
  {
    K: fo,
    V: <nil>,
    { K: e, V: 2 }
    { K: o, V: 1 }
  }
}`
	trie.Insert(k, v)

	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_7A(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)
	k = []byte("foe")

	if val, f := trie.Lookup(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_7B(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)
	k = []byte("foe")

	if val, f := trie.Remove(k); f {
		t.Errorf("is supposed to NOT find `%s`", k)
	} else if val != nil {
		t.Errorf("is supposed to be nil")
	}
}

func Test_Case_8(t *testing.T) {
	trie := New()
	k := []byte("foo")
	v := 1

	trie.Insert(k, v)

	k = []byte("foo")
	v = 2
	trie_str := `{
  K: ,
  V: <nil>,
  { K: foo, V: 2 }
}`
	trie.Insert(k, v)

	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_9A(t *testing.T) {
	trie := New()

	trie.Insert([]byte("foo"), 1)
	trie.Insert([]byte("foobar"), 2)
	trie.Insert([]byte("foogar"), 3)

	trie_str := `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: 1,
    { K: bar, V: 2 }
    { K: gar, V: 3 }
  }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}

	trie.Remove([]byte("foo"))
	trie_str = `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: <nil>,
    { K: bar, V: 2 }
    { K: gar, V: 3 }
  }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_9B(t *testing.T) {
	trie := New()

	trie.Insert([]byte("foo"), 1)
	trie.Insert([]byte("foobar"), 2)
	trie.Insert([]byte("foogar"), 3)

	trie_str := `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: 1,
    { K: bar, V: 2 }
    { K: gar, V: 3 }
  }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}

	trie.Remove([]byte("foobar"))
	trie_str = `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: 1,
    { K: gar, V: 3 }
  }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_9C(t *testing.T) {
	trie := New()

	trie.Insert([]byte("foo"), 1)
	trie.Insert([]byte("foobar"), 2)

	trie_str := `{
  K: ,
  V: <nil>,
  {
    K: foo,
    V: 1,
    { K: bar, V: 2 }
  }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}

	trie.Remove([]byte("foo"))
	trie_str = `{
  K: ,
  V: <nil>,
  { K: foobar, V: 2 }
}`
	if trie_str != trie.String() {
		t.Errorf("expected:\n%s\nactual:\n%s", trie_str, trie)
	}
}

func Test_Case_10(t *testing.T) {
	trie := New()
	k := []byte("foobar")
	v := 2

	trie.Insert([]byte("foo"), 1)
	trie.Insert([]byte("foobar"), 2)

	if val, f := trie.Lookup(k); !f {
		t.Errorf("is supposed to find `%s`", k)
	} else if s, ok := val.(int); !ok {
		t.Errorf("is supposed to find an int")
	} else if s != v {
		t.Errorf("is supposed to find `%s` instead of %+v", v, val)
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
			t.Fatalf("is supposed to find `%s`", word)
		} else if s, ok := val.(int); !ok {
			t.Fatalf("is supposed to find a int but was a %T", val)
		} else if s != i {
			t.Fatalf("is supposed to find `%d` instead of %+v", i, val)
		}
	}

	m := trie.ConsumedMemory()
	fmt.Printf("mem: %f (%f / N)\n", float64(m)/1024/1024, float64(m)/float64(len(words)))
}

func TestGobEncodeDecodeRoundtripSmall(t *testing.T) {

	trie_a := New()
	var trie_b *T

	trie_a.Insert([]byte("foo"), "a")
	trie_a.Insert([]byte("foo"), "b")
	trie_a.Insert([]byte("foos"), "c")
	trie_a.Insert([]byte("foor"), "f")
	trie_a.Insert([]byte("foe"), "d")
	trie_a.Insert([]byte("fo"), "e")
	trie_a.Insert([]byte("f"), "g")

	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(trie_a)
	if err != nil {
		t.Fatalf("err: encode: %s", err)
	}

	fmt.Printf("trie: %f\n", float64(len(buf.Bytes()))/1024/1024)

	err = gob.NewDecoder(&buf).Decode(&trie_b)
	if err != nil {
		t.Fatalf("err: decode: %s", err)
	}
}

func TestGobEncodeDecodeRoundtrip(t *testing.T) {

	words := words()
	trie_a := New()
	var trie_b *T

	for i, word := range words {
		trie_a.Insert(word, i)
	}

	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(trie_a)
	if err != nil {
		t.Fatalf("err: encode: %s", err)
	}

	fmt.Printf("list: %f <=> trie: %f\n", float64(v_words_len)/1024/1024, float64(len(buf.Bytes()))/1024/1024)

	err = gob.NewDecoder(&buf).Decode(&trie_b)
	if err != nil {
		t.Fatalf("err: decode: %s", err)
	}

	for i, word := range words {
		if val, f := trie_b.Lookup(word); !f {
			t.Fatalf("is supposed to find `%s`", word)
		} else if s, ok := val.(int); !ok {
			t.Fatalf("is supposed to find a int but was a %T", val)
		} else if s != i {
			t.Fatalf("is supposed to find `%d` instead of %+v", i, val)
		}
	}

}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	words := words()

	n := 20
	b.N = len(words) * n
	tries := make([]*T, n)
	for i := 0; i < n; i++ {
		tries[i] = New()
	}

	b.StartTimer()

	for _, trie := range tries {
		for i, word := range words {
			v, f := trie.Insert(word, i)
			if !f {
				fmt.Printf("v: %+v f: %+v\n", v, f)
			}
		}
	}
}

func BenchmarkInsertExisting(b *testing.B) {
	b.StopTimer()

	words := words()

	n := 20
	b.N = len(words) * n
	tries := make([]*T, n)
	for i := 0; i < n; i++ {
		trie := New()
		tries[i] = trie
		for i, word := range words {
			trie.Insert(word, i)
		}
	}

	b.StartTimer()

	for _, trie := range tries {
		for i, word := range words {
			v, f := trie.Insert(word, i)
			if !f {
				fmt.Printf("v: %+v f: %+v\n", v, f)
			}
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	b.StopTimer()

	words := words()

	n := 20
	b.N = len(words) * n
	tries := make([]*T, n)
	for i := 0; i < n; i++ {
		trie := New()
		tries[i] = trie
		for i, word := range words {
			trie.Insert(word, i)
		}
	}

	b.StartTimer()

	for _, trie := range tries {
		for _, word := range words {
			v, f := trie.Lookup(word)
			if !f {
				fmt.Printf("v: %+v f: %+v\n", v, f)
			}
		}
	}
}

func BenchmarkRemove(b *testing.B) {
	b.StopTimer()

	words := words()

	n := 20
	b.N = len(words) * n
	tries := make([]*T, n)
	for i := 0; i < n; i++ {
		trie := New()
		tries[i] = trie
		for i, word := range words {
			trie.Insert(word, i)
		}
	}

	b.StartTimer()

	for _, trie := range tries {
		for _, word := range words {
			trie.Remove(word)
		}
	}
}

var v_words [][]byte
var v_words_len int

func words() [][]byte {
	if len(v_words) > 0 {
		return v_words
	}

	w, err := ioutil.ReadFile("word_list.txt")
	if err != nil {
		panic(err)
	}

	v_words_len = len(w)
	v_words = bytes.Split(w, []byte{' '})

	return v_words
}
