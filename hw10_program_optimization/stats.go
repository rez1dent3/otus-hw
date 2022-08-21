package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/goccy/go-json"
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

type users chan User

func getUsers(r io.Reader) (users, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result := make(users)
	var user User
	go func() {
		defer close(result)
		for i := 0; ; i++ {
			if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
				return
			}

			result <- user

			scanner.Scan()
			if scanner.Err() != nil {
				return
			}
		}
	}()

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	regular, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	var key string
	for user := range u {
		if regular.Match([]byte(user.Email)) {
			key = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}

	return result, nil
}
