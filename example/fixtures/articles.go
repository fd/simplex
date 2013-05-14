package main

import (
	"encoding/json"
	"fmt"
	"time"
)

const C = 50

func main() {

	type article struct {
		Id          int64
		Title       string
		Body        string
		PublishedAt time.Time
	}

	var (
		o []interface{}
		n = time.Now()
	)

	for i := 0; i < C; i++ {
		o = append(o, &article{
			Id:          int64(i + 1),
			Title:       fmt.Sprintf("Article %d", i+1),
			Body:        LOREM,
			PublishedAt: n.AddDate(0, 0, 20-i),
		})
	}

	data, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(data))
}

const LOREM = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Suspendisse bibendum lacus a orci accumsan congue. Nulla ac quam diam. Nullam
quis elit libero, vitae sollicitudin turpis. Proin mattis augue neque. Fusce
rhoncus tortor a eros aliquet ultricies. Phasellus eu est eu quam posuere porta
a ut quam. Praesent vel ante massa, ut pretium neque. Aenean eleifend viverra
elit, nec pellentesque ligula rutrum id. Nulla facilisi. Nullam cursus quam
quam. Maecenas sagittis, magna ac volutpat fringilla, ipsum lectus gravida
dolor, et facilisis enim mauris et metus. Etiam sed nibh mauris, sit amet
consectetur lorem.`
