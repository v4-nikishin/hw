package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(t string) []string {
	if t == "" {
		return nil
	}

	sl := strings.Fields((t))
	mp := make(map[string]int)
	for _, v := range sl {
		mp[v]++
	}

	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if mp[keys[i]] == mp[keys[j]] {
			return keys[i] < keys[j]
		}
		return mp[keys[i]] > mp[keys[j]]
	})

	top := 10
	if len(keys) < top {
		top = len(keys)
	}
	return keys[:top]
}
