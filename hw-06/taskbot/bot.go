package main

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"strconv"
	"net/http"
)

var (
	// @BotFather gives you this
	BotToken   = "473604524:AAFMxNFvI-vgUJDH-tB7SkzeJJD-4M8zeMU"
	WebhookURL = "https://ed940f60.ngrok.io"

)

type Task struct {
	text string
	by string
	byID int64
	to string
	toID int64
}

var tasks = map[int] * Task{}
var taskID = 0


func sendmMsg(userID int64, bot *tgbotapi.BotAPI, answer string){
	bot.Send(tgbotapi.NewMessage(
		userID,
		answer,
	))
}

func getTasks(userID int64, bot *tgbotapi.BotAPI, name string) {
	var msg string
	lenTasks := len(tasks)
	if lenTasks == 0 {
		sendmMsg(userID, bot, "Нет задач")
	} else {
		size := 0
		for i, task := range tasks {
			msg += strconv.Itoa(i + 1) + ". " + task.text + " by @" + task.by
			if task.to != "" {
				if task.to == name {
					msg += "\nassignee: я\n/unassign_" + strconv.Itoa(i + 1) + " /resolve_" + strconv.Itoa(i + 1)
				} else {
					msg += "\nassignee: @" + task.to;
				}
			} else {
				msg += "\n/assign_" + strconv.Itoa(i + 1)
			}

			size++
			if size != lenTasks {
				msg += "\n\n"
			}
		}
		sendmMsg(userID, bot, msg)
	}
}

func addTasks(anyTasks map[int]*Task, user string){
	for i, task := range tasks {
		if (task.by == user) {
			anyTasks[i] = task
		}
	}	
}

func getMyTasks(userID int64, bot *tgbotapi.BotAPI, user string) {
	myTasks := map[int] * Task{}
	var msg string
	lenTasks := len(tasks)

	if lenTasks == 0 {
		sendmMsg(userID, bot, "Нет задач")
	} else {
		for i, task := range tasks {
			if (task.to == user) {
				myTasks[i] = task
			}
		}	
		if len(myTasks) == 0 {
			sendmMsg(userID, bot, "Нет задач")
		} else {
			size := 0;
			for i, task := range myTasks {
				msg += strconv.Itoa(i + 1) + ". " + task.text + " by @" + task.by
				msg += "\n/unassign_" + strconv.Itoa(i + 1) + " /resolve_" + strconv.Itoa(i + 1)
				size++
			}
			sendmMsg(userID, bot, msg)
		}
	}
}

func getOwnTasks(userID int64, bot *tgbotapi.BotAPI, user string) {
	var ownTasks = map[int] * Task{}
	var msg string
	lenTasks := len(tasks)
	

	if lenTasks == 0 {
		sendmMsg(userID, bot, "Нет задач")
	} else {
		for i, task := range tasks {
			if (task.by == user) {
				ownTasks[i] = task
			}
		}	
		if len(ownTasks) == 0 {
			sendmMsg(userID, bot, "Нет задач")
		} else {
			size := 0
			for i, task := range ownTasks {
				msg += strconv.Itoa(i + 1) + ". " + task.text + " by @" + task.by
				if (task.to != "") {
					if (task.to == user) {
						msg += "\nassignee: я\n/unassign_" + strconv.Itoa(i+ 1) + " /resolve_" + strconv.Itoa(i + 1)
					} else {
						msg += "\nassignee: @" + task.to;
					}
				} else {
					msg += "\n/assign_" + strconv.Itoa(i + 1)
				}
				size++
			}
			sendmMsg(userID, bot, msg)
		}
	}
}

func getNewTask(text string, userID int64, bot *tgbotapi.BotAPI, user string) {
	task := new(Task)
	
	task.text = text
	task.by = user
	task.to = ""
	task.byID = userID
	task.toID = 0

	taskID++
	tasks[taskID - 1] = task
	msg := "Задача \"" + text + "\" создана, id=" + strconv.Itoa(taskID)
	sendmMsg(userID, bot, msg)
}

func getAssignTask(assignID int, userID int64, bot *tgbotapi.BotAPI, user string) {
	var msg string
	if tasks[assignID - 1].to == "" {
		tasks[assignID - 1].to = user
		tasks[assignID - 1].toID = userID
		if userID != tasks[assignID - 1].byID {
			msg = "Задача \"" + tasks[assignID - 1].text + "\" назначена на @" + user
			sendmMsg(tasks[assignID - 1].byID, bot, msg)
		}
		msg = "Задача \"" + tasks[assignID - 1].text + "\" назначена на вас"
		sendmMsg(userID, bot, msg)
	} else {
		msg = "Задача \"" + tasks[assignID - 1].text + "\" назначена на @" + user
		sendmMsg(tasks[assignID - 1].toID, bot, msg)
		tasks[assignID - 1].to = user
		tasks[assignID - 1].toID = userID
		msg = "Задача \"" + tasks[assignID - 1].text + "\" назначена на вас"
		sendmMsg(userID, bot, msg)
	}

}

func getUnassignTask(unassignID int, userID int64, bot *tgbotapi.BotAPI, user string) {
	var msg string
	
	if (tasks[unassignID - 1].to != user) {
		msg = "Задача не на вас"
	} else {
		tasks[unassignID - 1].to = ""
		if userID != tasks[unassignID - 1].byID {
			msg = "Задача \"" + tasks[unassignID - 1].text + "\" осталась без исполнителя"
			sendmMsg(tasks[unassignID - 1].byID, bot, msg)
		}
		msg = "Принято"
	}
	sendmMsg(userID, bot, msg)
}

func getResolveTask(resolveID int, userID int64, bot *tgbotapi.BotAPI, user string) {
	
	var msg string

	if (tasks[resolveID - 1].to != user) {
		msg = "Задача не на вас"
	} else {
		if userID != tasks[resolveID - 1].byID {
			msg = "Задача \"" + tasks[resolveID - 1].text + "\" выполнена @" + tasks[resolveID - 1].to
			sendmMsg(tasks[resolveID - 1].byID, bot, msg)
		}
		msg = "Задача \"" + tasks[resolveID - 1].text + "\" выполнена"
		delete(tasks, resolveID - 1)
	}
	sendmMsg(userID, bot, msg)
}

func startTaskBot(ctx context.Context) error {

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8081", nil)
	fmt.Println("start listen :8081")

	// var tasks map[int]string
	
	
	for update := range updates {
		userID := update.Message.Chat.ID;
		textID := update.Message.Text;
		usernameID :=  update.Message.Chat.UserName

		if update.Message == nil {
            continue
		}
		if strings.HasPrefix(textID, "/new") {
			text := textID[5:]
			getNewTask(text, userID, bot, usernameID)

		}	else if strings.HasPrefix(textID, "/assign") {
			assignID, _ := strconv.Atoi(textID[len("/assign") + 1:len("/assign") + 2])
			getAssignTask(assignID, userID, bot, usernameID)

		} else if strings.HasPrefix(textID, "/unassign") {
			unassignID, _ := strconv.Atoi(textID[len("/unassign") + 1:len("/unassign") + 2])
			getUnassignTask(unassignID, userID, bot, usernameID)

		} else if strings.HasPrefix(textID, "/resolve") {
			resolveID, _ := strconv.Atoi(textID[len("/resolve") + 1:len("/resolve") + 2])
			getResolveTask(resolveID, userID, bot, usernameID)
		} 

		switch textID {
			case "/start":
				sendmMsg(userID, bot, "Welcome to taskmanager bot!")
			case "/tasks":
				getTasks(userID, bot, usernameID)
			case "/my":
				getMyTasks(userID, bot, usernameID)
			case "/owner":
				getOwnTasks(userID, bot, usernameID)
			default:
				//sendmMsg(userID, bot, "Choose right comand!")		
		}
	}	

	return nil
}


func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
