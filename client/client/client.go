package client

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"tokio/client/visual"
	pb "tokio/tokio" // Замените на правильный путь к вашему сгенерированному пакету
)

func DoCorrectString(message string)string{
	result := strings.Replace(message, " ", "_", -1)
	return result
}

func RemoveCorrectString(message string) string{
	result := strings.Replace(message, "_", " ", -1)
	return result
}

func ROUTE(client pb.ChatServiceClient, name string){
	md := metadata.Pairs("name", name)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	

	_, err := client.GiveClientName(ctx, &pb.HelloRequest{
		Name: name,
	})
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error(err.Error())
	}

}

func Client() {
	// Подключение к серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	

	// Получаем имя клиента
	var clientName string
	var client_to string

	fmt.Print("Enter your name: ")
	fmt.Scanln(&clientName)
	fmt.Println("Your name is ", clientName)
	fmt.Println("Enter name who you want to sedning")
	fmt.Scanln(&client_to)

	md := metadata.Pairs("name", clientName)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	
	// Присоединение к чату
	stream, err := client.JoinChat(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error on JoinChat: %v", err)
	}
	
	// Отправляем имя клиента на сервер
	go func() {
		for {
			// Ждем сообщения от сервера
			chatMessage, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error while receiving message: %v", err)
			}
			visual.NewMessage(chatMessage.UserFrom, chatMessage.UserTo, RemoveCorrectString(chatMessage.Message))
		}
	}()

	// Основной цикл для отправки сообщений
	for {
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		// Создаем chatMessage и отправляем его
		req := &pb.SendMessageRequest{
			ChatMessage: &pb.ChatMessage{
				UserFrom: clientName,
				Message:   DoCorrectString(message),
				UserTo:   client_to, // Можно указать получателя, если необходимо
			},
		}
		// Отправляем сообщение на сервер
		_, err := client.SendMessage(context.Background(), req)
		if err != nil {
			log.Printf("Error on SendMessage: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
