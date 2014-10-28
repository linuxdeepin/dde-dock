package search

type SearchResultList []SearchResult

func (self SearchResultList) Len() int {
	return len(self)
}

func (self SearchResultList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self SearchResultList) Less(i, j int) bool {
	if self[i].Score > self[j].Score {
		return true
	}

	if self[j].Score > self[i].Score {
		return false
	}

	return self[i].Name < self[j].Name
}
