# gw

一个Go语言Windows平台GUI框架。

## 使用方法

```go
package main

import (
    "github.com/mkch/gw/app"
    "github.com/mkch/gw/button"
    "github.com/mkch/gw/win32"
    "github.com/mkch/gw/win32/win32util"
    "github.com/mkch/gw/window"
)

func main() {
    win, _ := window.New(&window.Spec{
        Text:  "Hello, Go!",
        Style: win32.WS_OVERLAPPEDWINDOW,
        X:     win32.CW_USEDEFAULT,
        Width: 500, Height: 300,
        OnClose: func() { app.Quit(0) },
    })
    button.New(win.HWND(), &button.Spec{
        Text:  "Hello",
        Style: win32.WS_VISIBLE,
        X:     200, Y: 120,
        Width: 100, Height: 60,
        OnClick: func() {
            win32util.MessageBox(win.HWND(),
                "Hello GUI!", "Button clicked",
                win32.MB_ICONINFORMATION)
        },
    })
    win.Show(win32.SW_SHOW)
    app.Run()
}

```

上述 Go 程序创建了一个窗口和一个按钮，当按钮被点击时弹出一个消息框。

## 常见问题

1. 如何去掉exe运行时的控制台窗口？

    在执行`go build`时加上`-ldflags "-H=windowsgui"`参数，例如：`go build -ldflags "-H=windowsgui"`。

2. 如何为exe指定图标等资源？

    可以使用第三方工具，例如`rsrc`。

    首先使用`go get github.com/akavel/rsrc`命令来安装`rsrc`。

    然后使用诸如`rsrc -arch amd64 -ico FILE.ico`的命令来把\*.ico文件编译为\*.syso资源文件。

    最后把\*.syso文件放在\*.go源文件放在同一个目录下，然后执行`go build`即可。
