#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 2) in vec3 aNormal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform vec3 viewPos;

out vec3 FragPos;
out vec3 Normal;
out float FaceVisible;

void main() {
    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal;
    
    // Calculate face normal in world space
    vec3 faceNormal = normalize(Normal);
    
    // Calculate direction from fragment to camera
    vec3 viewDir = normalize(viewPos - FragPos);
    
    // Check if face is visible (dot product > 0 means face is pointing toward camera)
    // Use a small bias to avoid z-fighting at grazing angles
    float bias = 0.01;
    FaceVisible = step(-bias, dot(faceNormal, viewDir));
    
    gl_Position = projection * view * vec4(FragPos, 1.0);
}
