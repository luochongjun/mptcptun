package jsonv4

import (
	"fmt"

	handlerService "github.com/v2fly/v2ray-core/v5/app/proxyman/command"
	"github.com/v2fly/v2ray-core/v5/main/commands/all/api"
	"github.com/v2fly/v2ray-core/v5/main/commands/base"
)

var cmdListUsers = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api users [--server=127.0.0.1:8080] -tags <tag1> ...",
	Short:       "list bound users",
	Long: `
Remove inbounds from V2Ray.

> Make sure you have "HandlerService" set in "config.api.services" 
of server config.

Arguments:

	-tags
		The input are tags instead of config files

	-s, -server <server:port>
		The API server address. Default 127.0.0.1:8080

	-t, -timeout <seconds>
		Timeout seconds to call API. Default 3

Example:

    {{.Exec}} {{.LongName}} -tags tag1 tag2
`,
	Run: executeListUsers,
}

func executeListUsers(cmd *base.Command, args []string) {
	api.SetSharedFlags(cmd)
	api.SetSharedConfigFlags(cmd)
	isTags := cmd.Flag.Bool("tags", false, "")
	cmd.Flag.Parse(args)

	var tags []string
	if *isTags {
		tags = cmd.Flag.Args()
	}
	if len(tags) == 0 {
		base.Fatalf("no inbound to special")
	}

	conn, ctx, close := api.DialAPIServer()
	defer close()

	client := handlerService.NewHandlerServiceClient(conn)
	for _, tag := range tags {
		r := &handlerService.ListUsersRequest{
			Tag: tag,
		}
		ret, err := client.ListUsers(ctx, r)
		if err != nil {
			base.Fatalf("failed to list users: %s", err)
		}
		for k,v := range ret.Users {
			fmt.Println(tag," ",k," ",v)
		}
	}
}
