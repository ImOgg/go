# Agents & Skills 配置

本項目使用 Claude Code 的 agents 和 skills 系統來提供開發支持。

## 📍 配置位置

所有 agents 和 skills 定義在 `.claude/` 目錄中：
- `.claude/agents/` - 專業 agent 配置
- `.claude/skills/` - 開發技能和最佳實踐

## 🤖 Available Agents

### 後端開發
- **laravel-patterns** - Laravel 架構模式、最佳實踐
- **backend-patterns** - 通用後端設計模式

### 前端開發
- **vue3-nuxt** - Vue 3 / Nuxt 開發模式
- **frontend-patterns** - React 和現代前端模式

### 代碼質量
- **code-reviewer** - 代碼審查 (質量、安全、可維護性)
- **security-reviewer** - 安全漏洞檢測和修復
- **tdd-guide** - 測試驅動開發專家
- **build-error-resolver** - 構建和類型錯誤修復
- **refactor-cleaner** - 死代碼清理和重構

### 架構與規劃
- **architect** - 系統架構設計和技術決策
- **planner** - 複雜特性和重構規劃
- **doc-updater** - 文檔和代碼映射專家
- **e2e-runner** - E2E 測試和自動化

## 💡 Available Skills

### 安全與測試
- **security-review** - 安全檢查清單和模式
- **tdd-workflow** - TDD 工作流程

### 開發規範
- **laravel-security** - PHP/Laravel 安全最佳實踐
- **backend-patterns** - 後端架構模式
- **coding-standards** - 通用編碼標準
- **vue3-patterns** - Vue 3 / Nuxt 開發模式
- **frontend-patterns** - 前端設計模式
- **clickhouse-io** - ClickHouse 分析數據庫模式

## 🚀 使用方式

```bash
# Claude Code 會自動從 .claude/ 目錄讀取配置
# 無需額外配置
```

### 手動調用特定 Agent

```bash
claude run --task "描述你的任務"
# Claude 會根據任務內容自動選擇合適的 agent
```

### 使用特定 Skill

某些 skills 可以通過命令調用（如果已配置為 slash command）。

## 📝 配置信息

- **Agents**: 14 個（針對不同開發任務）
- **Skills**: 8 個（開發規範和最佳實踐指南）
- **默認模型**: sonnet（平衡速度和質量）
- **復雜任務**: 需要時自動升級到 opus

## 🔧 修改配置

要添加或修改 agents/skills：

1. 編輯 `.claude/agents/` 中的對應文件
2. 編輯 `.claude/skills/` 中的對應文件
3. 遵循現有的 YAML frontmatter 格式
4. 提交更改到 git

## 📚 相關文檔

- [CLAUDE.md](../CLAUDE.md) - 項目整體指南
- [.claude/agents/](../.claude/agents/) - Agent 配置詳情
- [.claude/skills/](../.claude/skills/) - Skill 配置詳情

---

**Generated with Claude Code** | 最後更新: 2026-01-22
