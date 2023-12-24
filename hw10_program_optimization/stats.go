package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	json "github.com/goccy/go-json"
)

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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	var user User

	i := 0
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			return result, fmt.Errorf("unmarshal error: %w", err)
		}
		result[i] = user
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	lowerDomain := strings.ToLower(domain)

	for _, user := range u {
		lowerEmail := strings.ToLower(user.Email)
		atIndex := strings.LastIndex(lowerEmail, "@")

		if atIndex != -1 && strings.HasSuffix(lowerEmail[atIndex+1:], lowerDomain) {
			domainPart := lowerEmail[atIndex+1:]
			result[domainPart]++
		}
	}

	return result, nil
}
