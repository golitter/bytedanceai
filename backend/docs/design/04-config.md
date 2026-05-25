# Config — 配置加载

## 实现了什么

通过 YAML 文件 + `.env` 环境变量双层机制加载配置，涵盖 MySQL、JWT、AgentEnd、七牛云、Redis 五个模块。敏感信息（七牛云密钥）从环境变量注入，不硬编码在 YAML 中。

## 怎么实现的

### Config 结构体 (`internal/conf/conf.go`)

```go
type Config struct {
	MySQL    MySQLConfig    `yaml:"mysql"`
	JWT      JWTConfig      `yaml:"jwt"`
	AgentEnd AgentEndConfig `yaml:"agentend"`
	Qiniu    QiniuConfig    `yaml:"qiniu"`
	Redis    RedisConfig    `yaml:"redis"`
}
```

各子配置结构体均提供辅助方法：

```go
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}
```

```go
type AgentEndConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *AgentEndConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
```

```go
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
```

```go
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type QiniuConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
	Region    string `yaml:"region"`
}
```

### 加载逻辑

`Load` 先尝试加载可选的 `.env` 文件，再读取 YAML 配置，最后用环境变量覆盖七牛云密钥：

```go
func Load(path string) (*Config, error) {
	_ = godotenv.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.Qiniu.AccessKey = os.Getenv("QINIU_ACCESS_KEY")
	cfg.Qiniu.SecretKey = os.Getenv("QINIU_SECRET_KEY")

	return &cfg, nil
}
```

### YAML 文件 (`configs/config.yaml`)

```yaml
mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: "123456"
  dbname: agenthub
  charset: utf8mb4

jwt:
  secret: agenthub-demo-secret
  expire_hours: 24

agentend:
  host: http://localhost
  port: 8001

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

qiniu:
  bucket: "agenthub"
  domain: "http://tfj4mvkda.hd-bkt.clouddn.com"
  region: z0    # z0=华东 z1=华北 z2=华南 na0=北美
```

七牛云 `access_key` / `secret_key` 不在 YAML 中配置，通过 `QINIU_ACCESS_KEY` / `QINIU_SECRET_KEY` 环境变量注入。
