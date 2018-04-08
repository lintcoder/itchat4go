package itchat4go

import (
    "fmt"
    "strings"
    "os"
    "image"
    "image/draw"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
)

const vertexShader = `
#version 330
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
out vec2 TexCoord;

void main()
{
	gl_Position = vec4(aPos, 1.0);
	TexCoord = vec2(aTexCoord.x, aTexCoord.y);
}
` + "\x00"

const fragmentShader = `
#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform sampler2D texture1;

void main()
{
	FragColor = texture(texture1, TexCoord);
}
` + "\x00"

func displayQRPic(status chan struct{}, file string) {
    defer chatter.wg.Done()
    if err := glfw.Init(); err != nil {
        fmt.Println("failed to init glfw:", err)
        return
    }
    defer glfw.Terminate()

    const windowWidth = 512
    const windowHeight = 512
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

    program, err := newProgram(vertexShader, fragmentShader)
    if err != nil {
        panic(err)
    }
    gl.UseProgram(program)

    vertices := []float32{0.5, 0.5, 0.0, 1.0, 1.0,
                          0.5, -0.5, 0.0, 1.0, 0.0,
                         -0.5, -0.5, 0.0, 0.0, 0.0,
                         -0.5, 0.5, 0.0, 0.0, 1.0}
    indices := []uint32{0, 1, 3, 1, 2, 3}
    var vbo, vao, ebo uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.GenBuffers(1, &vbo)
    gl.GenBuffers(1, &ebo)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW);
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW);

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
    gl.EnableVertexAttribArray(1)

    texture, err := newTexture(file)
    if err != nil {
        fmt.Println("load the texture failed ", file)
        return
    }

    flag := false
    for !flag {
        select {
        case <- status:
            flag = true
            break
        default:
            gl.ClearColor(1.0, 1.0, 1.0, 1.0)
            gl.Clear(gl.COLOR_BUFFER_BIT)
            gl.UseProgram(program)

            gl.BindTexture(gl.TEXTURE_2D, texture)
            gl.BindVertexArray(vao)
            gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

            window.SwapBuffers()
            glfw.PollEvents()
        }
    }
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
