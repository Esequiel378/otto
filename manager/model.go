package manager

import (
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/udhos/gwob"
)

const (
	FLOAT32_BYTES   = 4
	POSITION_FLOATS = 3
	TEXCOORD_FLOATS = 2
)

// Model represents a 3D model with its OpenGL buffers
type Model struct {
	Name     string
	VAO      uint32
	VBO      uint32
	EBO      uint32
	Indices  []uint32
	Vertices []float32
	Stride   int
	Bounds   mgl64.Vec3
	Volume   float64
}

// Texture represents a loaded texture
type Texture struct {
	Name string
	ID   uint32
}

// ModelManager manages loading and caching of 3D models and textures
type ModelManager struct {
	models   map[string]*Model
	textures map[string]*Texture
}

// NewModelManager creates a new instance of ModelManager
func NewModelManager() *ModelManager {
	return &ModelManager{
		models:   make(map[string]*Model),
		textures: make(map[string]*Texture),
	}
}

// Model retrieves a model by its name
func (m *ModelManager) Model(name string) (*Model, error) {
	model, exists := m.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s not found", name)
	}
	return model, nil
}

// Texture retrieves a texture by its name
func (m *ModelManager) Texture(name string) (*Texture, error) {
	texture, exists := m.textures[name]
	if !exists {
		return nil, fmt.Errorf("texture %s not found", name)
	}
	return texture, nil
}

// LoadModel loads a model from a file and caches it
func (m *ModelManager) LoadModel(name, filepath string) error {
	model, err := m.loadModelFromFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to load model %s from %s: %w", name, filepath, err)
	}

	model.Name = name
	m.models[name] = model
	return nil
}

// LoadTexture loads a texture from a file and caches it
func (m *ModelManager) LoadTexture(name, filepath string) error {
	textureID, err := m.loadTextureFromFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to load texture %s from %s: %w", name, filepath, err)
	}

	m.textures[name] = &Texture{
		Name: name,
		ID:   textureID,
	}
	return nil
}

// Init initializes the model manager by loading all models and textures from the specified paths
func (m *ModelManager) Init(modelsPath, texturesPath string) error {
	// Load models
	if err := m.loadModelsFromDirectory(modelsPath); err != nil {
		return fmt.Errorf("failed to load models: %w", err)
	}

	// Load textures
	if err := m.loadTexturesFromDirectory(texturesPath); err != nil {
		return fmt.Errorf("failed to load textures: %w", err)
	}

	return nil
}

// loadModelsFromDirectory loads all .obj files from a directory
func (m *ModelManager) loadModelsFromDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read models directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !isObjFile(filename) {
			continue
		}

		modelName := getFileNameWithoutExtension(filename)
		filepath := filepath.Join(path, filename)

		if err := m.LoadModel(modelName, filepath); err != nil {
			return err
		}
	}

	return nil
}

// loadTexturesFromDirectory loads all image files from a directory
func (m *ModelManager) loadTexturesFromDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read textures directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !isImageFile(filename) {
			continue
		}

		textureName := getFileNameWithoutExtension(filename)
		filepath := filepath.Join(path, filename)

		if err := m.LoadTexture(textureName, filepath); err != nil {
			return err
		}
	}

	return nil
}

// loadModelFromFile loads a single model from a file
func (m *ModelManager) loadModelFromFile(file string) (*Model, error) {
	options := &gwob.ObjParserOptions{
		LogStats: false,
	}

	objModel, err := gwob.NewObjFromFile(file, options)
	if err != nil {
		return nil, err
	}

	indices := make([]uint32, len(objModel.Indices))
	for i, val := range objModel.Indices {
		indices[i] = uint32(val)
	}

	var VAO, VBO, EBO uint32

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(objModel.Coord)*4, gl.Ptr(objModel.Coord), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	stride := int32(objModel.StrideSize)

	// Position attribute (location = 0)
	gl.VertexAttribPointerWithOffset(0, POSITION_FLOATS, gl.FLOAT, false, stride, 0)
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute (location = 1)
	if objModel.TextCoordFound {
		texOffset := POSITION_FLOATS * FLOAT32_BYTES
		gl.VertexAttribPointerWithOffset(1, TEXCOORD_FLOATS, gl.FLOAT, false, stride, uintptr(texOffset))
		gl.EnableVertexAttribArray(1)
	}

	// Normal attribute (location = 2) - comes after texture coordinates
	if objModel.NormCoordFound {
		normalOffset := (POSITION_FLOATS + TEXCOORD_FLOATS) * FLOAT32_BYTES
		gl.VertexAttribPointerWithOffset(2, POSITION_FLOATS, gl.FLOAT, false, stride, uintptr(normalOffset))
		gl.EnableVertexAttribArray(2)
	}

	gl.BindVertexArray(0)

	model := &Model{
		VAO:      VAO,
		VBO:      VBO,
		EBO:      EBO,
		Indices:  indices,
		Vertices: objModel.Coord,
		Stride:   objModel.StrideSize,
	}

	// Calculate bounds and volume
	model.Bounds = model.calculateBounds()
	model.Volume = model.calculateVolume()

	return model, nil
}

// loadTextureFromFile loads a single texture from a file
func (m *ModelManager) loadTextureFromFile(path string) (uint32, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open texture file: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return 0, fmt.Errorf("failed to decode JPEG: %v", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	return textureID, nil
}

// calculateBounds calculates the bounding box of the model
func (m *Model) calculateBounds() mgl64.Vec3 {
	if len(m.Vertices) == 0 {
		return mgl64.Vec3{0.5, 0.5, 0.5} // Default unit cube half-extents
	}

	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)

	// Convert stride from bytes to number of floats
	strideFloats := m.Stride / 4

	// Iterate through vertices (positions are first 3 floats of each vertex)
	for i := 0; i < len(m.Vertices); i += strideFloats {
		x := float64(m.Vertices[i])
		y := float64(m.Vertices[i+1])
		z := float64(m.Vertices[i+2])

		minX = math.Min(minX, x)
		minY = math.Min(minY, y)
		minZ = math.Min(minZ, z)

		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
		maxZ = math.Max(maxZ, z)
	}

	return mgl64.Vec3{
		(maxX - minX),
		(maxY - minY),
		(maxZ - minZ),
	}
}

// calculateVolume estimates the volume of the model using its bounding box
func (m *Model) calculateVolume() float64 {
	bounds := m.Bounds
	return bounds.X() * bounds.Y() * bounds.Z()
}

// Cleanup deletes all models and textures managed by the ModelManager
func (m *ModelManager) Cleanup() {
	for _, model := range m.models {
		gl.DeleteVertexArrays(1, &model.VAO)
		gl.DeleteBuffers(1, &model.VBO)
		gl.DeleteBuffers(1, &model.EBO)
	}
	m.models = make(map[string]*Model)

	for _, texture := range m.textures {
		gl.DeleteTextures(1, &texture.ID)
	}
	m.textures = make(map[string]*Texture)
}

// GetLoadedModels returns a list of all loaded model names
func (m *ModelManager) GetLoadedModels() []string {
	models := make([]string, 0, len(m.models))
	for name := range m.models {
		models = append(models, name)
	}
	return models
}

// GetLoadedTextures returns a list of all loaded texture names
func (m *ModelManager) GetLoadedTextures() []string {
	textures := make([]string, 0, len(m.textures))
	for name := range m.textures {
		textures = append(textures, name)
	}
	return textures
}

// Helper functions
func isObjFile(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".obj"
}

func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp"
}

func getFileNameWithoutExtension(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
