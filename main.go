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
	"github.com/joho/godotenv"
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

// GetAppearance 獲取第一筆外觀資料
func (s *ElizaServer) GetAppearance(ctx context.Context, req *connect.Request[elizav1.GetAppearanceRequest]) (*connect.Response[elizav1.GetAppearanceResponse], error) {
	log.Println("收到獲取外觀請求")

	appearance, err := s.appearanceRepo.GetFirst()
	if err != nil {
		log.Printf("獲取外觀失敗：%v", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// 轉換為 proto 消息，處理空值
	protoAppearance := &elizav1.Appearance{
		Id: appearance.ID,
	}

	// 安全地處理可能為空的字符串字段
	if appearance.LogoURL.Valid {
		protoAppearance.LogoUrl = appearance.LogoURL.String
	}
	if appearance.PrimaryColor.Valid {
		protoAppearance.PrimaryColor = appearance.PrimaryColor.String
	}
	if appearance.SecondaryColor.Valid {
		protoAppearance.SecondaryColor = appearance.SecondaryColor.String
	}
	if appearance.LogoKey.Valid {
		protoAppearance.LogoKey = appearance.LogoKey.String
	}
	if appearance.LogoType.Valid {
		protoAppearance.LogoType = appearance.LogoType.String
	}

	// 安全地處理可能為空的整數字段
	if appearance.LogoSize.Valid {
		protoAppearance.LogoSize = appearance.LogoSize.Int64
	}

	// 時間字段格式化
	protoAppearance.CreatedAt = appearance.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	protoAppearance.UpdatedAt = appearance.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")

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
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

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

	// 創建 HTTP 路由
	mux := http.NewServeMux()
	path, handler := elizav1connect.NewElizaServiceHandler(server)
	mux.Handle(path, handler)

	// 設置 CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3003"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// 啟動 HTTP 服務器
	log.Println("後端啟動在 :8080")
	http.ListenAndServe(":8080", c.Handler(mux))
}
