package main

import (
	"os"
	"os/signal"
	"time"
	"net/smtp"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"github.com/shirou/gopsutil/cpu"
)
const (
	CONFIG = "config.json"
)


type Config struct {
	Duration int    `json:"duration"`
	TimeUnit string `json:"timeUnit"`
	PercentAlarm int `json:"percentAlarm"`
	SmtpHost string `json:"smtpHost"`
	SmtpPort int `json:"smtpPort"`
	From    string `json:"from"`
	Pass    string `json:"pass"`
	To     string `json:"to"`
	Message string `json:"message"`
}

var conf Config

func UpdateConfig() {
	file, err := ioutil.ReadFile(CONFIG)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &conf)
}

func initConfig() {
	if (Config{}) != conf {
		return
	}
	UpdateConfig()
}
func getDuration() int {
	initConfig()
	return conf.Duration
}
func getTimeUnit() string {
	initConfig()
	return conf.TimeUnit
}
func getSmtpHost() string {
	initConfig()
	return conf.SmtpHost
}
func getSmtpPort() int {
	initConfig()
	return conf.SmtpPort
}
func getFrom() string {
	initConfig()
	return conf.From
}
func getPass() string {
	initConfig()
	return conf.Pass
}
func getTo() string {
	initConfig()
	return conf.To
}
func getMessage() string {
	initConfig()
	return conf.Message
}
func getPercentAlarm() int {
	initConfig()
	return conf.PercentAlarm
}
func main() {

	sigintCh := make(chan os.Signal)
	signal.Notify(sigintCh, os.Interrupt)

	go cpuCheck()
	<-sigintCh
}

func cpuCheck()  {

	for {

		var dur time.Duration

		if getTimeUnit() == `second` {
			dur = time.Duration(getDuration()) * time.Second
		} else if getTimeUnit() == `minute`{
			dur = time.Duration(getDuration()) * time.Minute
		}

		float, err := cpu.Percent(dur, false)

		if err != nil {
			log.Fatal(err)
		}

		for _, usage := range float {
			println(int(usage))
			if int(usage) > getPercentAlarm() {
				send(getMessage())
				}
			}
	}
}

func send(body string)  {
	from := getFrom()
	pass := getPass()
	to := getTo()

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: CPU usage overload\n\n" +
		body

	err := smtp.SendMail(getSmtpHost() +":"+ strconv.Itoa(getSmtpPort()),
		smtp.PlainAuth("", from, pass, getSmtpHost()),
		from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}