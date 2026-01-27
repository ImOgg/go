---
name: laravel-patterns
description: Laravel architecture patterns, API design, database optimization, and server-side best practices for Laravel 11 development.
tools: Read, Grep, Glob, Bash
model: sonnet
---

# Laravel Development Patterns

Laravel 11 應用開發的架構模式和最佳實踐。

## 項目上下文

本項目是 Laravel 11 課程示範應用，核心功能：
- 文章通知生成並透過 Redis 隊列非同步處理
- Laravel Task Scheduler 執行排程任務
- Docker Compose 容器化部署
- GitHub Actions 自動化測試 (CI/CD)

## Service Layer Pattern

業務邏輯與數據訪問分離：

```php
// app/Services/NotificationService.php
class NotificationService
{
    public function __construct(
        private PostRepository $postRepo,
        private UserRepository $userRepo,
    ) {}

    public function generateNotifications(int $userCount, int $postCount): void
    {
        $users = $this->userRepo->getRandomUsers($userCount);
        $posts = $this->postRepo->getLatestPosts($postCount);

        foreach ($users as $user) {
            foreach ($posts as $post) {
                SendNewPostNotification::dispatch($user, $post);
            }
        }
    }
}
```

## Repository Pattern

數據訪問邏輯抽象：

```php
// app/Repositories/PostRepository.php
interface PostRepository
{
    public function all(): Collection;
    public function find(int $id): ?Post;
    public function create(array $data): Post;
    public function update(int $id, array $data): Post;
    public function delete(int $id): bool;
    public function getLatestPosts(int $limit = 10): Collection;
}

// app/Repositories/EloquentPostRepository.php
class EloquentPostRepository implements PostRepository
{
    public function getLatestPosts(int $limit = 10): Collection
    {
        return Post::latest()
            ->limit($limit)
            ->get();
    }
}
```

## Queue Jobs Pattern

非同步任務處理：

```php
// app/Jobs/SendNewPostNotification.php
class SendNewPostNotification implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    public function __construct(
        public User $user,
        public Post $post,
    ) {}

    public function handle(): void
    {
        // 生成通知邏輯
        $notification = $this->generateNotification();

        // 寫入存儲
        Storage::disk('public')->put(
            "notifications/{$this->user->id}/{$this->post->id}.json",
            json_encode($notification)
        );

        // 更新數據庫
        Notification::create([
            'user_id' => $this->user->id,
            'post_id' => $this->post->id,
            'data' => $notification,
        ]);
    }

    public function failed(\Throwable $exception): void
    {
        Log::error('Notification job failed', [
            'user_id' => $this->user->id,
            'post_id' => $this->post->id,
            'error' => $exception->getMessage(),
        ]);
    }
}
```

## Scheduling Tasks

定期執行任務：

```php
// routes/console.php
Schedule::command('inspire')->hourly();

// 自定義命令
Schedule::call(function () {
    // 每天執行一次清理任務
    Storage::disk('public')->deleteDirectory('temp');
})->daily();

// 執行 Artisan 命令
Schedule::command('generate:notifications 2 3 normal')->everyFiveMinutes();
```

## API Controller Pattern

RESTful API 端點：

```php
// app/Http/Controllers/PostController.php
class PostController extends Controller
{
    public function __construct(private PostService $postService) {}

    public function index()
    {
        $posts = $this->postService->getAllPosts();

        return response()->json([
            'success' => true,
            'data' => $posts,
        ]);
    }

    public function show(Post $post)
    {
        return response()->json([
            'success' => true,
            'data' => $post->load('user'),
        ]);
    }

    public function store(StorePostRequest $request)
    {
        $post = $this->postService->createPost($request->validated());

        return response()->json([
            'success' => true,
            'data' => $post,
        ], 201);
    }

    public function latestReport()
    {
        $report = $this->postService->generateLatestReport();

        return response()->json([
            'success' => true,
            'data' => $report,
        ]);
    }
}
```

## Model Relationships

模型關聯：

```php
// app/Models/Post.php
class Post extends Model
{
    protected $fillable = ['title', 'content', 'user_id'];
    protected $casts = ['created_at' => 'datetime'];

    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    public function notifications(): HasMany
    {
        return $this->hasMany(Notification::class);
    }
}

// app/Models/User.php
class User extends Model
{
    public function posts(): HasMany
    {
        return $this->hasMany(Post::class);
    }

    public function notifications(): HasMany
    {
        return $this->hasMany(Notification::class);
    }
}
```

## Database Migrations

遷移管理：

```php
// database/migrations/2024_01_01_000000_create_posts_table.php
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('posts', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->constrained()->onDelete('cascade');
            $table->string('title', 200);
            $table->text('content');
            $table->timestamps();

            // 索引
            $table->index('user_id');
            $table->index('created_at');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('posts');
    }
};
```

## Testing Pattern

Feature 測試：

```php
// tests/Feature/PostControllerTest.php
class PostControllerTest extends TestCase
{
    public function test_can_list_posts(): void
    {
        $posts = Post::factory(5)->create();

        $response = $this->getJson('/api/posts');

        $response->assertStatus(200)
                 ->assertJsonCount(5, 'data');
    }

    public function test_can_create_post(): void
    {
        $user = User::factory()->create();

        $response = $this->actingAs($user)
                        ->postJson('/api/posts', [
                            'title' => 'Test Post',
                            'content' => 'Test content',
                        ]);

        $response->assertStatus(201)
                 ->assertJsonPath('data.title', 'Test Post');

        $this->assertDatabaseHas('posts', [
            'title' => 'Test Post',
            'user_id' => $user->id,
        ]);
    }

    public function test_queue_job_sends_notification(): void
    {
        Queue::fake();

        $user = User::factory()->create();
        $post = Post::factory()->create();

        SendNewPostNotification::dispatch($user, $post);

        Queue::assertPushed(SendNewPostNotification::class);
    }
}
```

## Docker Support

### API 容器配置

```dockerfile
# Dockerfile - API Stage
FROM php:8.3-fpm as api

WORKDIR /usr/src

# 安裝擴展
RUN docker-php-ext-install pdo pdo_mysql bcmath

# 複製應用
COPY api .

# 安裝依賴
RUN composer install --no-dev

EXPOSE 9000
CMD ["php-fpm"]
```

### Worker 容器配置

```dockerfile
# Dockerfile - Worker Stage
FROM php:8.3-fpm as worker

WORKDIR /usr/src

# 安裝 supervisor
RUN apt-get update && apt-get install -y supervisor

COPY api .
RUN composer install --no-dev

COPY deployment/docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
```

## Environment Variables

```bash
# .env.example
APP_NAME=HiskioApp
APP_ENV=local
APP_DEBUG=true
APP_KEY=base64:...

# Database
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=hiskio
DB_USERNAME=root
DB_PASSWORD=

# Queue
QUEUE_CONNECTION=redis
REDIS_HOST=127.0.0.1
REDIS_PORT=6379

# Storage
FILESYSTEM_DISK=public
```

## Performance Best Practices

### Database Queries

```php
// ❌ BAD: N+1 query
$posts = Post::all();
foreach ($posts as $post) {
    echo $post->user->name;
}

// ✅ GOOD: Eager loading
$posts = Post::with('user')->get();
foreach ($posts as $post) {
    echo $post->user->name;
}

// ✅ GOOD: Select specific columns
$posts = Post::select('id', 'title', 'user_id')
    ->with(['user' => fn($q) => $q->select('id', 'name')])
    ->get();
```

### Caching

```php
// Cache queries
$posts = Cache::remember('posts.all', 3600, function () {
    return Post::with('user')->get();
});

// Cache results
$report = Cache::tags(['reports'])->remember('latest_report', 1800, function () {
    return $this->generateReport();
});

// Invalidate cache
Cache::tags(['reports'])->flush();
```

### Batch Operations

```php
// ✅ GOOD: Bulk insert
Post::insert($arrayOfPosts);

// ✅ GOOD: Batch update
Post::whereIn('id', $ids)->update(['status' => 'published']);

// ❌ BAD: Loop and save
foreach ($posts as $post) {
    $post->save();
}
```

## Common Artisan Commands

```bash
# 遷移
php artisan migrate
php artisan migrate:fresh
php artisan migrate:rollback

# 數據填充
php artisan db:seed
php artisan db:seed --class=PostSeeder

# 隊列
php artisan queue:work
php artisan queue:work --queue=high,default --tries=3

# 排程
php artisan schedule:run
php artisan schedule:list

# 命令
php artisan make:command GeneratePostNotifications
php artisan generate:notifications 2 3 normal

# 測試
php artisan test
php artisan test tests/Feature/PostControllerTest.php
php artisan test --coverage
```

## Code Organization

```
api/
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   ├── Requests/
│   │   └── Resources/
│   ├── Models/
│   ├── Services/
│   ├── Repositories/
│   ├── Jobs/
│   ├── Commands/
│   └── Events/
├── database/
│   ├── migrations/
│   └── seeders/
├── routes/
│   ├── api.php
│   ├── web.php
│   └── console.php
├── tests/
│   ├── Feature/
│   └── Unit/
└── storage/
    └── app/public/
```

---

**Remember**: Laravel's elegance comes from following its conventions and patterns. Clean architecture, proper separation of concerns, and testing from the start enable maintainable, scalable applications.
