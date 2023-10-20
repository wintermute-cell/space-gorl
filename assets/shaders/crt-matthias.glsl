#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

//uniform vec2 resolution;
#define CHROMATIC_ABBERATION 1.1
#define RESOLUTION vec2(640, 480)
//uniform float time;
//#define TIME 0.0 // i guess if time is always zero we just miss out on animation?
uniform float TIME;
uniform sampler2D texture0;
uniform sampler2D texture1;

// these are never used??
uniform float cfg_curvature;
uniform float cfg_scanlines;
uniform float cfg_shadow_mask; 
uniform float cfg_separation; 
uniform float cfg_ghosting; 
uniform float cfg_noise; 
uniform float cfg_flicker;
uniform float cfg_vignette; 
uniform float cfg_distortion; 
uniform float cfg_aspect_lock;
uniform float cfg_hpos;
uniform float cfg_vpos; 
uniform float cfg_hsize; 
uniform float cfg_vsize; 
uniform float cfg_contrast; 
uniform float cfg_brightness; 
uniform float cfg_saturation;
uniform float cfg_degauss;

vec3 tsample( sampler2D samp, vec2 tc, float offs, vec2 resolution )
{
    vec3 s = pow( abs( texture2D( samp, vec2( tc.x, tc.y ) ).rgb), vec3( 2.2 ) );
    return s*vec3(1.25);
}

vec3 filmic( vec3 LinearColor )
{
    vec3 x = max( vec3(0.0), LinearColor-vec3(0.004));
    return (x*(6.2*x+0.5))/(x*(6.2*x+1.7)+0.06);
}

vec2 curve( vec2 uv )
{
    // This curvature here will cause tear artifacts when the shader runs on a
    // texture with a low resolution.
    uv = (uv - 0.5) * 2.0;
    uv *= 1.1;  
    // 5.0 and 4.0 control the curve amount
    uv.x *= 1.0 + pow((abs(uv.y) / 5.0), 2.0);
    uv.y *= 1.0 + pow((abs(uv.x) / 4.0), 2.0);
    uv  = (uv / 2.0) + 0.5;
    uv =  uv *0.92 + 0.04;
    return uv;
}

float rand(vec2 co)
    {
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
    }
    
void main(void){
        /* Curve */
        vec2 curved_uv = mix( curve( fragTexCoord ), fragTexCoord, 0.8 ); // mix value adjusts curve amount. higher value, less curve

        // this appears to zoom the screen
        float scale = 0.0;
        vec2 scuv = curved_uv*(1.0-scale)+scale/2.0+vec2(0.003, -0.001);

        /* Main color, Bleed */
        vec3 col;

        // this adds chromatic abberation that is animated if TIME is changing
        float x =  sin(0.1*TIME+curved_uv.y*13.0)*sin(0.23*TIME+curved_uv.y*19.0)*sin(0.3+0.11*TIME+curved_uv.y*23.0)*0.0012;
        float o = sin(gl_FragCoord.y*1.5)/RESOLUTION.x;
        x+=o*0.25;

        float brightener = 0.01;
        col.r = tsample(
            texture0,vec2(
                x+scuv.x+0.0005*CHROMATIC_ABBERATION,scuv.y+0.0005*CHROMATIC_ABBERATION),RESOLUTION.y/800.0, RESOLUTION ).x+brightener;
        col.g = tsample(                                                                                                  
            texture0,vec2(                                                                                               
                x+scuv.x+0.0000*CHROMATIC_ABBERATION,scuv.y-0.0007*CHROMATIC_ABBERATION),RESOLUTION.y/800.0, RESOLUTION ).y+brightener;
        col.b = tsample(                                                                                                
            texture0,vec2(                                                                                             
                x+scuv.x-0.0011*CHROMATIC_ABBERATION,scuv.y+0.0000*CHROMATIC_ABBERATION),RESOLUTION.y/800.0, RESOLUTION ).z+brightener+0.002;

        /* Ghosting */

        float i = clamp(col.r*0.299 + col.g*0.587 + col.b*0.114, 0.0, 1.0 );        
        i = pow( 1.0 - pow(i,2.0), 1.0 );
        i = (1.0-i) * 0.85 + 0.15;  

        float ghs = 0.13;
        vec3 r = tsample(texture1, vec2(x-0.014*1.0, -0.027)*0.85+0.007*vec2( 0.35*sin(1.0/7.0 + 15.0*curved_uv.y + 0.9*TIME), 
            0.35*sin( 2.0/7.0 + 10.0*curved_uv.y + 1.37*TIME) )+vec2(scuv.x+0.001,scuv.y+0.001),
            5.5+1.3*sin( 3.0/9.0 + 31.0*curved_uv.x + 1.70*TIME),RESOLUTION).xyz*vec3(0.5,0.25,0.25);
        vec3 g = tsample(texture1, vec2(x-0.019*1.0, -0.020)*0.85+0.007*vec2( 0.35*cos(1.0/9.0 + 15.0*curved_uv.y + 0.5*TIME), 
            0.35*sin( 2.0/9.0 + 10.0*curved_uv.y + 1.50*TIME) )+vec2(scuv.x+0.000,scuv.y-0.002),
            5.4+1.3*sin( 3.0/3.0 + 71.0*curved_uv.x + 1.90*TIME),RESOLUTION).xyz*vec3(0.25,0.5,0.25);
        vec3 b = tsample(texture1, vec2(x-0.017*1.0, -0.003)*0.85+0.007*vec2( 0.35*sin(2.0/3.0 + 15.0*curved_uv.y + 0.7*TIME), 
            0.35*cos( 2.0/3.0 + 10.0*curved_uv.y + 1.63*TIME) )+vec2(scuv.x-0.002,scuv.y+0.000),
            5.3+1.3*sin( 3.0/7.0 + 91.0*curved_uv.x + 1.65*TIME),RESOLUTION).xyz*vec3(0.25,0.25,0.5);
        
        col += vec3(ghs*(1.0-0.299))*pow(clamp(vec3(3.0)*r,vec3(0.0),vec3(1.0)),vec3(2.0))*vec3(i);
        col += vec3(ghs*(1.0-0.587))*pow(clamp(vec3(3.0)*g,vec3(0.0),vec3(1.0)),vec3(2.0))*vec3(i);
        col += vec3(ghs*(1.0-0.114))*pow(clamp(vec3(3.0)*b,vec3(0.0),vec3(1.0)),vec3(2.0))*vec3(i);


        // TODO: add adjustment parameters to all these effects
        /* Level adjustment (curves) */
        col *= vec3(0.75,0.79,0.75);
        col = clamp(col*1.3 + 0.75*col*col + 1.25*col*col*col*col*col,vec3(0.0),vec3(10.0));

        /* Vignette */
        float vig = (0.1 + 1.0*16.0*curved_uv.x*curved_uv.y*(1.0-curved_uv.x)*(1.0-curved_uv.y));
        vig = 1.3*pow(vig,0.5);
        col *= vig;

        /* Scanlines */
        float scans = clamp( 0.35+0.18*sin(6.0*TIME+curved_uv.y*RESOLUTION.y*1.5), 0.0, 1.0);
        float s = pow(scans,0.07);
        col = col * vec3(s);

        /* Vertical lines (shadow mask) */
        col*=1.0-0.10*(clamp((mod(gl_FragCoord.xy.x, 3.0))/2.0,0.0,1.0));

        /* Tone map */
        col = filmic( col );

        /* Noise */
        /*vec2 seed = floor(curved_uv*RESOLUTION.xy*vec2(0.5))/RESOLUTION.xy;*/
        vec2 seed = curved_uv*RESOLUTION.xy;;
        /* seed = curved_uv; */
        col -= 0.015*pow(vec3(rand( seed +TIME ), rand( seed +TIME*2.0 ), rand( seed +TIME * 3.0 ) ), vec3(1.5) );

        /* Flicker */
        col *= (1.0-0.005*(sin(50.0*TIME+curved_uv.y*2.0)*0.5+0.5));

        /* Clamp */
        if (curved_uv.x < 0.0 || curved_uv.x > 1.0)
            col *= 0.0;
        if (curved_uv.y < 0.0 || curved_uv.y > 1.0)
            col *= 0.0;

        fragColor = vec4( col, 1.0 );
    }
