#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
#define modulate 1.0
//uniform float modulate;

void main( void )
{
    vec4 a = texture2D( texture0, fragTexCoord ) * vec4( modulate );
    vec4 b = texture2D( texture1, fragTexCoord );

    fragColor = max( a, b * 0.32 );
    //fragColor = mix( a, b, 0.2 );
}   
