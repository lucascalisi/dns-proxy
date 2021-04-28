package blocker

import (
	"bufio"
	"dns-proxy/pkg/domain/proxy"
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

func (u *updater) UpdateAll() map[string]bool {
	if u.nextRefresh.Before(time.Now()) {

		var blocklist map[string]bool
		blocklist = make(map[string]bool)
		for _, s := range u.sources {
			resp, _ := http.Get(s)
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				expresion, err := u.parseRegister(scanner.Text())
				if err == nil && expresion != nil {
					blocklist[*expresion] = true
				}
			}
			defer resp.Body.Close()
		}

		u.nextRefresh = time.Now().Add(u.refresh)
		return blocklist
	}
	return nil
}

func (u *updater) parseRegister(register string) (*string, error) {
	if register != "" && register[:1] != "#" {
		splited := strings.Split(register, " ")
		if len(splited) == 1 {
			return &register, nil
		}
	}

	return nil, nil
}
