package main

import (
	"auth_user/app/config"
	"auth_user/app/db"
	"auth_user/app/pb"
	"auth_user/app/services"
	"auth_user/app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	//"gorm.io/driver/postgres"
	//"gorm.io/driver/sqlite"

	//"gorm.io/gorm"
	"log"
	"net"
)

//func createPgDb() {
//	cmd := exec.Command("createdb", "-p", "5432", "-h", "127.0.0.1", "-U", "superuser", "-e", "test_db")
//
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	if err := cmd.Run(); err != nil {
//		log.Printf("Error: %v", err)
//	}
//	log.Printf("Output: %q\n", out.String())
//}

func main() {
	//createPgDb()
	//db1, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	//fmt.Println("sdsdfsd")
	//if err != nil {
	//	panic("failed to connect database")
	//}
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", fmt.Sprintf("localhost:%d", 50051))

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}
	rg := gin.Default()
	s.StartHttp(rg.Group("user/"))
	grpcServer := grpc.NewServer()
	go rg.Run(":8082")
	pb.RegisterAuthServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
