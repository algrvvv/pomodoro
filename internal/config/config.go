package config

import "time"

type Integration struct {
	Name            string        `yaml:"name"    json:"name"`
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	ApiKey          string        `yaml:"api_key" json:"-"`
	TimeoutInSecond int           `yaml:"timeout" json:"-"`
	Timeout         time.Duration `yaml:"-"       json:"-"`
}

type PomodoroConfig struct {
	// StartWorkReminder колво минут через которое будет высвечено напоминание о начале
	// новой рабочей сессии
	StartSessionReminder int `yaml:"startSessionReminder"`

	// WorkMinutes колво минут выделенное на работу
	WorkMinutes int `yaml:"workMinutes"`

	// ShortBreakMinutes колво минут выделенное на короткий отдых между сессиями
	ShortBreakMinutes int `yaml:"shortBreakMinutes"`

	// LongBreakMinutes колво минут выделенное на длинный отдых, который начинается
	// после достижения нужного колва пройденных сессий
	LongBreakMinutes int `yaml:"longBreakMinutes"`

	// BreakAfterSessions колво сессий после которых будет запущен длинный отдых
	BreakAfterSessions int `yaml:"breakAfterSessions"`

	// SessionsGoal цель сессий на день
	SessionsGoal int `yaml:"sessionsGoal"`

	// SessionGoalMinutes цель сессий на день в минутах, если какая то из сессий, например, не была полностью закончена
	SessionGoalMinutes int
}

type AppConfig struct {
	Port            string `yaml:"port"`
	DisplayZeroDays bool   `yaml:"displayZeroDays"`
}

type DatabaseConfig struct {
	Addr   string `yaml:"addr"`
	User   string `yaml:"user"`
	Passwd string `yaml:"pass"`
	DBName string `yaml:"dbname"`
}

type Cfg struct {
	App          AppConfig      `yaml:"app"`
	DB           DatabaseConfig `yaml:"db"`
	Pomodoro     PomodoroConfig `yaml:"pomodoro"`
	Intergations []Integration  `yaml:"integrations"`
}

var Config Cfg
