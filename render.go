package main

import (
	"errors"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const vertexShaderSrc = `#version 330 core
layout (location = 0) in vec3 aPos;
void main() {
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}` + "\x00"

const fragmentShaderSrc = `#version 330 core
out vec4 FragColor;
void main() {
	FragColor = vec4(1.0f, 0.0f, 0.0f, 1.0f);
}` + "\x00"

const yellowFragSrc = `#version 330 core
out vec4 FragColor;
void main() {
	FragColor = vec4(1.0f, 1.0f, 0.0f, 1.0f);
}` + "\x00"

func init() {
	runtime.LockOSThread()
}

func compileShader(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	srcGlString, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, srcGlString, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var msgLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &msgLength)
		logBuf := strings.Repeat("\x00", int(msgLength+1))
		gl.GetShaderInfoLog(shader, msgLength, nil, gl.Str(logBuf))

		return 0, errors.New("couldn't compile shader: " + logBuf)
	}

	return shader, nil
}

func createProgram(vertexSrc string, fragmentSrc string) (uint32, error) {
	vertexShader, err := compileShader(vertexSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var msgLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &msgLength)
		logBuf := strings.Repeat("\x00", int(msgLength+1))
		gl.GetProgramInfoLog(shaderProgram, msgLength, nil, gl.Str(logBuf))

		return 0, errors.New("couldn't link shaders: " + logBuf)
	}
	return shaderProgram, nil
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 640, "Hello world!", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	redShaderProgram, err := createProgram(vertexShaderSrc, fragmentShaderSrc)
	if err != nil {
		panic(err)
	}
	yellowShaderProgram, err := createProgram(vertexShaderSrc, yellowFragSrc)
	if err != nil {
		panic(err)
	}

	vertices := [...]float32{
		0.5, 0.5, 0.0, // top right
		0.5, -.5, 0.0, // bottom right
		-.5, -.5, 0.0, // bottom left
		1.0, -.5, 0.0, // further bottom right
	}
	indices := [...]uint32{
		0, 1, 2,
		0, 1, 3,
	}

	var VBO, VAO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)
	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, unsafe.Pointer(&indices[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.ClearColor(0.0, 1.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindVertexArray(VAO)
		gl.UseProgram(redShaderProgram)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.UseProgram(yellowShaderProgram)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(3*4))

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteBuffers(1, &EBO)

	glfw.Terminate()
}
