package blocker

import (
	"time"
)

type Updater interface {
	Update(source string) error
	UpdateAll(source []string) error
}

type updater struct {
	refresh time.Duration
}

func NewUpdater(refresh time.Duration) Updater {
	return &updater{refresh}
}

func (u *updater) Update(source string) error {
	return nil
}

func (u *updater) UpdateAll(source []string) error {
	return nil
}
