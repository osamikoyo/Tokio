package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "tokio/tokio" // Замените на правильный путь к вашему сгенерированному пакету
)
type Client struct {
	name   string
	stream pb.ChatService_JoinChatServer
}
type server struct {
	pb.UnimplementedChatServiceServer
	clients map[string]*Client
	mu      sync.Mutex
}
func (s *server) GiveClientName(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error){
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Отправка сообщения только клиентам с соответствующим именем
	for _, client := range s.clients {
		if client.name == req.ChatMessage.UserTo {
			err := client.stream.Send(req.ChatMessage)
			if err != nil {
				log.Printf("Failed to send message to client %s: %v", client.name, err)
			}
		}
	}
	return &pb.SendMessageResponse{Status: "Message sent"}, nil
}
func MetadataFromIncomingContext(ctx context.Context) (metadata.MD, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	return md, ok
}

func (s *server) JoinChat(empty *emptypb.Empty, stream pb.ChatService_JoinChatServer) error {
	// Получение имени клиента (например, из метаданных или передача от клиента)
	md, ok := MetadataFromIncomingContext(stream.Context())
	var clientName string

	if ok {
		clientName = md["name"][0] // Предполагаем, что имя передается в метаданных
	} else {
		clientName = "Unknown" //Если имя не передано
	}

	s.mu.Lock()
	client := &Client{name: clientName, stream: stream}
	s.clients[clientName] = client
	s.mu.Unlock()

	fmt.Printf("New client joined: %s\n", clientName)

	// Обработка клиента до его отключения
	<-stream.Context().Done()

	// Удаляем клиента при отключении
	s.mu.Lock()
	delete(s.clients, clientName)
	s.mu.Unlock()

	fmt.Printf("Client left: %s\n", clientName)
	return nil
}


func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterChatServiceServer(s, &server{clients: make(map[string]*Client)})

	fmt.Println("Chat server is running on port: 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}