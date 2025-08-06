package main

import (
	"github.com/asticode/go-astilectron"
	"log"
	"os"
)

func main() {
	// 1. 创建 logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// 2. 初始化 astilectron
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:           "Simple Notifier",
		BaseDirectoryPath: ".",
		// 添加以下配置强制使用64位
		ElectronSwitches: []string{"--arch=x64"},
		// 可选：使用更新的Electron版本
		VersionElectron: "28.1.0", // 最新LTS版本
	})
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	// 3. 启动 astilectron
	if err := a.Start(); err != nil {
		log.Fatal(err)
	}

	// 4. 创建窗口 (使用指针直接赋值)
	w, err := a.NewWindow("", &astilectron.WindowOptions{
		Width:       ptrInt(300),
		Height:      ptrInt(150),
		Center:      ptrBool(true),
		Frame:       ptrBool(false),
		Transparent: ptrBool(true),
		AlwaysOnTop: ptrBool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 5. 创建窗口
	if err := w.Create(); err != nil {
		log.Fatal(err)
	}

	// 6. 加载HTML内容
	html := `
	<div style="padding:20px;background:#4285f4;color:white;border-radius:5px;font-family:Arial;">
		<h3>微信消息</h3>
		<p>您收到一条新消息</p>
	</div>
	`
	if err := w.ExecuteJavaScript("document.body.innerHTML = `" + html + "`"); err != nil {
		log.Fatal(err)
	}

	// 7. 保持运行
	a.Wait()
}

// 辅助函数 - 创建int指针
func ptrInt(i int) *int {
	return &i
}

// 辅助函数 - 创建bool指针
func ptrBool(b bool) *bool {
	return &b
}
