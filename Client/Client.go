
package main

import (
	"context"
	"fmt"

	h "github.com/MieMilvang/DISYSMockExam/HelperMethod"
	"github.com/MieMilvang/DISYSMockExam/Proto"
	"google.golang.org/grpc"
)

type User struct {
	userId int64
}

func main() {
	//init
	//setup a connection, this is a blocking call
	conn, err := grpc.Dial(h.FRONT_END_ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	h.CheckError(err, "Main when to FRONT_END_ADDRESS")
	defer conn.Close()

	//create client
	ctx := context.Background()
	client := Proto.NewProtoServiceClient(conn)

	//create user
	id := ChooseUserId()
	user := User{userId: int64(id)}

	//wait for input
	user.ListenForInput(client, ctx)
	fmt.Println("main has ended")
}

// listen for user input as a string and send request accordingly, implement logic in seperate methods as see fit
func (u *User) ListenForInput(client Proto.ProtoServiceClient, ctx context.Context) {
	for {
		var input string
		fmt.Scanln(&input)
		if len(input) > 0 {
			switch input {
			case "joinchat":
				response, err := client.JoinService(ctx, &Proto.JoinRequest{UserId: u.userId})
				h.CheckError(err, "listenForInput joinchat")
				fmt.Println(response.GetMsg())
			case "getvalue":
				response, err := client.GetValue(ctx, &Proto.GetRequest{})
				h.CheckError(err, "ListenForInput getvalue")
				fmt.Println(response)
			case "setvalue":
				fmt.Println("Choose a none 0, integer you want to set the value")
				var value int64
				fmt.Scanln(&value)
				if value == 0 {
					fmt.Println("try a none 0 integer")
					break
				}
				response, err := client.SetValue(ctx, &Proto.SetRequest{UserId: u.userId, RequestedValue: value})
				h.CheckError(err, "ListenForInput setvalue")
				fmt.Println(response.GetMsg())
			default:
				fmt.Printf("unrecognised command: %v\n", input)
				fmt.Println("try, 'setvalue', 'getvalue', or 'joinchat'")
			}
		}
	}
}

// helper method wait for user to input a wanted userid
func ChooseUserId() int {
	for {
		var userId int
		fmt.Println("Choose a positive integer as id:")
		fmt.Scanln(&userId)
		if userId > 0 {
			return userId
		} else {
			fmt.Printf("\nSomething went wrong when processing your input: %v\nTry a positive integer like 1, 2, 3, ...", userId)
		}
	}
}
