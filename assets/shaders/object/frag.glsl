#version 410 core

in vec2 TexCoord;
in vec3 FragPos;
in vec3 Normal;

out vec4 FragColor;

uniform sampler2D textureColor;

// Lighting uniforms
uniform vec3 lightPos;      // Light position in world space
uniform vec3 lightColor;    // Light color
uniform vec3 viewPos;       // Camera position in world space

// Material properties
uniform float ambientStrength;
uniform float specularStrength;
uniform int shininess;

void main() {
    // Sample the texture
    vec3 textureColor = texture(textureColor, TexCoord).rgb;
    
    // Normalize the normal vector
    vec3 norm = normalize(Normal);
    
    // Ambient lighting
    vec3 ambient = ambientStrength * lightColor;
    
    // Diffuse lighting
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;
    
    // Specular lighting (Phong)
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
    vec3 specular = specularStrength * spec * lightColor;
    
    // Combine all lighting components
    vec3 result = (ambient + diffuse + specular) * textureColor;
    
    FragColor = vec4(result, 1.0);
}
