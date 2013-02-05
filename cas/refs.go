package cas

import (
	"io/ioutil"
)

func SetRef(s Setter, ref string, addr Addr) error {
	ref_addr := make(Addr, len(ref)+1)
	ref_addr[0] = byte(addr_kind__ref)
	copy(ref_addr[1:], ref)

	w, err := s.Set()
	if err != nil {
		return err
	}

	w.Write(addr)
	return w.Commit(ref_addr)
}

func GetRef(s Getter, ref string) (Addr, error) {
	ref_addr := make(Addr, len(ref)+1)
	ref_addr[0] = byte(addr_kind__ref)
	copy(ref_addr[1:], ref)

	r, err := s.Get(ref_addr)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	addr, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Addr(addr), nil
}
