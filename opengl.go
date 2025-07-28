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

// RenderEntity renders an entity using OpenGL 4.1 core profile
func RenderEntity(shaderManager *manager.ShaderManager, modelManager *manager.ModelManager, entity *physics.EntityRigidBody, camera *system.Camera) {
	// Skip rendering if no model name is specified (invisible entities like player, camera)
	if entity.ModelName == "" {
		return
	}

	shaderProgram, err := shaderManager.Program("camera")
	if err != nil {
		log.Printf("Failed to get shader program: %v", err)
		return
	}

	// Get the model
	model, err := modelManager.Model(entity.ModelName)
	if err != nil {
		log.Printf("Failed to get model %s: %v", entity.ModelName, err)
		return
	}

	// Use shader program
	gl.UseProgram(shaderProgram.PID)

	// Bind VAO
	gl.BindVertexArray(model.VAO)

	// Set up transformation matrices
	position := util.Vec64ToVec32(entity.Position)
	scale := util.Vec64ToVec32(entity.Scale)
	rotation := util.Vec64ToVec32(entity.Rotation)

	// Create model matrix
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(position.X(), position.Y(), position.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.Scale3D(scale.X(), scale.Y(), scale.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.X(), mgl32.Vec3{1, 0, 0}))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.Y(), mgl32.Vec3{0, 1, 0}))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3D(rotation.Z(), mgl32.Vec3{0, 0, 1}))

	// Create view matrix using camera data
	cameraPos := util.Vec64ToVec32(camera.Position)

	// Convert 2D rotation (pitch, yaw) to 3D direction vectors
	// camera.Rotation is mgl64.Vec2{pitch, yaw}
	pitch := camera.Rotation.X()
	yaw := camera.Rotation.Y()

	// Calculate forward direction from pitch and yaw
	cosPitch := float32(math.Cos(pitch))
	sinPitch := float32(math.Sin(pitch))
	cosYaw := float32(math.Cos(yaw))
	sinYaw := float32(math.Sin(yaw))

	forward := mgl32.Vec3{
		cosPitch * sinYaw,
		sinPitch,
		cosPitch * cosYaw,
	}

	// Calculate up vector (world up)
	up := mgl32.Vec3{0, 1, 0}

	// Calculate the target point (where camera is looking)
	target := cameraPos.Add(forward)

	// Create proper view matrix
	view := mgl32.LookAtV(
		cameraPos, // Camera position
		target,    // Look at point (camera position + forward direction)
		up,        // Up vector
	)

	// Create projection matrix with camera zoom
	fov := float32(45.0 / camera.Zoom) // Zoom affects field of view
	if fov < 5.0 {
		fov = 5.0
	}
	if fov > 90.0 {
		fov = 90.0
	}
	projection := mgl32.Perspective(mgl32.DegToRad(fov), 1200.0/900.0, 0.1, 10_000.0)

	// Set uniform matrices
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("model\x00")), 1, false, &modelMatrix[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("view\x00")), 1, false, &view[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("projection\x00")), 1, false, &projection[0])

	// Set material color (white)
	gl.Uniform4f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("color\x00")), 1.0, 1.0, 1.0, 1.0)

	// Set view position (needed for both vertex and fragment shaders)
	gl.Uniform3f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("viewPos\x00")), cameraPos.X(), cameraPos.Y(), cameraPos.Z())

	// Set ambient strength
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("ambientStrength\x00")), 0.3)

	// Set occlusion strength
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("occlusionStrength\x00")), 1.0)

	// Set up lighting (single light)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.PID, gl.Str("numLights\x00")), 1)

	// Light position - closer to the scene for better lighting
	lightPos := mgl32.Vec3{1.0, 1.0, 1.0}
	gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightPositions\x00")), 1, &lightPos[0])

	// Light color - brighter white light
	lightColor := mgl32.Vec3{1.0, 1.0, 1.0}
	gl.Uniform3fv(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightColors\x00")), 1, &lightColor[0])

	// Light intensity - increased for better visibility
	lightIntensity := float32(2.0)
	gl.Uniform1f(gl.GetUniformLocation(shaderProgram.PID, gl.Str("lightIntensities\x00")), lightIntensity)

	// Check for OpenGL errors
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error before drawing: %v", err)
	}

	// Draw the model
	gl.DrawElements(gl.TRIANGLES, int32(len(model.Indices)), gl.UNSIGNED_INT, nil)

	// Check for OpenGL errors after drawing
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error after drawing: %v", err)
	}

	// Unbind VAO
	gl.BindVertexArray(0)

	// Unuse shader program
	gl.UseProgram(0)
}

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
