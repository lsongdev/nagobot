package cron

type ScheduleKind string

const (
	ScheduleAt    ScheduleKind = "at"
	ScheduleEvery ScheduleKind = "every"
	ScheduleCron  ScheduleKind = "cron"
)

type Schedule struct {
	Kind    ScheduleKind `yaml:"kind"`
	AtMs    int64        `yaml:"at_ms,omitempty"`
	EveryMs int64        `yaml:"every_ms,omitempty"`
	Expr    string       `yaml:"expr,omitempty"`
	Tz      string       `yaml:"tz,omitempty"`
}

type Payload struct {
	Message string `yaml:"message"`
	Deliver bool   `yaml:"deliver,omitempty"`
	Channel string `yaml:"channel,omitempty"`
	To      string `yaml:"to,omitempty"`
}

type Job struct {
	ID             string   `yaml:"id"`
	Name           string   `yaml:"name"`
	Enabled        bool     `yaml:"enabled"`
	Schedule       Schedule `yaml:"schedule"`
	Payload        Payload  `yaml:"payload"`
	DeleteAfterRun bool     `yaml:"delete_after_run,omitempty"`
}

type Config struct {
	Jobs []Job `yaml:"jobs"`
}
