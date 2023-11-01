#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;


vec3 blendSoftLight(vec3 base, vec3 blend) {
    return mix(
        sqrt(base) * (2.0 * blend - 1.0) + 2.0 * base * (1.0 - blend), 
        2.0 * base * blend + base * base * (1.0 - 2.0 * blend), 
        step(base, vec3(0.5))
    );
}

void main( void )
{
    vec4 a = texture2D( texture1, fragTexCoord );
    vec4 b = texture2D( texture0, fragTexCoord );
    fragColor = vec4(blendSoftLight(a.rgb, b.rgb), a.a);
}   
