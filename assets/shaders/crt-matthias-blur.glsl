#version 330

out vec4 fragColor;

uniform vec2 BLUR;

uniform sampler2D texture0;
in vec2 fragTexCoord;

void main( void )
{
    vec4 sum = texture2D( texture0, fragTexCoord ) * 0.2270270270;
    sum += texture2D(texture0, vec2( fragTexCoord.x - 4.0 * BLUR.x, fragTexCoord.y - 4.0 * BLUR.y ) ) * 0.0162162162;
    sum += texture2D(texture0, vec2( fragTexCoord.x - 3.0 * BLUR.x, fragTexCoord.y - 3.0 * BLUR.y ) ) * 0.0540540541;
    sum += texture2D(texture0, vec2( fragTexCoord.x - 2.0 * BLUR.x, fragTexCoord.y - 2.0 * BLUR.y ) ) * 0.1216216216;
    sum += texture2D(texture0, vec2( fragTexCoord.x - 1.0 * BLUR.x, fragTexCoord.y - 1.0 * BLUR.y ) ) * 0.1945945946;
    sum += texture2D(texture0, vec2( fragTexCoord.x + 1.0 * BLUR.x, fragTexCoord.y + 1.0 * BLUR.y ) ) * 0.1945945946;
    sum += texture2D(texture0, vec2( fragTexCoord.x + 2.0 * BLUR.x, fragTexCoord.y + 2.0 * BLUR.y ) ) * 0.1216216216;
    sum += texture2D(texture0, vec2( fragTexCoord.x + 3.0 * BLUR.x, fragTexCoord.y + 3.0 * BLUR.y ) ) * 0.0540540541;
    sum += texture2D(texture0, vec2( fragTexCoord.x + 4.0 * BLUR.x, fragTexCoord.y + 4.0 * BLUR.y ) ) * 0.0162162162;
    fragColor = sum;
}   
