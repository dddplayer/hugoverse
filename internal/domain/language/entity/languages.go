package entity

// Languages is a sortable list of language.
type Languages []*Language

func (l Languages) Len() int { return len(l) }
func (l Languages) Less(i, j int) bool {
	wi, wj := l[i].Weight, l[j].Weight

	if wi == wj {
		return l[i].Lang < l[j].Lang
	}

	return wj == 0 || wi < wj
}

func (l Languages) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
