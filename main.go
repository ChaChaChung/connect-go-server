// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"connectrpc.com/connect"
// 	"github.com/rs/cors"

// 	elizav1 "myapp/backend/gen/connectrpc/eliza/v1"                // proto messages
// 	elizav1connect "myapp/backend/gen/connectrpc/eliza/v1/v1connect" // connect service handlers
// )

// // ElizaServer 是 ElizaService 的實作
// type ElizaServer struct{}

// // Say 實作了 proto 中定義的 Say RPC
// func (s *ElizaServer) Say(
// 	ctx context.Context,
// 	req *connect.Request[elizav1.SayRequest],
// ) (*connect.Response[elizav1.SayResponse], error) {
// 	log.Println("Request headers: ", req.Header())
// 	res := connect.NewResponse(&elizav1.SayResponse{
// 		Sentence: fmt.Sprintf("Go server says: You said '%s'", req.Msg.Sentence),
// 	})
// 	res.Header().Set("Eliza-Version", "v1.0.0")
// 	return res, nil
// }

// func main() {
// 	// 建立服務實作
// 	server := &ElizaServer{}

// 	// 建立 Connect HTTP 處理器
// 	path, handler := elizav1connect.NewElizaServiceHandler(server)

// 	// 設定 CORS 中介軟體
// 	// 這是關鍵！允許來自 Next.js 開發伺服器 (e.g., localhost:3003) 的請求
// 	corsHandler := cors.New(cors.Options{
// 		AllowedMethods: []string{"POST", "GET"},
// 		AllowedOrigins: []string{"http://localhost:3003"}, // 根據你的 Next.js port 設定
// 		AllowedHeaders: []string{"Content-Type", "Connect-Protocol-Version"},
// 		ExposedHeaders: []string{"Connect-Content-Encoding", "Connect-Accept-Encoding"},
// 	}).Handler(handler)

// 	mux := http.NewServeMux()
// 	mux.Handle(path, corsHandler) // 使用帶有 CORS 的處理器

// 	fmt.Println("Go server listening on :8080...")
// 	err := http.ListenAndServe(":8080", mux)
// 	if err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }

package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	elizav1 "myapp/backend/gen/connectrpc/eliza"                     // proto messages
	elizav1connect "myapp/backend/gen/connectrpc/eliza/elizaconnect" // connect service handlers

	"connectrpc.com/connect"
	"github.com/rs/cors"
)

type ElizaServer struct{}

func (s *ElizaServer) Say(ctx context.Context, req *connect.Request[elizav1.SayRequest]) (*connect.Response[elizav1.SayResponse], error) {
	log.Println("收到請求：", req.Msg.Sentence)
	return connect.NewResponse(&elizav1.SayResponse{
		Sentence: "你說的是：" + req.Msg.Sentence,
	}), nil
}

func (s *ElizaServer) GetRandomPerson(ctx context.Context, req *connect.Request[elizav1.GetRandomPersonRequest]) (*connect.Response[elizav1.GetRandomPersonResponse], error) {
	log.Println("收到 GetRandomPerson 請求")

	// 預定義一些人員資料
	people := []*elizav1.Person{
		{
			Id:         "001",
			Name:       "張小明",
			Age:        28,
			Email:      "zhang.xiaoming@company.com",
			Department: "工程部",
			Position:   "軟體工程師",
		},
		{
			Id:         "002",
			Name:       "李小華",
			Age:        32,
			Email:      "li.xiaohua@company.com",
			Department: "產品部",
			Position:   "產品經理",
		},
		{
			Id:         "003",
			Name:       "王小美",
			Age:        25,
			Email:      "wang.xiaomei@company.com",
			Department: "設計部",
			Position:   "UI/UX 設計師",
		},
		{
			Id:         "004",
			Name:       "陳大強",
			Age:        35,
			Email:      "chen.daqiang@company.com",
			Department: "管理部",
			Position:   "技術總監",
		},
		{
			Id:         "005",
			Name:       "林小芳",
			Age:        29,
			Email:      "lin.xiaofang@company.com",
			Department: "行銷部",
			Position:   "行銷專員",
		},
	}

	// 隨機選擇一個人
	randomIndex := rand.Intn(len(people))
	randomPerson := people[randomIndex]

	log.Printf("回傳隨機人員：%s (%s)", randomPerson.Name, randomPerson.Position)

	return connect.NewResponse(&elizav1.GetRandomPersonResponse{
		Person: randomPerson,
	}), nil
}

func main() {
	mux := http.NewServeMux()
	path, handler := elizav1connect.NewElizaServiceHandler(&ElizaServer{})
	mux.Handle(path, handler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3003"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	log.Println("後端啟動在 :8080")
	http.ListenAndServe(":8080", c.Handler(mux))
}
