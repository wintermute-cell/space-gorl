#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

#define size 16.0
#define OCTAVES 6
#define seed 6.0
#define pixels 180.0
#define should_tile false
#define reduce_background false
#define uv_correct vec2(1.0)

uniform sampler2D colorscheme;
uniform float time;

float rand(vec3 coord) {
    // Adapted for 3D
    return fract(sin(dot(coord, vec3(12.9898, 78.233, 271.9))) * 43758.5453);
}

float noise(vec3 coord){
    // 3D Perlin noise
    vec3 i = floor(coord);
    vec3 f = fract(coord);
    vec3 u = f * f * (3.0 - 2.0 * f);

    // Hashing the corner coordinates
    float n000 = rand(i);
    float n100 = rand(i + vec3(1.0, 0.0, 0.0));
    float n010 = rand(i + vec3(0.0, 1.0, 0.0));
    float n001 = rand(i + vec3(0.0, 0.0, 1.0));
    float n101 = rand(i + vec3(1.0, 0.0, 1.0));
    float n011 = rand(i + vec3(0.0, 1.0, 1.0));
    float n110 = rand(i + vec3(1.0, 1.0, 0.0));
    float n111 = rand(i + vec3(1.0, 1.0, 1.0));

    // Interpolate between the corner values
    float nx00 = mix(n000, n100, u.x);
    float nx01 = mix(n001, n101, u.x);
    float nx10 = mix(n010, n110, u.x);
    float nx11 = mix(n011, n111, u.x);

    float nxy0 = mix(nx00, nx10, u.y);
    float nxy1 = mix(nx01, nx11, u.y);

    return mix(nxy0, nxy1, u.z);
}

float fbm(vec3 coord){
    float value = 0.0;
    float amplitude = 0.5;

    for(int i = 0; i < OCTAVES; i++){
        value += amplitude * noise(coord);
        coord *= 2.0;
        amplitude *= 0.5;
    }

    return value;
}

bool dither(vec2 uv1, vec2 uv2) {
	return mod(uv1.y+uv2.x,2.0/pixels) <= 1.0 / pixels;
}

float circleNoise(vec3 uv) {
    if (should_tile) {
        uv.xy = mod(uv.xy, uv.z);
    }

    float uv_y = floor(uv.y);
    uv.x += uv_y * 0.31;
    vec2 f = fract(uv.xy);
    float h = rand(vec3(floor(uv.x), floor(uv_y), 0));
    float m = (length(f - 0.25 - (h * 0.5)));
    float r = h * 0.25;
    return smoothstep(0.0, r, m * 0.75);
}

float cloud_alpha(vec3 uv) {
    float c_noise = 0.0;

    int iters = 2;
    for (int i = 0; i < iters; i++) {
        c_noise += circleNoise(uv * 0.5 + vec3(float(i + 1), -0.3, 0.0));
    }

    return fbm(uv + vec3(c_noise, 0.0, 0.0));
}


vec2 rotate(vec2 vec, float angle) {
	vec -=vec2(0.5);
	vec *= mat2(vec2(cos(angle),-sin(angle)), vec2(sin(angle),cos(angle)));
	vec += vec2(0.5);
	return vec;
}

void main() {
    vec3 uv3 = vec3(fragTexCoord * size, time/10);
	bool dith = false;//dither(uv3, fragTexCoord);

    // Updated function calls with 3D coordinates
    float n_alpha = fbm(uv3 + vec3(2, 2, 0));
    float n_dust = cloud_alpha(uv3);
    float n_dust2 = fbm(uv3 * vec3(0.2) - vec3(2, 2, 0));
    float n_dust_lerp = n_dust2 * n_dust;	// apply dithering
	if (dith) {
		n_dust_lerp *= 0.95;
	}

	// choose alpha value
	float a_dust = step(n_alpha , n_dust_lerp * 1.8);
	n_dust_lerp = pow(n_dust_lerp, 3.2) * 56.0;
	if (dith) {
		n_dust_lerp *= 1.1;
	}
	
	// choose & apply colors
	if (reduce_background) {
		n_dust_lerp = pow(n_dust_lerp, 0.8) * 0.7;
	}
	
	float col_value = floor(n_dust_lerp) / 7.0;
	vec3 col = texture(colorscheme, vec2(col_value, 0.0)).rgb;
	
	
	fragColor = vec4(col, a_dust);
}
