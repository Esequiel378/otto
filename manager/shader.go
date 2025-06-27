package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderFile struct {
	Name    string
	Path    string
	Content string
	Type    uint32
}

type ShaderProgram struct {
	Name string
	PID  uint32
}

type ShaderManager struct {
	programs map[string]ShaderProgram
}

// NewShaderManager creates a new instance of ShaderManager.
func NewShaderManager() *ShaderManager {
	instance := &ShaderManager{
		programs: make(map[string]ShaderProgram),
	}
	return instance
}

// Program retrieves a shader program by its name.
func (s *ShaderManager) Program(name string) (ShaderProgram, error) {
	program, exists := s.programs[name]
	if !exists {
		return ShaderProgram{}, fmt.Errorf("shader program %s not found", name)
	}
	return program, nil
}

// Cleanup deletes all shader programs managed by the ShaderManager.
func (s *ShaderManager) Cleanup() {
	for _, program := range s.programs {
		gl.DeleteProgram(program.PID)
	}
	s.programs = make(map[string]ShaderProgram)
}

// Init initializes the shader manager by loading all shader files from the specified path.
// Each subdirectory will be treated as a separate shader program, and the shader files
func (s *ShaderManager) Init(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read root directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		folderName := entry.Name()
		folderPath := filepath.Join(path, folderName)

		shaderFiles, err := s.loadShadersFromFolder(folderPath)
		if err != nil {
			return err
		}

		if len(shaderFiles) == 0 {
			return fmt.Errorf("no shader files found in %s", folderPath)
		}

		programHandle, err := s.createProgram(shaderFiles)
		if err != nil {
			return err
		}

		s.programs[folderName] = ShaderProgram{
			Name: folderName,
			PID:  programHandle,
		}
	}

	return nil
}

func (s *ShaderManager) loadShadersFromFolder(folderPath string) ([]ShaderFile, error) {
	var shaders []ShaderFile

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %s: %w", folderPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if !strings.HasSuffix(strings.ToLower(fileName), ".glsl") {
			continue
		}

		filePath := filepath.Join(folderPath, fileName)

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read shader file %s: %w", filePath, err)
		}

		shaderType := s.determineShaderType(fileName, string(content))

		name := strings.TrimSuffix(fileName, ".glsl")

		shaders = append(shaders, ShaderFile{
			Name:    name,
			Path:    filePath,
			Content: string(content),
			Type:    shaderType,
		})
	}

	return shaders, nil
}

func (s *ShaderManager) determineShaderType(filename, content string) uint32 {
	lower := strings.ToLower(filename)

	if strings.Contains(lower, "vertex") || strings.Contains(lower, "vert") {
		return gl.VERTEX_SHADER
	}
	if strings.Contains(lower, "fragment") || strings.Contains(lower, "frag") {
		return gl.FRAGMENT_SHADER
	}
	if strings.Contains(lower, "geometry") || strings.Contains(lower, "geom") {
		return gl.GEOMETRY_SHADER
	}
	if strings.Contains(lower, "compute") || strings.Contains(lower, "comp") {
		return gl.COMPUTE_SHADER
	}
	if strings.Contains(lower, "tess_ctrl") || strings.Contains(lower, "tesc") {
		return gl.TESS_CONTROL_SHADER
	}
	if strings.Contains(lower, "tess_eval") || strings.Contains(lower, "tese") {
		return gl.TESS_EVALUATION_SHADER
	}

	contentLower := strings.ToLower(content)
	if strings.Contains(contentLower, "gl_position") {
		return gl.VERTEX_SHADER
	}
	if strings.Contains(contentLower, "gl_fragcolor") || strings.Contains(contentLower, "gl_fragdata") {
		return gl.FRAGMENT_SHADER
	}

	return gl.VERTEX_SHADER
}

func (s *ShaderManager) createProgram(shaderFiles []ShaderFile) (uint32, error) {
	program := gl.CreateProgram()
	var shaderHandles []uint32

	for _, shaderFile := range shaderFiles {
		shader := gl.CreateShader(shaderFile.Type)

		csource, free := gl.Strs(shaderFile.Content + "\x00")
		gl.ShaderSource(shader, 1, csource, nil)
		free()

		gl.CompileShader(shader)

		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

			return 0, fmt.Errorf("failed to compile shader %s: %v", shaderFile.Path, log)
		}

		gl.AttachShader(program, shader)
		shaderHandles = append(shaderHandles, shader)
	}

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

	for _, shader := range shaderHandles {
		gl.DeleteShader(shader)
	}

	return program, nil
}
