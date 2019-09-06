package calendar

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

type JobType struct {
	gorm.Model

	Name  string
	Color string
}

type JobTypeJSON struct {
	ID    uint
	Name  string
	Color string
}

func (j *JobType) toJobTypeJSON() *JobTypeJSON {
	if j == nil {
		return nil
	}
	return &JobTypeJSON{
		ID:    j.ID,
		Name:  j.Name,
		Color: j.Color,
	}
}

func (j *JobTypeJSON) toJobType() *JobType {
	if j == nil {
		return nil
	}
	jt := &JobType{
		Name:  j.Name,
		Color: j.Color,
	}
	jt.ID = j.ID
	return jt
}

var colorReg = regexp.MustCompile(`^#[0-9a-f]+$`)

func (j *JobType) validate() error {
	if strings.TrimSpace(j.Name) == "" {
		return errors.New("name is empty")
	}

	color := strings.ToLower(j.Color)
	if !colorReg.MatchString(color) {
		return errors.New("invalid color")
	}
	switch len(color) - 1 {
	case 3, 4, 6, 8:
		// rgb 6 缩写 3
		// rgba 8 缩写 4
		//pass
	default:
		return errors.New("invalid color length")
	}

	return nil
}
