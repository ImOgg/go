# GitHub Actions 工作流程說明

本專案使用 GitHub Actions 進行自動化測試 (CI - Continuous Integration)。

---

## 目錄

- [工作流程檔案](#工作流程檔案)
- [觸發條件](#觸發條件)
- [測試流程說明](#測試流程說明)
- [環境配置](#環境配置)
- [如何使用](#如何使用)
- [常見問題](#常見問題)

---

## 工作流程檔案

**檔案位置：** `.github/workflows/test.yml`

這個工作流程會在每次 Pull Request 到 `main` 分支時自動執行測試，確保新的程式碼不會破壞現有功能。

---

## 觸發條件

```yaml
on:
  workflow_dispatch:  # 手動觸發
  pull_request:       # PR 到 main 分支時自動觸發
    branches:
      - main
```

### 觸發方式

1. **自動觸發：** 當建立 Pull Request 到 `main` 分支時
2. **手動觸發：** 在 GitHub Actions 頁面手動執行

---

## 測試流程說明

### 執行環境

- **作業系統：** Ubuntu Latest
- **PHP 版本：** 8.3
- **資料庫：** MySQL 8.0.35

### 服務容器 (Services)

工作流程會自動啟動 MySQL 容器作為測試資料庫：

```yaml
services:
  mysql:
    image: mysql:8.0.35
    env:
      MYSQL_DATABASE: laravel
      MYSQL_ROOT_PASSWORD: root
    ports:
      - 3306:3306
    options: --health-cmd="mysqladmin ping" --health-interval=10s --health-timeout=5s --health-retries=3
```

### 測試步驟

#### 1. 設定 PHP 環境
```yaml
- uses: shivammathur/setup-php@v2
  with:
    php-version: '8.3'
```
安裝 PHP 8.3 及必要的擴展。

#### 2. 檢出代碼
```yaml
- uses: actions/checkout@v4
```
從儲存庫下載最新的程式碼。

#### 3. 複製環境變數檔案
```yaml
- name: COPY .env
  run: cp .env.ci .env
```
使用 CI 專用的環境變數檔案 (`.env.ci`)，其中包含測試用的資料庫配置。

**`.env.ci` 應包含：**
```env
APP_ENV=testing
APP_KEY=
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=root
```

#### 4. 安裝 Composer 依賴
```yaml
- name: Install Dependencies
  run: composer install -q --no-ansi --no-interaction --no-scripts --no-progress --prefer-dist
```

**參數說明：**
- `-q`: 安靜模式（減少輸出）
- `--no-ansi`: 不使用 ANSI 顏色
- `--no-interaction`: 非互動模式
- `--no-scripts`: 不執行 Composer 腳本
- `--no-progress`: 不顯示進度條
- `--prefer-dist`: 優先使用發行版本（速度較快）

#### 5. 生成應用程式金鑰
```yaml
- name: Generate Application Key
  run: php artisan key:generate
```
生成 Laravel 應用程式加密金鑰。

#### 6. 執行資料庫遷移
```yaml
- name: Run Migrations
  run: php artisan migrate
```
在測試資料庫中建立所有資料表。

#### 7. 設定目錄權限
```yaml
- name: directory permission
  run: sudo chmod -R 777 storage bootstrap/cache
```
確保 Laravel 可以寫入 storage 和 cache 目錄。

#### 8. 執行測試
```yaml
- name: Run Tests
  run: php artisan test
```
執行所有的 PHPUnit 測試。

---

## 環境配置

### 必要檔案

在 `api/` 目錄下需要創建 `.env.ci` 檔案：

```bash
# 範例 .env.ci
APP_NAME=Laravel
APP_ENV=testing
APP_KEY=
APP_DEBUG=true
APP_TIMEZONE=UTC
APP_URL=http://localhost

LOG_CHANNEL=stack
LOG_STACK=single
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=root

BROADCAST_CONNECTION=log
FILESYSTEM_DISK=local
QUEUE_CONNECTION=sync
SESSION_DRIVER=array
SESSION_LIFETIME=120

CACHE_STORE=array
REDIS_CLIENT=phpredis

MAIL_MAILER=log
```

### 工作目錄設定

```yaml
defaults:
  run:
    working-directory: ./api
```

所有命令都在 `./api` 目錄下執行，因為 Laravel 專案位於該目錄。

---

## 如何使用

### 1. 建立 Pull Request

當你完成功能開發並推送到 GitHub 後：

```bash
git checkout -b feature/new-feature
git add .
git commit -m "Add new feature"
git push origin feature/new-feature
```

然後在 GitHub 上建立 Pull Request 到 `main` 分支。

### 2. 自動執行測試

Pull Request 建立後，GitHub Actions 會自動：
1. 啟動測試環境
2. 執行所有測試
3. 在 PR 頁面顯示測試結果

### 3. 查看測試結果

在 Pull Request 頁面可以看到：
- ✅ 綠色勾勾：測試通過
- ❌ 紅色叉叉：測試失敗

點擊 "Details" 可以查看詳細的測試日誌。

### 4. 手動觸發測試

你也可以手動觸發測試：

1. 前往 GitHub 儲存庫
2. 點擊 "Actions" 標籤
3. 選擇 "TESTS" 工作流程
4. 點擊 "Run workflow"
5. 選擇分支並執行

---

## 常見問題

### 1. 測試失敗：資料庫連線錯誤

**錯誤訊息：**
```
SQLSTATE[HY000] [2002] Connection refused
```

**解決方法：**
- 確認 `.env.ci` 中的資料庫設定正確
- MySQL 服務容器需要時間啟動，已設定 health check 等待

### 2. 測試失敗：缺少 .env.ci 檔案

**錯誤訊息：**
```
cp: cannot stat '.env.ci': No such file or directory
```

**解決方法：**
在 `api/` 目錄下建立 `.env.ci` 檔案並提交到儲存庫。

### 3. 測試失敗：APP_KEY 未設定

**錯誤訊息：**
```
No application encryption key has been specified.
```

**解決方法：**
工作流程會自動執行 `php artisan key:generate`，如果仍然失敗，檢查 `.env.ci` 是否有 `APP_KEY=` 這一行。

### 4. 權限問題

**錯誤訊息：**
```
file_put_contents(/path/to/storage): failed to open stream: Permission denied
```

**解決方法：**
工作流程已包含權限設定步驟：
```yaml
- name: directory permission
  run: sudo chmod -R 777 storage bootstrap/cache
```

如果仍有問題，檢查是否有其他目錄需要寫入權限。

### 5. Composer 安裝超時

**錯誤訊息：**
```
The operation timed out
```

**解決方法：**
- 使用 `--prefer-dist` 參數（已包含）
- 檢查是否有私有套件需要 token

---

## 擴展建議

### 1. 增加程式碼品質檢查

可以在測試之前加入程式碼風格檢查：

```yaml
- name: Run PHP CS Fixer
  run: ./vendor/bin/php-cs-fixer fix --dry-run --diff

- name: Run PHPStan
  run: ./vendor/bin/phpstan analyse
```

### 2. 增加測試覆蓋率報告

```yaml
- name: Run Tests with Coverage
  run: php artisan test --coverage

- name: Upload Coverage to Codecov
  uses: codecov/codecov-action@v3
```

### 3. 多版本 PHP 測試

```yaml
strategy:
  matrix:
    php-version: ['8.2', '8.3']
```

### 4. 快取 Composer 依賴

```yaml
- name: Cache Composer dependencies
  uses: actions/cache@v3
  with:
    path: vendor
    key: ${{ runner.os }}-composer-${{ hashFiles('**/composer.lock') }}
```

---

## 未實作功能

### CD (Continuous Deployment) - 持續部署

本專案**目前沒有實作自動部署**，原因如下：

1. **需要額外配置：**
   - SSH Key 設定
   - 伺服器 IP 位址
   - 部署腳本權限

2. **如要實作 CD，需要：**

#### 設定 GitHub Secrets

在儲存庫設定中新增：
- `SSH_PRIVATE_KEY`: 部署用的 SSH 私鑰
- `SERVER_HOST`: 伺服器 IP 或域名
- `SERVER_USER`: SSH 使用者名稱
- `SERVER_PATH`: 專案部署路徑

#### 建立部署工作流程

建立 `.github/workflows/deploy.yml`：

```yaml
name: Deploy to Production

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Server
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd ${{ secrets.SERVER_PATH }}
            git pull origin main
            cd api
            composer install --no-dev --optimize-autoloader
            php artisan migrate --force
            php artisan optimize
            sudo systemctl reload php8.3-fpm
            sudo systemctl reload nginx
```

**但目前專案不包含此功能，需要手動部署。**

---

## 相關資源

- [GitHub Actions 文件](https://docs.github.com/en/actions)
- [Laravel 測試文件](https://laravel.com/docs/10.x/testing)
- [PHPUnit 文件](https://phpunit.de/documentation.html)

---

**最後更新：** 2025-10-24
