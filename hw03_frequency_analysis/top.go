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
	sort.SliceStable(keys, func(i, j int) bool {
		return mp[keys[i]] > mp[keys[j]]
	})

	var sk []string
	n := 10
	for i := 0; i < len(keys); i++ {
		sk = append(sk, keys[i])
		if i >= n-1 {
			break
		}
	}

	var prev int
	var start int
	for i, k := range sk {
		if i != 0 && prev != mp[k] {
			sort.Strings(sk[start:i])
			start = i
		}
		prev = mp[k]
	}
	sort.Strings(sk[start:])

	return sk
}
