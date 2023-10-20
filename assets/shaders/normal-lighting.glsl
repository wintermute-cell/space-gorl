#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;
uniform vec2 light_position; // the lights position in worlds space
uniform vec2 self_position; // the sprites position in world space
uniform vec3 light_color;
uniform float rotation_angle; // in radians
uniform vec2 self_origin; // the sprites origin (in worldspace? whatever raylib uses)
uniform vec3 light_range;

#define light_height 0.66

mat2 rotationMatrix(float angle) {
    return mat2(cos(angle), -sin(angle), sin(angle), cos(angle));
}

vec2 rotatePointAroundPivot(vec2 point, float angle, vec2 pivot) {
    // Translate the point to the origin (relative to pivot)
    vec2 translatedPoint = point - pivot;
    
    float radAngle = angle;
    float cosA = cos(radAngle);
    float sinA = sin(radAngle);
    
    // Perform the rotation around the origin
    float x = translatedPoint.x * cosA - translatedPoint.y * sinA;
    float y = translatedPoint.x * sinA + translatedPoint.y * cosA;
    
    // Translate the rotated point back to its original position
    return vec2(x, y) + pivot;
}

void main(void) {
    vec4 normal_map = texture2D(texture0, fragTexCoord);
    vec3 normal_vec = normal_map.rgb * 2.0 - 1.0;

    vec2 rotated_normal = rotationMatrix(-rotation_angle) * normal_vec.xy;
    vec3 rotated_normal_vec = vec3(rotated_normal, normal_vec.z);

    vec2 origin = self_origin / textureSize(texture0, 0); // calculate 0 .. 1 origin in sprite space
    vec2 adjustedFragTexCoord = fragTexCoord - origin; // offset the UVs by the origin
    vec2 pixelWorldPosition = self_position + (adjustedFragTexCoord) * vec2(textureSize(texture0, 0));
    pixelWorldPosition = rotatePointAroundPivot(pixelWorldPosition, rotation_angle, self_position);

    vec3 light_dir = vec3(pixelWorldPosition - light_position, light_height);

    vec3 N = normalize(rotated_normal_vec);
    vec3 L = normalize(light_dir);

    float distanceToLight = length(light_dir.xy);
    // Quadratic attenuation
    float influenceX = smoothstep(light_range.x, light_range.x - light_range.z, abs(light_dir.x));
    float influenceY = smoothstep(light_range.y, light_range.y - light_range.z, abs(light_dir.y));
    float attenuation = influenceX * influenceY;

    float light = max(dot(N, L), 0.0) * attenuation;
    fragColor = vec4(vec3(light) * light_color, normal_map.a);
}
