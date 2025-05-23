package config

type Config struct {
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Email   Email   `json:"email" yaml:"email"`
	ES      ES      `json:"es" yaml:"es"`
	Gaode   Gaode   `json:"gaode" yaml:"gaode"`
	Jwt     Jwt     `json:"jwt" yaml:"jwt"`
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Qiniu   Qiniu   `json:"qiniu" yaml:"qiniu"`
	QQ      QQ      `json:"qq" yaml:"qq"`
	Redis   Redis   `json:"redis" yaml:"redis"`
	System  System  `json:"system" yaml:"system"`
	Upload  Upload  `json:"upload" yaml:"upload"`
	Website Website `json:"website" yaml:"website"`
	Zap     Zap     `json:"zap" yaml:"zap"` // zap 是一个高性能的 Go 语言日志库，专为需要高效日志记录的应用程序设计。它提供了快速、结构化的日志记录功能，支持多种日志级别、日志滚动和多种输出格式（如 JSON）。InitLogger 函数通过配置初始化并返回一个 zap.Logger 实例，用于记录日志。
}

// 在Go语言中，yaml:"qiniu"是一个标签（tag），用于指定在YAML格式的配置文件中，Qiniu字段的名称应为qiniu。YAML是一种人类可读的数据序列化格式，常用于配置文件。YAML支持的数据类型包括：
// 标量类型：字符串、整数、浮点数、布尔值、null等。
// 序列类型：数组或列表。
// 映射类型：键值对，类似于JSON对象。
// 在YAML中，qiniu会被解析为一个键，其对应的值可以是上述任意一种YAML数据类型。
