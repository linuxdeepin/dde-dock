package calendar

import "github.com/jinzhu/gorm"

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

func (j *JobType) validate() error {
	// TODO
	return nil
}
