package chatbots

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/bot/irc"
	_ "github.com/go-chat-bot/plugins/catfacts"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
	"github.com/jtaylorcpp/quaditor"

	log "github.com/sirupsen/logrus"
	// Import all the commands you wish to use
)

var (
	subjectRegex   = regexp.MustCompile(`s:(\w+)`)
	objectRegex    = regexp.MustCompile(`o:(\w+)`)
	predicateRegex = regexp.MustCompile(`p:(\w+)`)
)

var auditor quaditor.Auditor = nil

func init() {
	bot.RegisterCommand(
		"query",
		"Sends a query to Quaditor.",
		"",
		query,
	)

	bot.RegisterCommandV3(
		"queryv3",
		"Sends a query to Quaditor.",
		"",
		queryv3,
	)
}

type IRCBot struct {
	config *irc.Config
}

func NewIRCBot(auditor quaditor.Auditor) *IRCBot {
	auditor = auditor
	return &IRCBot{
		config: &irc.Config{
			Server:   "192.168.86.95:6667",
			Channels: []string{"#quaditor"},
			Nick:     "quadbot",
			UseTLS:   false,
			Debug:    true,
			User:     "quadbot",
			Password: "",
		},
	}
}

func (bot *IRCBot) Run() {
	irc.Run(bot.config)
}

func query(command *bot.Cmd) (msg string, err error) {
	log.Infof("cmd: %#v\n", *command)
	queriesLines := strings.Split(command.RawArgs, "|")
	queryStatements := [][]string{}
	for _, queryLine := range queriesLines {
		queryStatements = append(queryStatements, strings.Split(queryLine, ";"))
	}
	queries := []quaditor.Query{}
	for _, queryStatement := range queryStatements {
		query := quaditor.Query{}
		for idx, quadString := range queryStatement {
			subjects := subjectRegex.FindStringSubmatch(quadString)
			objects := objectRegex.FindStringSubmatch(quadString)
			predicates := predicateRegex.FindStringSubmatch(quadString)

			quad := quaditor.Quad{}
			if len(subjects) == 2 {
				quad.Subject = subjects[1]
			}
			if len(objects) == 2 {
				quad.Object = objects[1]
			}
			if len(predicates) == 2 {
				quad.Predicate = predicates[1]
			}
			switch idx {
			case 0:
				query.Constraint = quad
			case 1:
				query.Assignment = quad
			default:
				panic("should not be more than 2 quads per query")
			}
		}
		queries = append(queries, query)
	}
	msg = fmt.Sprintf("Hello %#v", queries)
	return
}

func queryv3(command *bot.Cmd) (bot.CmdResultV3, error) {
	// setup return chan
	var channel string
	if command.ChannelData.IsPrivate {
		channel = command.User.Nick
	} else {
		channel = command.Channel
	}
	returnCmd := bot.CmdResultV3{
		Channel: channel,
		Message: make(chan string),
		Done:    make(chan bool, 1),
	}

	// do query
	go func(command *bot.Cmd, returnCommand bot.CmdResultV3) {
		log.Infof("cmd: %#v\n", *command)
		returnCmd.Message <- fmt.Sprintf("processing query: %#v\n", command.RawArgs)
		queriesLines := strings.Split(command.RawArgs, "|")
		queryStatements := [][]string{}
		for _, queryLine := range queriesLines {
			queryStatements = append(queryStatements, strings.Split(queryLine, ";"))
		}
		queries := []quaditor.Query{}
		for _, queryStatement := range queryStatements {
			query := quaditor.Query{}
			for idx, quadString := range queryStatement {
				subjects := subjectRegex.FindStringSubmatch(quadString)
				objects := objectRegex.FindStringSubmatch(quadString)
				predicates := predicateRegex.FindStringSubmatch(quadString)

				quad := quaditor.Quad{}
				if len(subjects) == 2 {
					quad.Subject = subjects[1]
				}
				if len(objects) == 2 {
					quad.Object = objects[1]
				}
				if len(predicates) == 2 {
					quad.Predicate = predicates[1]
				}
				switch idx {
				case 0:
					query.Constraint = quad
				case 1:
					query.Assignment = quad
				default:
					returnCommand.Message <- ("should not be more than 2 quads per query")
					returnCommand.Done <- true
					return
				}
			}
			returnCommand.Message <- fmt.Sprintf("processing query: %#v\n", query)
			queries = append(queries, query)
		}
		returnCommand.Message <- "all queries processed"
		returnCommand.Message <- "results:\n"
		paths, err := auditor.Query(queries...)
		if err != nil {
			returnCommand.Message <- "ran into error: " + err.Error()
		} else {
			for idx, path := range paths {
				returnCommand.Message <- fmt.Sprintf("%v: %s\n", idx, path.String())
			}
		}
		returnCommand.Done <- true
	}(command, returnCmd)

	return returnCmd, nil
}
