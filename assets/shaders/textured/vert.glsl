#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec3 aNormal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec2 TexCoord;
out vec3 FragPos;
out vec3 Normal;

void main() {
    // Transform position to world space
    FragPos = vec3(model * vec4(aPos, 1.0));
    
    // Transform normal to world space (using normal matrix)
    Normal = mat3(transpose(inverse(model))) * aNormal;
    
    // Pass texture coordinates
    TexCoord = aTexCoord;
    
    // Transform to clip space
    gl_Position = projection * view * vec4(FragPos, 1.0);
} 