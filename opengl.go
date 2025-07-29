package otto

import (
	"log"
	"math"
	"otto/manager"
	"otto/system"
	"otto/system/physics"
	"otto/util"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// RenderEntityBatch renders multiple entities of the same model type in a single batch
// This is more efficient than individual RenderEntity calls for many objects
func RenderEntityBatch(shaderManager *manager.ShaderManager, modelManager *manager.ModelManager, entities []*physics.EntityRigidBody, camera *system.Camera) {
	if len(entities) == 0 {
		return
	}

	// Group entities by model name for batch rendering
	modelGroups := make(map[string][]*physics.EntityRigidBody)
	for _, entity := range entities {
		if entity.ModelName == "" {
			continue // Skip invisible entities
		}
		modelGroups[entity.ModelName] = append(modelGroups[entity.ModelName], entity)
	}

	// Render each model group in batches
	for modelName, modelEntities := range modelGroups {
		shaderProgram, err := shaderManager.Program("camera")
		if err != nil {
			log.Printf("Failed to get shader program: %v", err)
			continue
		}

		model, err := modelManager.Model(modelName)
		if err != nil {
			log.Printf("Failed to get model %s: %v", modelName, err)
			continue
		}

		// Use shader program once for the entire batch
		gl.UseProgram(shaderProgram.PID)

		// Bind VAO once for the entire batch
		gl.BindVertexArray(model.VAO)

		// Set up view and projection matrices once (same for all entities)
		cameraPos := util.Vec64ToVec32(camera.Position)
		pitch := camera.Rotation.X()
		yaw := camera.Rotation.Y()

		cosPitch := float32(math.Cos(pitch))
		sinPitch := float32(math.Sin(pitch))
		cosYaw := float32(math.Cos(yaw))
		sinYaw := float32(math.Sin(yaw))

		forward := mgl32.Vec3{
			cosPitch * sinYaw,
			sinPitch,
			cosPitch * cosYaw,
		}

		up := mgl32.Vec3{0, 1, 0}
		target := cameraPos.Add(forward)
		view := mgl32.LookAtV(cameraPos, target, up)

		fov := float32(45.0 / camera.Zoom)
		if fov < 5.0 {
			fov = 5.0
		}
		if fov > 90.0 {
			fov = 90.0
		}
		projection := mgl32.Perspective(mgl32.DegToRad(fov), 1200.0/900.0, 0.1, 10_000.0)

		// Set view and projection uniforms once
		gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("projection\x00")), 1, false, &projection[0])

		// Set lighting uniforms once
		gl.Uniform4f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("color\x00")), 1.0, 1.0, 1.0, 1.0)
		gl.Uniform3f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("viewPos\x00")), cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("ambientStrength\x00")), 0.3)
		gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("occlusionStrength\x00")), 1.0)
		gl.Uniform1i(gl.GetUniformLocation(shaderProgram.PID, gl.Str("numLights\x00")), 1)

		lightPos := mgl32.Vec3{1.0, 1.0, 1.0}
		gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightPositions\x00")), 1, &lightPos[0])

		lightColor := mgl32.Vec3{1.0, 1.0, 1.0}
		gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightColors\x00")), 1, &lightColor[0])

		lightIntensity := float32(2.0)
		gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightIntensities\x00")), lightIntensity)

		// Render each entity in the batch
		renderedEntities := 0
		for _, entity := range modelEntities {
			// Simple frustum culling: skip entities too far from camera
			entityPos := util.Vec64ToVec32(entity.Position)
			distance := cameraPos.Sub(entityPos).Len()
			if distance > 1200.0 { // Skip entities more than 50 units away
				continue
			}

			// Set up transformation matrices for this entity
			position := util.Vec64ToVec32(entity.Position)
			scale := util.Vec64ToVec32(entity.Scale)
			rotation := util.Vec64ToVec32(entity.Rotation)

			modelMatrix := mgl32.Ident4()
			modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(position.X(), position.Y(), position.Z()))
			modelMatrix = modelMatrix.Mul4(mgl32.Scale3D(scale.X(), scale.Y(), scale.Z()))
			modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.X(), mgl32.Vec3{1, 0, 0}))
			modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.Y(), mgl32.Vec3{0, 1, 0}))
			modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.Z(), mgl32.Vec3{0, 0, 1}))

			// Set model matrix for this entity
			gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("model\x00")), 1, false, &modelMatrix[0])

			// Draw the model
			gl.DrawElements(gl.TRIANGLES, int32(len(model.Indices)), gl.UNSIGNED_INT, nil)
			renderedEntities++
		}

		// Unbind VAO
		gl.BindVertexArray(0)
	}

	// Unuse shader program
	gl.UseProgram(0)
}

// RenderGridFloor renders a grid floor at Y=0 using line primitives
func RenderGridFloor(shaderManager *manager.ShaderManager, modelManager *manager.ModelManager, floor *physics.EntityRigidBody, camera *system.Camera) {
	shaderProgram, err := shaderManager.Program("camera")
	if err != nil {
		log.Printf("Failed to get camera shader program: %v", err)
		return
	}

	// Use shader program
	gl.UseProgram(shaderProgram.PID)

	// Set up view and projection matrices
	cameraPos := util.Vec64ToVec32(camera.Position)
	pitch := camera.Rotation.X()
	yaw := camera.Rotation.Y()

	cosPitch := float32(math.Cos(pitch))
	sinPitch := float32(math.Sin(pitch))
	cosYaw := float32(math.Cos(yaw))
	sinYaw := float32(math.Sin(yaw))

	forward := mgl32.Vec3{
		cosPitch * sinYaw,
		sinPitch,
		cosPitch * cosYaw,
	}

	up := mgl32.Vec3{0, 1, 0}
	target := cameraPos.Add(forward)
	view := mgl32.LookAtV(cameraPos, target, up)

	fov := float32(45.0 / camera.Zoom)
	if fov < 5.0 {
		fov = 5.0
	}
	if fov > 90.0 {
		fov = 90.0
	}
	projection := mgl32.Perspective(mgl32.DegToRad(fov), 1200.0/900.0, 0.1, 10_000.0)

	// Set up model matrix to position grid at floor position
	floorPos := util.Vec64ToVec32(floor.Position)
	modelMatrix := mgl32.Translate3D(floorPos.X(), floorPos.Y(), floorPos.Z())

	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("model\x00")), 1, false, &modelMatrix[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("view\x00")), 1, false, &view[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("projection\x00")), 1, false, &projection[0])

	// Set lighting uniforms
	gl.Uniform4f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("color\x00")), 0.8, 0.8, 0.8, 1.0) // Light gray grid
	gl.Uniform3f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("viewPos\x00")), cameraPos.X(), cameraPos.Y(), cameraPos.Z())
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("ambientStrength\x00")), 0.3)
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("occlusionStrength\x00")), 1.0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.PID, gl.Str("numLights\x00")), 1)

	lightPos := mgl32.Vec3{1.0, 1.0, 1.0}
	gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightPositions\x00")), 1, &lightPos[0])

	lightColor := mgl32.Vec3{1.0, 1.0, 1.0}
	gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightColors\x00")), 1, &lightColor[0])

	lightIntensity := float32(2.0)
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightIntensities\x00")), lightIntensity)

	// Create grid data centered at floor position
	gridSize := float32(floor.Scale.X()) // Grid size in world units (100x100)
	divisions := int(floor.Scale.Z())    // Number of grid lines (100 divisions for 1-unit cells)
	step := gridSize / float32(divisions)
	halfSize := gridSize / 2.0

	var coordData []float32
	var indices []uint32
	vertexIndex := uint32(0)

	// Create horizontal lines (along X-axis)
	for i := 0; i <= divisions; i++ {
		z := float32(-halfSize + float32(i)*step)

		// Start point of line
		coordData = append(coordData, float32(-halfSize), 0.0, z, 0.0, 0.0, 0.0, 1.0, 0.0) // position, texcoord, normal
		indices = append(indices, vertexIndex)
		vertexIndex++

		// End point of line
		coordData = append(coordData, float32(halfSize), 0.0, z, 1.0, 0.0, 0.0, 1.0, 0.0) // position, texcoord, normal
		indices = append(indices, vertexIndex)
		vertexIndex++
	}

	// Create vertical lines (along Z-axis)
	for i := 0; i <= divisions; i++ {
		x := float32(-halfSize + float32(i)*step)

		// Start point of line
		coordData = append(coordData, x, 0.0, float32(-halfSize), 0.0, 0.0, 0.0, 1.0, 0.0) // position, texcoord, normal
		indices = append(indices, vertexIndex)
		vertexIndex++

		// End point of line
		coordData = append(coordData, x, 0.0, float32(halfSize), 0.0, 1.0, 0.0, 1.0, 0.0) // position, texcoord, normal
		indices = append(indices, vertexIndex)
		vertexIndex++
	}

	// Create OpenGL buffers
	var VAO, VBO, EBO uint32

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(coordData)*4, gl.Ptr(coordData), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	stride := int32(8 * 4) // 8 floats per vertex * 4 bytes per float

	// Position attribute (location = 0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, 0)
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute (location = 1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	// Normal attribute (location = 2)
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, stride, uintptr(5*4))
	gl.EnableVertexAttribArray(2)

	// Draw the grid using lines
	gl.DrawElements(gl.LINES, int32(len(indices)), gl.UNSIGNED_INT, nil)

	// Cleanup
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteBuffers(1, &EBO)

	// Unbind VAO
	gl.BindVertexArray(0)

	// Unuse shader program
	gl.UseProgram(0)
}
