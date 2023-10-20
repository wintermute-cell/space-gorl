#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;
uniform float ambient_light_level;

#define BLUR_AMOUNT 2

void main(void) {
    vec2 texelSize = 1.0 / vec2(textureSize(texture0, 0));
    
    vec4 color = vec4(0.0);
    int samples = 0;
    
    for(int x = -BLUR_AMOUNT; x <= BLUR_AMOUNT; x++) {
        for(int y = -BLUR_AMOUNT; y <= BLUR_AMOUNT; y++) {
            color += texture(texture0, fragTexCoord + vec2(x, y) * texelSize);
            samples++;
        }
    }
    
    fragColor = max(color / float(samples), ambient_light_level);
}
