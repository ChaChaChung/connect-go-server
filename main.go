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
	"database/sql"
	"log"
	"math/rand"
	"net/http"

	"myapp/backend/config"
	elizav1 "myapp/backend/gen/connectrpc/eliza"                     // proto messages
	elizav1connect "myapp/backend/gen/connectrpc/eliza/elizaconnect" // connect service handlers
	"myapp/backend/models"

	"connectrpc.com/connect"
	"github.com/rs/cors"
)

type ElizaServer struct {
	db             *sql.DB
	appearanceRepo models.AppearanceRepository
}

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

// GetAppearance 根據 ID 獲取外觀
func (s *ElizaServer) GetAppearance(ctx context.Context, req *connect.Request[elizav1.GetAppearanceRequest]) (*connect.Response[elizav1.GetAppearanceResponse], error) {
	log.Println("收到獲取外觀請求：", req.Msg.Id)

	appearance, err := s.appearanceRepo.GetByID(req.Msg.Id)
	if err != nil {
		log.Printf("獲取外觀失敗：%v", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// 轉換為 proto 消息
	protoAppearance := &elizav1.Appearance{
		Id:             appearance.ID,
		LogoUrl:        appearance.LogoURL.String,
		PrimaryColor:   appearance.PrimaryColor.String,
		SecondaryColor: appearance.SecondaryColor.String,
		CreatedAt:      appearance.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      appearance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LogoKey:        appearance.LogoKey.String,
		LogoSize:       appearance.LogoSize.Int64,
		LogoType:       appearance.LogoType.String,
	}

	return connect.NewResponse(&elizav1.GetAppearanceResponse{
		Appearance: protoAppearance,
	}), nil
}

// UpdateAppearance 更新外觀信息
func (s *ElizaServer) UpdateAppearance(ctx context.Context, req *connect.Request[elizav1.UpdateAppearanceRequest]) (*connect.Response[elizav1.UpdateAppearanceResponse], error) {
	log.Println("收到更新外觀請求：", req.Msg.Appearance.Id)

	// 先從資料庫讀取現有資料
	existingAppearance, err := s.appearanceRepo.GetByID(req.Msg.Appearance.Id)
	if err != nil {
		log.Printf("獲取現有外觀失敗：%v", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// 更新需要修改的欄位，保留 CreatedAt
	existingAppearance.LogoURL = sql.NullString{String: req.Msg.Appearance.LogoUrl, Valid: req.Msg.Appearance.LogoUrl != ""}
	existingAppearance.PrimaryColor = sql.NullString{String: req.Msg.Appearance.PrimaryColor, Valid: req.Msg.Appearance.PrimaryColor != ""}
	existingAppearance.SecondaryColor = sql.NullString{String: req.Msg.Appearance.SecondaryColor, Valid: req.Msg.Appearance.SecondaryColor != ""}
	existingAppearance.LogoKey = sql.NullString{String: req.Msg.Appearance.LogoKey, Valid: req.Msg.Appearance.LogoKey != ""}
	existingAppearance.LogoSize = sql.NullInt64{Int64: req.Msg.Appearance.LogoSize, Valid: req.Msg.Appearance.LogoSize > 0}
	existingAppearance.LogoType = sql.NullString{String: req.Msg.Appearance.LogoType, Valid: req.Msg.Appearance.LogoType != ""}
	// UpdatedAt 會在 repository 層自動更新

	// 更新到數據庫
	err = s.appearanceRepo.Update(existingAppearance)
	if err != nil {
		log.Printf("更新外觀失敗：%v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 轉換回 proto 消息
	protoAppearance := &elizav1.Appearance{
		Id:             existingAppearance.ID,
		LogoUrl:        existingAppearance.LogoURL.String,
		PrimaryColor:   existingAppearance.PrimaryColor.String,
		SecondaryColor: existingAppearance.SecondaryColor.String,
		CreatedAt:      existingAppearance.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      existingAppearance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LogoKey:        existingAppearance.LogoKey.String,
		LogoSize:       existingAppearance.LogoSize.Int64,
		LogoType:       existingAppearance.LogoType.String,
	}

	return connect.NewResponse(&elizav1.UpdateAppearanceResponse{
		Appearance: protoAppearance,
		Message:    "外觀更新成功",
	}), nil
}

func main() {
	// 初始化配置
	appConfig := config.NewAppConfig()

	// 連接到數據庫
	db, err := appConfig.Database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 創建外觀倉庫
	appearanceRepo := models.NewAppearanceRepository(db.DB)

	// 創建服務器實例
	server := &ElizaServer{
		db:             db.DB,
		appearanceRepo: appearanceRepo,
	}

	mux := http.NewServeMux()
	path, handler := elizav1connect.NewElizaServiceHandler(server)
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
