package hw10programoptimization

import (
	"bufio"
	"io"
	"log"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

//go:generate easyjson -all stats.go
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	res := make(DomainStat)
	var user User

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return DomainStat{}, err
		}
		if strings.Contains(user.Email, domain) {
			k := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			res[k]++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res, nil
}
