# Connect Go Server with PostgreSQL CRUD

這是一個基於 Connect-Go 和 PostgreSQL 的 API 服務器，提供完整的 CRUD 操作。

## 功能特性

- ✅ 現有的 Eliza 服務（Say, GetRandomPerson）
- ✅ 外觀數據的獲取和更新操作
- ✅ PostgreSQL 數據庫集成
- ✅ Connect-Go 協議支持
- ✅ CORS 支持

## 數據庫配置

在運行服務器之前，請確保設置以下環境變量：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=your_database_name
export DB_SSLMODE=disable
```

或者修改 `config/database.go` 中的默認值。

## 安裝依賴

```bash
go mod tidy
```

## 運行服務器

```bash
go run main.go
```

服務器將在 `:8080` 端口啟動。

## API 端點

### 現有端點

- `POST /connectrpc.eliza.v1.ElizaService/Say` - 回顯輸入的句子
- `POST /connectrpc.eliza.v1.ElizaService/GetRandomPerson` - 獲取隨機人員信息

### 新增的外觀相關端點

- `POST /connectrpc.eliza.v1.ElizaService/GetAppearance` - 根據 ID 獲取外觀
- `POST /connectrpc.eliza.v1.ElizaService/UpdateAppearance` - 更新外觀信息

## 數據庫表結構

服務器期望有一個名為 `appearance` 的表，包含以下字段：

```sql
CREATE TABLE appearance (
    id VARCHAR(255) PRIMARY KEY,
    logo_url TEXT,
    primary_color VARCHAR(50),
    secondary_color VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    logo_key VARCHAR(255),
    logo_size BIGINT,
    logo_type VARCHAR(100)
);
```

## 使用示例

### 獲取外觀

```bash
curl -X POST http://localhost:8080/connectrpc.eliza.v1.ElizaService/GetAppearance \
  -H "Content-Type: application/json" \
  -d '{
    "id": "appearance_001"
  }'
```

### 更新外觀

```bash
curl -X POST http://localhost:8080/connectrpc.eliza.v1.ElizaService/UpdateAppearance \
  -H "Content-Type: application/json" \
  -d '{
    "appearance": {
      "id": "appearance_001",
      "logo_url": "https://example.com/new-logo.png",
      "primary_color": "#FF6B6B",
      "secondary_color": "#4ECDC4",
      "logo_key": "logo_key_001",
      "logo_size": 1024,
      "logo_type": "image/png"
    }
  }'
```

## 項目結構

```
connect-go-server/
├── config/           # 配置管理
│   ├── app.go       # 應用配置
│   └── database.go  # 數據庫配置
├── models/           # 數據模型
│   ├── company.go   # 公司模型和倉庫（已棄用）
│   └── appearance.go # 外觀模型和倉庫
├── proto/            # Protocol Buffers 定義
│   └── connectrpc/eliza/eliza.proto
├── gen/              # 生成的代碼
├── main.go           # 主服務器文件
└── README.md         # 本文檔
```

## 開發說明

- 修改 `proto/connectrpc/eliza/eliza.proto` 後，運行 `buf generate` 重新生成代碼
- 外觀相關的數據庫操作在 `models/appearance.go` 中實現
- 服務器邏輯在 `main.go` 中的 `ElizaServer` 結構體中實現
