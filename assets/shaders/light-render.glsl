#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;
uniform vec2 resolution;
uniform vec3 light_color;

#define PI 3.14
#define numBlurSamples 7

//sample from the 1D distance map
float sample(vec2 coord, float r) {
    return step(r, texture2D(texture0, coord).r);
}

void main(void) {
    //rectangular to polar
    vec2 norm = fragTexCoord.st * 2.0 - 1.0;
    float theta = atan(norm.y, norm.x);
    float r = length(norm);	
    float coord = (theta + PI) / (2.0*PI);

    //the tex coord to sample our 1D lookup texture	
    //always 0.0 on y axis
    vec2 tc = vec2(coord, 0.0);

    //the center tex coord, which gives us hard shadows
    float center = sample(tc, r);        

    //we multiply the blur amount by our distance from center
    //this leads to more blurriness as the shadow "fades away"
    float smoothness = 2.5;
    float blur = (smoothness/resolution.x)  * smoothstep(0., 1., r); 

    //now we use a simple gaussian blur
    float sum = 0.0;

    sum += sample(vec2(tc.x - 4.0*blur, tc.y), r) * 0.03;
    sum += sample(vec2(tc.x - 3.5*blur, tc.y), r) * 0.04;
    sum += sample(vec2(tc.x - 3.0*blur, tc.y), r) * 0.05;
    sum += sample(vec2(tc.x - 2.5*blur, tc.y), r) * 0.07;
    sum += sample(vec2(tc.x - 2.0*blur, tc.y), r) * 0.09;
    sum += sample(vec2(tc.x - 1.5*blur, tc.y), r) * 0.11;
    sum += sample(vec2(tc.x - 1.0*blur, tc.y), r) * 0.13;

    sum += center * 0.16;

    sum += sample(vec2(tc.x + 1.0*blur, tc.y), r) * 0.13; // symmetric weights for the other side
    sum += sample(vec2(tc.x + 1.5*blur, tc.y), r) * 0.11; 
    sum += sample(vec2(tc.x + 2.0*blur, tc.y), r) * 0.09;
    sum += sample(vec2(tc.x + 2.5*blur, tc.y), r) * 0.07; 
    sum += sample(vec2(tc.x + 3.0*blur, tc.y), r) * 0.05;
    sum += sample(vec2(tc.x + 3.5*blur, tc.y), r) * 0.04;
    sum += sample(vec2(tc.x + 4.0*blur, tc.y), r) * 0.03; 


    //sum of 1.0 -> in light, 0.0 -> in shadow

    //multiply the summed amount by our distance, which gives us a radial falloff
    //then multiply by vertex (light) color  
    vec3 positive = vec3(sum * smoothstep(1.0, 0.0, r));
    fragColor = vec4(positive * light_color, 1.0);
}
