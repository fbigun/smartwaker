package tests

import (
	"os"
	"testing"
)

// TestMain 是测试的入口点
func TestMain(m *testing.M) {
	// 在所有测试开始前的准备工作
	setup()

	// 运行测试
	code := m.Run()

	// 在所有测试结束后的清理工作
	teardown()

	// 退出，返回测试结果
	os.Exit(code)
}

// setup 在测试开始前执行的准备工作
func setup() {
	// 可以在这里进行全局测试环境的设置
	// 例如创建临时文件、设置环境变量等
}

// teardown 在测试结束后执行的清理工作
func teardown() {
	// 可以在这里进行全局测试环境的清理
	// 例如删除临时文件、恢复环境变量等
}
