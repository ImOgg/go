---
name: laravel-security
description: PHP and Laravel security best practices, OWASP vulnerabilities, and secure coding patterns for Laravel applications.
---

# Laravel Security Best Practices

PHP/Laravel 應用安全開發指南。

## 環境變量與密鑰管理

### ❌ NEVER

```php
// config/app.php
'key' => 'base64:xxxxxxxxxxxxx',  // 不要硬編碼

$db_password = 'mysql_password';   // 不要在代碼中

define('API_KEY', 'sk-xxxxxxx');   // 不要定義常量
```

### ✅ ALWAYS

```php
// 使用 .env 文件
DB_PASSWORD=secure_password
APP_KEY=base64:xxxxxxxxxxxxx
EXTERNAL_API_KEY=${EXTERNAL_API_KEY}

// config/app.php
'key' => env('APP_KEY'),

// 在代碼中使用
$apiKey = config('services.external.key');
$dbPassword = env('DB_PASSWORD');

// 驗證密鑰存在
if (!env('APP_KEY')) {
    throw new Exception('APP_KEY is not configured');
}
```

## SQL 注入防護

### ❌ VULNERABLE

```php
// 直接字符串拼接
$posts = DB::select("SELECT * FROM posts WHERE user_id = {$userId}");

// 使用字符串插值
$posts = Post::whereRaw("status = '$status'")->get();
```

### ✅ SECURE

```php
// 使用參數綁定
$posts = DB::select('SELECT * FROM posts WHERE user_id = ?', [$userId]);

// 使用 Query Builder
$posts = Post::where('user_id', $userId)->get();

// 使用 Eloquent
$posts = Post::whereUserId($userId)->get();

// 使用 Eloquent 關聯
$user = User::find($userId);
$posts = $user->posts()->get();
```

## XSS 防護

### ❌ VULNERABLE

```blade
{!! $userInput !!}
<div>{{ $post->content }}</div> {{-- 如果 content 包含 HTML --}}
```

### ✅ SECURE

```blade
{{-- Blade 默認轉義輸出 --}}
{{ $userInput }}
{{ $post->title }}

{{-- 如果需要轉義 HTML --}}
{{ htmlspecialchars($userInput, ENT_QUOTES, 'UTF-8') }}

{{-- 使用 HTML Purifier --}}
{{ Purifier::clean($userContent) }}
```

## CSRF 防護

### ✅ 自動 CSRF 保護

```blade
<form method="POST" action="/posts">
    @csrf
    <input type="text" name="title" />
    <button type="submit">Create</button>
</form>
```

```php
// 在控制器中驗證
class PostController extends Controller
{
    public function store(Request $request)
    {
        // 已自動驗證 CSRF 令牌（由 VerifyCsrfToken 中間件）
        $post = Post::create($request->validated());
        return redirect()->route('posts.show', $post);
    }
}
```

## 身份驗證與授權

### ✅ 安全的身份驗證

```php
// app/Http/Controllers/LoginController.php
class LoginController extends Controller
{
    public function login(Request $request)
    {
        $validated = $request->validate([
            'email' => 'required|email',
            'password' => 'required|min:8',
        ]);

        // 使用認證門面
        if (Auth::attempt($validated, $request->boolean('remember'))) {
            $request->session()->regenerate();
            return redirect()->route('home');
        }

        return back()->withErrors(['email' => 'Invalid credentials']);
    }

    public function logout(Request $request)
    {
        Auth::logout();
        $request->session()->invalidate();
        $request->session()->regenerateToken();

        return redirect('/');
    }
}
```

### ✅ 授權檢查

```php
// app/Policies/PostPolicy.php
class PostPolicy
{
    public function update(User $user, Post $post): bool
    {
        return $user->id === $post->user_id || $user->isAdmin();
    }

    public function delete(User $user, Post $post): bool
    {
        return $user->id === $post->user_id || $user->isAdmin();
    }
}

// 在控制器中使用
class PostController extends Controller
{
    public function update(Request $request, Post $post)
    {
        // 自動授權檢查
        $this->authorize('update', $post);

        $post->update($request->validated());
        return redirect()->route('posts.show', $post);
    }
}

// 在 Blade 中使用
@can('update', $post)
    <a href="{{ route('posts.edit', $post) }}">Edit</a>
@endcan
```

## 输入驗證

### ✅ 使用 Form Request

```php
// app/Http/Requests/StorePostRequest.php
class StorePostRequest extends FormRequest
{
    public function authorize(): bool
    {
        return auth()->check();
    }

    public function rules(): array
    {
        return [
            'title' => 'required|string|max:200',
            'content' => 'required|string|max:5000',
            'category_id' => 'required|exists:categories,id',
            'tags' => 'array|max:5',
            'tags.*' => 'string|max:50',
            'published_at' => 'nullable|date|after:today',
        ];
    }

    public function messages(): array
    {
        return [
            'title.required' => 'Title is required',
            'title.max' => 'Title cannot exceed 200 characters',
            'category_id.exists' => 'Selected category does not exist',
        ];
    }
}

// 在控制器中使用
class PostController extends Controller
{
    public function store(StorePostRequest $request)
    {
        // 已驗證和清理
        Post::create($request->validated());
        return redirect()->route('posts.index');
    }
}
```

### ✅ 文件上傳驗證

```php
class ProfileController extends Controller
{
    public function updateAvatar(Request $request)
    {
        $request->validate([
            'avatar' => 'required|image|mimes:jpeg,png,gif|max:2048',
        ]);

        // 檢查文件是否真的是圖片
        if ($request->file('avatar')->guessExtension() !== 'jpg') {
            return back()->withErrors(['avatar' => 'Invalid file']);
        }

        $path = $request->file('avatar')->store('avatars', 'public');

        auth()->user()->update(['avatar_path' => $path]);

        return back()->with('success', 'Avatar updated');
    }
}
```

## 密碼安全

### ✅ 密碼哈希與驗證

```php
// 新建用戶時哈希密碼
class User extends Model
{
    protected $hidden = ['password'];

    protected function password(): Attribute
    {
        return Attribute::make(
            set: fn(string $value) => Hash::make($value),
        );
    }
}

// 登錄時驗證
if (Hash::check($inputPassword, $storedHash)) {
    // 密碼正確
}

// 在請求驗證中
'password' => 'required|min:8|confirmed',
'password_confirmation' => 'required',
```

## 日誌與監控

### ✅ 安全的日誌記錄

```php
// ❌ 記錄敏感信息
Log::info('User login', ['password' => $password]);

// ✅ 只記錄必要信息
Log::info('User login', [
    'user_id' => $user->id,
    'email' => $user->email,
    'ip' => request()->ip(),
    'timestamp' => now(),
]);

// ❌ 記錄 API 密鑰
Log::info('API call', ['api_key' => env('EXTERNAL_API_KEY')]);

// ✅ 只記錄成功/失敗
Log::info('External API call successful', ['endpoint' => $endpoint]);
```

## 速率限制

### ✅ API 速率限制

```php
// routes/api.php
Route::post('/login', [LoginController::class, 'login'])
    ->middleware('throttle:5,1');  // 每分鐘 5 次

Route::post('/posts', [PostController::class, 'store'])
    ->middleware('throttle:30,1');  // 每分鐘 30 次

// 自定義速率限制
Route::group(['middleware' => 'throttle:api'], function () {
    Route::post('/api/trade', [TradeController::class, 'store']);
});

// config/cache.php
'api' => '60,1',  // 每分鐘 60 次
```

## 數據庫安全

### ✅ Row Level Security (Supabase)

```php
// 即使是 Laravel，也應確保用戶只能訪問自己的數據
class PostPolicy
{
    public function view(User $user, Post $post): bool
    {
        return $user->id === $post->user_id || $user->isAdmin();
    }
}
```

## HTTP 安全頭

### ✅ 在 Middleware 中設置

```php
// app/Http/Middleware/SecurityHeaders.php
class SecurityHeaders
{
    public function handle(Request $request, Closure $next)
    {
        $response = $next($request);

        // 防止 Clickjacking
        $response->header('X-Frame-Options', 'SAMEORIGIN');

        // 防止 MIME 類型嗅探
        $response->header('X-Content-Type-Options', 'nosniff');

        // 啟用 XSS 保護
        $response->header('X-XSS-Protection', '1; mode=block');

        // Content Security Policy
        $response->header(
            'Content-Security-Policy',
            "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
        );

        // HSTS (僅 HTTPS)
        if ($request->secure()) {
            $response->header('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
        }

        return $response;
    }
}
```

## 隊列安全

### ✅ 敏感數據在隊列中

```php
// app/Jobs/SendNotificationJob.php
class SendNotificationJob implements ShouldQueue
{
    // 不要存儲密碼或令牌
    public function __construct(
        protected User $user,  // ✅ 存儲模型 ID
        protected Post $post,
    ) {}

    public function handle()
    {
        // 重新查詢用戶以獲取最新數據
        $user = User::findOrFail($this->user->id);

        // 安全操作
        Notification::create([
            'user_id' => $user->id,
            'post_id' => $this->post->id,
        ]);
    }
}
```

## 依賴安全

### ✅ 保持依賴更新

```bash
# 檢查漏洞
composer audit

# 自動修復
composer update

# 檢查特定包的版本
composer show package-name
```

## 預部署安全檢查清單

```bash
# ✅ 檢查清單

# 1. 環境變量
- [ ] .env 文件不在版本控制中
- [ ] 所有敏感數據都在 .env 中
- [ ] 生產環境使用強 APP_KEY

# 2. 身份驗證
- [ ] 使用安全的密碼哈希
- [ ] CSRF 保護啟用
- [ ] 會話配置安全

# 3. 授權
- [ ] 所有受保護路由都有授權檢查
- [ ] 策略正確實現
- [ ] API 令牌驗證正確

# 4. 輸入驗證
- [ ] 所有用戶輸入都驗證
- [ ] 文件上傳受限
- [ ] 表單請求類完整

# 5. 輸出編碼
- [ ] 用戶數據在 Blade 中轉義
- [ ] 不使用 {!! !!} 用於用戶內容
- [ ] API 響應正確編碼

# 6. 依賴
- [ ] composer audit 通過
- [ ] 沒有已知漏洞

# 7. 日誌
- [ ] 沒有敏感數據在日誌中
- [ ] 日誌安全存儲

# 8. HTTPS
- [ ] 生產環境使用 HTTPS
- [ ] 重定向 HTTP 到 HTTPS
- [ ] 設置 HSTS 頭

# 9. 速率限制
- [ ] 登錄端點限制
- [ ] API 端點限制
- [ ] 敏感操作限制

# 10. 錯誤處理
- [ ] 生產環境 APP_DEBUG=false
- [ ] 沒有詳細的錯誤信息暴露
```

---

**Remember**: 安全是一個過程，不是目的地。定期審計代碼、保持依賴更新、並遵循最佳實踐。
