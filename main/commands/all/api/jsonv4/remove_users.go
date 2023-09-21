package jsonv4

import (
	"encoding/json"

	"github.com/v2fly/v2ray-core/v5/common/protocol"

	handlerService "github.com/v2fly/v2ray-core/v5/app/proxyman/command"
	"github.com/v2fly/v2ray-core/v5/main/commands/all/api"
	"github.com/v2fly/v2ray-core/v5/main/commands/base"
)

var cmdRmUsers = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api rmu [--server=127.0.0.1:8080] -users <users json obj> ...",
	Short:       "remove bound users",
	Long: `
Remove inbounds from V2Ray.

> Make sure you have "HandlerService" set in "config.api.services" 
of server config.

Arguments:

	-users
		The input are json array

	-s, -server <server:port>
		The API server address. Default 127.0.0.1:8080

	-t, -timeout <seconds>
		Timeout seconds to call API. Default 3

Example:

    {{.Exec}} {{.LongName}} -users {tag":"mytag","users":["user1","user2"]}
`,
	Run: executeRmUsers,
}

type  rmUserJsonArry struct {
	Tag string
	Users []string
}

func makeUser(user string) *protocol.User {
	u := new(protocol.User)
	u.Email = user
	return u
}

func executeRmUsers(cmd *base.Command, args []string) {
	api.SetSharedFlags(cmd)
	api.SetSharedConfigFlags(cmd)
	isUsers := cmd.Flag.Bool("users", false, "")
	cmd.Flag.Parse(args)

	var users []string
	if *isUsers {
		users = cmd.Flag.Args()
	}
	if len(users) == 0 {
		base.Fatalf("no user info to special")
	}

	conn, ctx, close := api.DialAPIServer()
	defer close()

	client := handlerService.NewHandlerServiceClient(conn)

	for _, user := range users {
		userArry := new(rmUserJsonArry)
		if err := json.Unmarshal([]byte(user),userArry); err == nil {
			var reqUsers []*protocol.User
			for _,v := range userArry.Users {
				u := makeUser(v)
				reqUsers = append(reqUsers,u)
			}
			if len(reqUsers) > 0 {
				r := &handlerService.RmUsersRequest{
					Tag: userArry.Tag,
					Users: reqUsers,
				}
				_, err := client.RmUsers(ctx, r)
				if err != nil {
					base.Fatalf("failed to add user: %s", err)
				}
			}
		}
	}
}
