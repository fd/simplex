package router

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"simplex.sh/store/cas"
	"sort"
)

type Rule struct {
	Key         cas.Addr
	Path        string
	Host        string
	Language    string
	ContentType string
	Status      int
	Header      http.Header
	Address     cas.Addr
}

func (r *Rule) calculate_key() {
	var (
		headers = make([]string, 0, len(r.Header))
		sha     = sha1.New()
	)

	fmt.Fprintf(sha, "%s %s %s %s %d %s",
		r.Path,
		r.Host,
		r.Language,
		r.ContentType,
		r.Status,
		r.Address.String(),
	)

	for h := range r.Header {
		headers = append(headers, h)
	}

	sort.Strings(headers)

	for _, h := range headers {
		fmt.Fprintf(sha, "%s %v", h, r.Header[h])
	}

	key := sha.Sum(nil)
	r.Key = cas.Addr(key)
}
