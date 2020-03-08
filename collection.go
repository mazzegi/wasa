package wasa

func (e *Elt) replaceAllChilds(cs []*Elt) {
	e.RemoveAll()
	e.Childs = cs
	e.Invalidate()
}

func (e *Elt) ReplaceCollection(new *Elt) {
	e.replaceAllChilds(new.Childs)
	// if len(e.Childs) != len(new.Childs) {
	// 	e.replaceAllChilds(new.Childs)
	// 	return
	// }

	// for i, c := range e.Childs {
	// 	nc := new.Childs[i]
	// 	if c.key == "" || c.key != nc.key {
	// 		e.replaceAllChilds(new.Childs)
	// 		return
	// 	}
	// 	if !bytes.Equal(c.Hash(), nc.Hash()) {
	// 		e.Replace(c, nc)
	// 		nc.Invalidate()
	// 	}
	// }
}
