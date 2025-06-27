#version 410 core

in vec3 FragPos;
in vec3 Normal;

out vec4 FragColor;

uniform vec4 color;
uniform vec3 viewPos;
uniform float ambientStrength;
uniform float occlusionStrength;

#define MAX_LIGHTS 8
uniform int numLights;
uniform vec3 lightPositions[MAX_LIGHTS];
uniform vec3 lightColors[MAX_LIGHTS];
uniform float lightIntensities[MAX_LIGHTS];

void main() {
    vec3 norm = normalize(Normal);
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 result = vec3(0.0);

    for (int i = 0; i < numLights; ++i) {
        // Ambient
        vec3 ambient = ambientStrength * lightColors[i] * lightIntensities[i];

        // Diffuse
        vec3 lightDir = normalize(lightPositions[i] - FragPos);
        float diff = max(dot(norm, lightDir), 0.0);
        vec3 diffuse = diff * lightColors[i] * lightIntensities[i];

        // Specular
        float specularStrength = 0.5;
        vec3 reflectDir = reflect(-lightDir, norm);
        float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
        vec3 specular = specularStrength * spec * lightColors[i] * lightIntensities[i];

        result += ambient + diffuse + specular;
    }

    // If no lights, optionally keep ambient (or set to black)
    if (numLights == 0) {
        result = vec3(0.0); // or: result = ambientStrength * color.rgb;
    }

    // Apply material color and occlusion
    result = result * color.rgb;
    result = mix(result, result * occlusionStrength, 0.3); // Blend occlusion effect
    FragColor = vec4(result, color.a);
}
