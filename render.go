package main

import (
	"errors"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const vertexShaderSrc = `#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 ourColor;

void main() {
	gl_Position = projection * view * model * vec4(aPos, 1.0);
	ourColor = aColor;
}` + "\x00"

const fragmentShaderSrc = `#version 330 core
in vec3 ourColor;
out vec4 FragColor;
void main() {
	FragColor = vec4(ourColor, 1.0f);
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
		return 0, err
	}
	fragmentShader, err := compileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
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

	shaderProgram, err := createProgram(vertexShaderSrc, fragmentShaderSrc)
	if err != nil {
		panic(err)
	}

	vertices := [...]float32{
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0, 0.0,

		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,

		0.5, 0.5, 0.5, 1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 1.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 1.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 1.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 1.0,
	}

	cubePositions := [...]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{2, 5, -15},
		mgl32.Vec3{-1.5, -2.2, -2.5},
		mgl32.Vec3{-3.8, -2, -12.3},
		mgl32.Vec3{2.4, -0.4, -3.5},
		mgl32.Vec3{-1.7, 3, -7.5},
		mgl32.Vec3{1.3, -2, -2.5},
		mgl32.Vec3{1.5, 2, -2.5},
		mgl32.Vec3{1.5, 0.2, -1.5},
		mgl32.Vec3{-1.3, 1, -1.5},
	}

	gl.Enable(gl.DEPTH_TEST)

	var VBO, VAO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)
	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	modelLocation := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(shaderProgram, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))

	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.BindVertexArray(VAO)
		gl.UseProgram(shaderProgram)

		view := mgl32.Ident4().Mul4(mgl32.Translate3D(0.0, 0.0, -3.0))
		projection := mgl32.Ident4().Mul4(mgl32.Perspective(mgl32.DegToRad(45.0), 800.0/640.0, 0.1, 100.0))
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
		gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

		for i := 0; i < 10; i++ {
			angle := 20.0*float32(i) + float32(glfw.GetTime())*50.0
			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(cubePositions[i][0], cubePositions[i][1], cubePositions[i][2]))
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{0.5, 1.0, 0.0}.Normalize()))
			gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteBuffers(1, &EBO)

	glfw.Terminate()
}
