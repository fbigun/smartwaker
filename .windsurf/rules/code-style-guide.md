---
trigger: always_on
description: SmartWaker 项目代码风格指南
---

# SmartWaker 项目代码风格指南

本文档定义了 SmartWaker 项目的代码风格和开发规范，旨在保持代码的一致性、可读性和可维护性。所有贡献者都应遵循这些规范。

## 目录

- [通用规范](#通用规范)
- [命名约定](#命名约定)
- [代码组织](#代码组织)
- [注释规范](#注释规范)
- [错误处理](#错误处理)
- [并发处理](#并发处理)
- [测试规范](#测试规范)
- [版本控制](#版本控制)

## 通用规范

1. **遵循 Go 官方规范**：遵循 [Effective Go](https://golang.org/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) 中的建议。

2. **代码格式化**：使用 `gofmt` 或 `go fmt` 格式化所有代码。提交前必须确保代码已格式化。

3. **代码简洁性**：遵循 KISS（Keep It Simple, Stupid）原则，避免不必要的复杂性。

4. **代码一致性**：保持项目中的代码风格一致，包括缩进、括号位置等。

5. **导入排序**：导入包应按以下顺序排列：
   - 标准库
   - 第三方库
   - 项目内部包

   各组之间用空行分隔，如：
   ```go
   import (
       "fmt"
       "log"
       "os"
       
       "github.com/eclipse/paho.mqtt.golang"
       
       "github.com/fbigun/smartwaker/internal/config"
   )
   ```

6. **避免魔法数字**：使用常量或变量替代代码中的魔法数字。

## 命名约定

1. **包名**：
   - 使用小写单词，不使用下划线或混合大小写。
   - 包名应简短且有意义，通常为单个单词。
   - 避免使用常见的变量名作为包名。

2. **文件名**：
   - 使用小写字母，可以包含下划线。
   - 文件名应反映其包含的主要内容。

3. **变量和函数**：
   - 使用驼峰命名法（camelCase）。
   - 局部变量应简短但有意义。
   - 导出的函数和变量（首字母大写）应有文档注释。

4. **常量**：
   - 使用全大写，单词间用下划线分隔（如 `MAX_RETRY_COUNT`）。
   - 相关常量应组织在一个 `const` 块中。

5. **接口名**：
   - 通常以 "er" 结尾（如 `Reader`, `Writer`）。
   - 单方法接口的名称应该是方法名加上 "er"（如 `Reader` 对应 `Read` 方法）。

6. **结构体**：
   - 使用驼峰命名法，首字母大写表示导出。
   - 字段名同样使用驼峰命名法。

## 代码组织

1. **项目结构**：
   - `/cmd`：主要应用程序入口点。
   - `/internal`：私有应用程序和库代码。
   - `/pkg`：可以被外部应用程序使用的库代码。
   - `/api`：API 定义和客户端库。
   - `/configs`：配置文件模板或默认配置。
   - `/test`：额外的外部测试应用程序和测试数据。

2. **包的组织**：
   - 按功能组织包，而不是按类型。
   - 避免循环依赖。

3. **文件组织**：
   - 相关的声明应该放在同一个文件中。
   - 文件不应过长，考虑拆分超过 500 行的文件。

## 注释规范

1. **包注释**：
   - 每个包都应有一个包注释，位于 `package` 语句之前。
   - 描述包的用途和提供的功能。

2. **函数和方法注释**：
   - 所有导出的函数和方法必须有注释。
   - 注释应以函数名开始，描述其功能、参数和返回值。

3. **变量和常量注释**：
   - 导出的变量和常量应有注释说明其用途。
   - 复杂的非导出变量也应有注释。

4. **代码内注释**：
   - 对于复杂的代码逻辑，添加内联注释解释实现细节。
   - 避免注释显而易见的内容。

5. **TODO 和 FIXME**：
   - 使用 `// TODO: ...` 标记待完成的工作。
   - 使用 `// FIXME: ...` 标记需要修复的问题。
   - 包含足够的上下文信息，以便其他开发者理解。

## 错误处理

1. **错误检查**：
   - 始终检查返回的错误。
   - 不要使用 `_` 忽略错误，除非有充分理由。

2. **错误传播**：
   - 使用 `fmt.Errorf("failed to ...: %w", err)` 包装错误并添加上下文。
   - 保持错误链，使用 `%w` 而不是 `%v` 来包装错误。

3. **错误日志**：
   - 在错误发生的地方记录详细信息。
   - 避免在多个地方记录同一个错误。

4. **自定义错误**：
   - 为特定错误定义自定义类型。
   - 实现 `Error()` 方法提供有意义的错误消息。

## 并发处理

1. **Goroutine 管理**：
   - 确保所有启动的 goroutine 都能正确退出。
   - 使用 context 包传递取消信号。

2. **共享资源访问**：
   - 使用适当的同步原语（如 `sync.Mutex`）保护共享资源。
   - 考虑使用 channel 进行通信而不是共享内存。

3. **并发模式**：
   - 优先使用 Go 的并发模式，如 fan-out/fan-in、worker pool 等。
   - 避免过度使用 goroutine，考虑资源限制。

## 测试规范

1. **测试覆盖率**：
   - 争取高测试覆盖率，特别是核心功能。
   - 使用 `go test -cover` 检查覆盖率。

2. **测试命名**：
   - 测试函数命名为 `Test<Function>`，如 `TestWakeOnLAN`。
   - 基准测试命名为 `Benchmark<Function>`。

3. **表驱动测试**：
   - 使用表驱动测试处理多个测试用例。
   - 为每个测试用例提供清晰的描述。

4. **模拟和存根**：
   - 使用接口进行依赖注入，便于测试时模拟依赖。
   - 考虑使用 testify 等测试辅助库。

5. **测试环境**：
   - 测试应该是可重复的，不依赖于特定环境。
   - 使用临时目录和文件进行文件系统测试。

## 版本控制

1. **提交信息**：
   - 提交信息应简洁明了，描述变更内容。
   - 使用现在时态（如 "Add feature" 而不是 "Added feature"）。

2. **分支策略**：
   - 主分支（main）应始终保持可部署状态。
   - 使用功能分支进行新功能开发。
   - 使用 Pull Request 进行代码审查。

3. **版本标签**：
   - 遵循语义化版本（SemVer）规范。
   - 为每个发布版本添加标签。

---

本指南将随着项目的发展而更新。如有疑问或建议，请提交 Issue 或 Pull Request。