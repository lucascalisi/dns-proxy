package blocker

import (
	"bufio"
	"dns-proxy/pkg/domain/proxy"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type updater struct {
	refresh     time.Duration
	nextRefresh time.Time
	sources     []string
}

func NewUpdater(refresh time.Duration, sources []string) proxy.ListUpdater {
	return &updater{refresh, time.Time{}, sources}
}

func (u *updater) Update(source string) error {
	return nil
}

func (u *updater) UpdateAll() (map[string]bool, int) {
	if u.nextRefresh.Before(time.Now()) {
		errors := 0
		var blocklist map[string]bool
		blocklist = make(map[string]bool)
		for _, s := range u.sources {
			resp, _ := http.Get(s)
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				expresion, err := u.parseRegister(scanner.Text())
				if err != nil {
					errors++
				} else if expresion != nil {
					blocklist[*expresion] = true
				}
			}
			defer resp.Body.Close()
		}

		u.nextRefresh = time.Now().Add(u.refresh)
		return blocklist, errors
	}
	return nil, 0
}

func (u *updater) parseRegister(register string) (*string, error) {
	if register != "" && register[:1] != "#" {
		i := strings.Index(register, "#")
		if i > 0 {
			register = register[:i]
		}
		splited := strings.Split(register, " ")
		lenSplited := len(splited)
		if lenSplited == 1 {
			return &register, nil
		} else if lenSplited == 2 {
			return &splited[1], nil
		}
		return nil, fmt.Errorf("could not parse register %s", register)
	}

	return nil, nil
}
