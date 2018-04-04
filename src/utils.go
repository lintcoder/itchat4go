package itchat4go

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
)

func displayQRPic(status chan struct{}, file string) {

    if err := glfw.Init(); err != nil {
        fmt.Println("failed to init glfw:", err)
        return
    }
    defer glfw.Terminate()

    const windowWidth = 640
    const windowHeight = 480
    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    window, err := glfw.CreateWindow(windowWidth, windowHeight, file, nil, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    window.MakeContextCurrent()
    if err := gl.Init(); err != nil {
        fmt.Println(err)
        return
    }

    flag := false
    for !flag {
        select {
        case <- status:
            flag = true
            break
        }
    }
    chatter.wg.Done()
}
