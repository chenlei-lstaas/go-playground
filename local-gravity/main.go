package main

import (
	context "context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/MrVegeta/go-playground/local-gravity/fast"

	"google.golang.org/grpc"
)

var (
	voyager = flag.String("voyager", "127.0.0.1:8079", "Voyager address.")
	execute = flag.String("exec", "query", "Command type, add/delete/query.")
	rule    = flag.String("rule", "{\"id\":\"1\",\"protocol\":1,\"type\":1,\"address\":\"www.uzing.io\",\"port\":443}", "Rule(s) to add.")
	enable  = flag.String("enable", "", "Enable all Custom/Gravity rules.")
	disable = flag.String("disable", "", "Disable all Custom/Gravity rules.")
)

// Rules ...
type Rule struct {
	ID       string `yaml:"id"`
	Protocol int    `yaml:"protocol"`
	Type     int    `yaml:"type"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
}

func ruleOperation(exec, addr string, rule *Rule) {
	log.Println("connecting voyager:", addr)

	conn, err := grpc.Dial(*voyager, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	client := fast.NewCustomRuleManagerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	switch exec {
	case "add":
		request := &fast.RuleRequest{
			Rules: []*fast.Rule{
				&fast.Rule{
					ID:       rule.ID,
					Protocol: int32(rule.Protocol),
					Type:     int32(rule.Type),
					Address:  rule.Address,
					Port:     int32(rule.Port),
				},
			},
		}
		if response, err := client.AddRules(ctx, request); err != nil {
			log.Panic("add rule failed, error:", err)
		} else {
			log.Println("add rule result:", response)
		}
	case "delete":
		request := &fast.RuleRequest{
			Rules: []*fast.Rule{
				&fast.Rule{
					ID:       rule.ID,
					Protocol: int32(rule.Protocol),
					Type:     int32(rule.Type),
					Address:  rule.Address,
					Port:     int32(rule.Port),
				},
			},
		}
		if response, err := client.DeleteRules(ctx, request); err != nil {
			log.Panic("delete rule failed, error:", err)
		} else {
			log.Println("delete rule result:", response)
		}
	case "query":
		request := &fast.QueryRequest{
			Type: int32(rule.Type),
		}
		if response, err := client.QueryRules(ctx, request); err != nil {
			log.Panic("query rule failed, error:", err)
		} else {
			log.Println("query rule result:", response)
		}

	}
}

func rulesSwitch(ruleType string, enable int32) { //ruleType int32, enable int32) {
	var request *fast.SwitchRequest
	ruleSwitch := fast.Switch{
		Enable: enable,
	}

	switch ruleType {
	case "custom":
		ruleSwitch.Type = 1
	case "gravity":
		ruleSwitch.Type = 0
	default:
		log.Panic("unsupported rule type:", ruleType)
	}

	log.Println("connecting voyager:", *voyager)
	conn, err := grpc.Dial(*voyager, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	client := fast.NewCustomRuleManagerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	response, err := client.RuleSwitch(ctx, request)
	if err != nil {
		log.Panic(err)
	}

	log.Println(response)
}

func main() {
	flag.Parse()

	rules := &Rule{}
	if err := json.Unmarshal([]byte(*rule), rules); err != nil {
		log.Panic(err)
	}

	if *execute != "" {
		ruleOperation(*execute, *voyager, rules)
	}
	if *enable != "" {
		rulesSwitch(*enable, 1)
	}
	if *disable != "" {
		rulesSwitch(*disable, 0)
	}
}
