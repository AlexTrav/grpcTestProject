package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "grpcTestProject/userpb"

	"google.golang.org/grpc"
)

// gRPC клиент с интерактивным меню в терминале

// Поддерживает: Create, Get, Update, Delete, List пользователей

func main() {
	// Устанавливаем соединение с gRPC-сервером по адресу localhost:50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close() // Гарантированное закрытие соединения при завершении работы

	// Создаём gRPC-клиент по сгенерированному интерфейсу
	client := pb.NewUserServiceClient(conn)

	// Создаём читателя для ввода с консоли
	reader := bufio.NewReader(os.Stdin)

	// Бесконечный цикл меню
	for {
		// Отображение интерактивного меню
		fmt.Println("=== Меню ===")
		fmt.Println("1. Создать пользователя")
		fmt.Println("2. Получить пользователя по ID")
		fmt.Println("3. Обновить пользователя")
		fmt.Println("4. Удалить пользователя")
		fmt.Println("5. Показать всех пользователей")
		fmt.Println("0. Выход")
		fmt.Print("Выберите действие: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// Считывание данных для создания пользователя
			fmt.Print("Введите имя: ")
			name := readLine(reader)
			fmt.Print("Введите email: ")
			email := readLine(reader)
			fmt.Print("Введите возраст: ")
			var age int32
			fmt.Scan(&age)

			// Устанавливаем контекст с таймаутом 3 секунды
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// gRPC-вызов метода CreateUser
			res, err := client.CreateUser(ctx, &pb.UserRequest{
				Name:  name,
				Email: email,
				Age:   age,
			})
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("✅ Пользователь создан:", res.User)
			}
			reader.ReadString('\n') // Очищаем буфер после fmt.Scan

		case "2":
			fmt.Print("Введите ID пользователя: ")
			id := readLine(reader)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// gRPC-вызов метода GetUser
			res, err := client.GetUser(ctx, &pb.UserID{Id: id})
			if err != nil {
				fmt.Println("❌ Ошибка:", err)
			} else {
				fmt.Println("👤 Пользователь:", res.User)
			}

		case "3":
			fmt.Print("Введите ID пользователя: ")
			id := readLine(reader)
			fmt.Print("Новое имя: ")
			name := readLine(reader)
			fmt.Print("Новый email: ")
			email := readLine(reader)
			fmt.Print("Новый возраст: ")
			var age int32
			fmt.Scan(&age)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// gRPC-вызов метода UpdateUser
			res, err := client.UpdateUser(ctx, &pb.UpdateRequest{
				Id:    id,
				Name:  name,
				Email: email,
				Age:   age,
			})
			if err != nil {
				fmt.Println("❌ Ошибка:", err)
			} else {
				fmt.Println("✏️ Пользователь обновлён:", res.User)
			}
			reader.ReadString('\n')

		case "4":
			fmt.Print("Введите ID пользователя для удаления: ")
			id := readLine(reader)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// gRPC-вызов метода DeleteUser
			res, err := client.DeleteUser(ctx, &pb.UserID{Id: id})
			if err != nil {
				fmt.Println("❌ Ошибка:", err)
			} else {
				fmt.Println("🗑️", res.Message)
			}

		case "5":
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// gRPC-вызов метода ListUsers
			res, err := client.ListUsers(ctx, &pb.Empty{})
			if err != nil {
				fmt.Println("❌ Ошибка:", err)
			} else if len(res.Users) == 0 {
				fmt.Println("👥 Пользователи не найдены.")
			} else {
				fmt.Println("📋 Список пользователей:")
				for _, user := range res.Users {
					fmt.Printf("- ID: %s | Имя: %s | Email: %s | Возраст: %d\n", user.Id, user.Name, user.Email, user.Age)
				}
			}

		case "0":
			fmt.Println("👋 Выход из программы.")
			return

		default:
			fmt.Println("❓ Неверный выбор")
		}
	}
}

// readLine — утилита для считывания строки из stdin
func readLine(r *bufio.Reader) string {
	text, _ := r.ReadString('\n')
	return strings.TrimSpace(text)
}
