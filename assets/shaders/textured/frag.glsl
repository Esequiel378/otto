#version 410 core

in vec2 TexCoord;
in vec3 FragPos;
in vec3 Normal;

out vec4 FragColor;

uniform sampler2D textureColor;
uniform vec4 color;
uniform vec3 lightPos;
uniform vec3 lightColor;
uniform vec3 viewPos;
uniform float ambientStrength;
uniform float occlusionStrength;
uniform bool useTexture;

void main() {
    vec3 materialColor;
    
    if (useTexture) {
        // Sample the texture
        materialColor = texture(textureColor, TexCoord).rgb;
    } else {
        // Use the color uniform
        materialColor = color.rgb;
    }
    
    // Normalize the normal vector
    vec3 norm = normalize(Normal);
    
    // Ambient lighting
    vec3 ambient = ambientStrength * lightColor;
    
    // Diffuse lighting
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;
    
    // Specular lighting (Phong)
    float specularStrength = 0.5;
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    vec3 specular = specularStrength * spec * lightColor;
    
    // Combine all lighting components
    vec3 result = (ambient + diffuse + specular) * materialColor;
    
    // Apply ambient occlusion effect
    result = mix(result, result * occlusionStrength, 0.3);
    
    FragColor = vec4(result, color.a);
} 