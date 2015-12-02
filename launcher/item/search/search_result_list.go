package search

// ResultList is type alias for []Result used for sort.
type ResultList []Result

func (self ResultList) Len() int {
	return len(self)
}

func (self ResultList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self ResultList) Less(i, j int) bool {
	if self[i].Score > self[j].Score {
		return true
	}

	if self[j].Score > self[i].Score {
		return false
	}

	if self[i].Freq > self[j].Freq {
		return true
	}

	if self[j].Freq > self[i].Freq {
		return false
	}

	return self[i].Name < self[j].Name
}
