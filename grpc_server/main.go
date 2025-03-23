package main

import (
	"context"
	"log"
	"net"

	pb "grpcTestProject/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Структура сервера
type server struct {
	pb.UnimplementedUserServiceServer
}

// Реализация методов gRPC-сервиса

// CreateUser — создание нового пользователя с UUID
func (s *server) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	id := uuid.New().String()
	user := User{ID: id, Name: req.Name, Email: req.Email, Age: req.Age}
	// Логируем попытку
	log.Println("➡️ Попытка создать пользователя:", user)

	if err := DB.Create(&user).Error; err != nil {
		log.Printf("❌ Ошибка при создании пользователя: %v\n", err)
		return nil, err
	}

	// Подтверждение в логах
	log.Printf("✅ Пользователь создан: ID=%s, Name=%s\n", user.ID, user.Name)

	return &pb.UserResponse{User: &pb.User{
		Id: user.ID, Name: user.Name, Email: user.Email, Age: user.Age,
	}}, nil
}

// GetUser — получить пользователя по ID
func (s *server) GetUser(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	var user User
	if err := DB.First(&user, "id = ?", req.Id).Error; err != nil {
		log.Printf("❌ Пользователь с ID=%s не найден\n", req.Id)
		return nil, err
	}

	log.Printf("🔍 Получен пользователь: ID=%s\n", user.ID)

	return &pb.UserResponse{User: &pb.User{
		Id: user.ID, Name: user.Name, Email: user.Email, Age: user.Age,
	}}, nil
}

// UpdateUser — обновление пользователя по ID
func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {
	log.Printf("🔄 Обновление пользователя с ID=%s", req.Id)

	var user User
	if err := DB.First(&user, "id = ?", req.Id).Error; err != nil {
		log.Printf("❌ Не удалось найти пользователя: ID=%s\n", req.Id)
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	if err := DB.Save(&user).Error; err != nil {
		log.Printf("❌ Ошибка при обновлении: %v", err)
		return nil, err
	}

	log.Printf("✅ Пользователь обновлён: %+v\n", user)

	return &pb.UserResponse{User: &pb.User{
		Id: user.ID, Name: user.Name, Email: user.Email, Age: user.Age,
	}}, nil
}

// DeleteUser — удаление пользователя по ID
func (s *server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.DeleteResponse, error) {
	if err := DB.Delete(&User{}, "id = ?", req.Id).Error; err != nil {
		log.Printf("❌ Ошибка при удалении: %v\n", err)
		return nil, err
	}
	log.Printf("🗑️ Пользователь удалён: ID=%s", req.Id)
	return &pb.DeleteResponse{Message: "Пользователь удалён"}, nil
}

// ListUsers — вывод всех пользователей
func (s *server) ListUsers(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		log.Println("❌ Ошибка при получении пользователей:", err)
		return nil, err
	}

	log.Printf("📋 Получено пользователей: %d", len(users))

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id: u.ID, Name: u.Name, Email: u.Email, Age: u.Age,
		})
	}
	return &pb.UserList{Users: pbUsers}, nil
}

func main() {
	// Инициализация подключения к базе данных и миграция модели
	InitDB()

	// Создание TCP-соединения, на котором будет слушать gRPC-сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}

	// Создание нового gRPC-сервера
	s := grpc.NewServer()

	// Регистрация реализации сервиса UserService на сервере
	pb.RegisterUserServiceServer(s, &server{})

	log.Println("🚀 gRPC сервер запущен на порту 50051")

	// Запуск gRPC-сервера на указанном порту (в данном случае 50051)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ Ошибка сервера: %v", err)
	}
}
