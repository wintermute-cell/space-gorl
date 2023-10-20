// ---------------------------------
// CRT Matthias Post Processing Shader
// ---------------------------------
// 
// Works with the default Raylib vertex shader.
// Worse performance than crt-lottes. Works with dedicated (but weak) GPU,
// but tanks FPS with integrated GPU module.
//

#version 330

#ifdef GL_ES
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif
#define COMPAT_PRECISION mediump
#else
#define COMPAT_PRECISION
#endif

uniform COMPAT_PRECISION int FrameDirection;
uniform COMPAT_PRECISION int FrameCount;
uniform COMPAT_PRECISION vec2 OutputSize;
uniform COMPAT_PRECISION vec2 TextureSize;
uniform COMPAT_PRECISION vec2 InputSize;

#define OutputSize  vec2(1280, 960)
#define TextureSize vec2(640, 480) // Doesnt affect RGB pixel effect but affects sharpness
#define InputSize   vec2(4000, 4000)

uniform sampler2D texture0;


// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

// fragment compatibility #defines
#define Source texture0
#define vTexCoord TEX0.xy

#define SourceSize vec4(TextureSize, 1.0 / TextureSize) //either TextureSize or InputSize
#define outsize vec4(OutputSize, 1.0 / OutputSize)

#define hardScan -8.0
#define hardPix -3.0
#define warpX 0.0
#define warpY 0.0
#define maskDark 0.5
#define maskLight 1.5
#define scaleInLinearGamma 1.0
#define shadowMask 3.0
#define brightBoost 1.0
#define hardBloomPix -1.5
#define hardBloomScan -2.0
#define bloomAmount 0.15
#define shape 3.0

#define DO_BLOOM

// ------------- //

// sRGB to Linear.
// Assuming using sRGB typed textures this should not be needed.
float ToLinear1(float c)
{
    if (scaleInLinearGamma == 0.) 
        return c;
    
    return(c<=0.04045) ? c/12.92 : pow((c + 0.055)/1.055, 2.4);
}

vec3 ToLinear(vec3 c)
{
    if (scaleInLinearGamma==0.) 
        return c;
    
    return vec3(ToLinear1(c.r), ToLinear1(c.g), ToLinear1(c.b));
}

// Linear to sRGB.
// Assuming using sRGB typed textures this should not be needed.
float ToSrgb1(float c)
{
    if (scaleInLinearGamma == 0.) 
        return c;
    
    return(c<0.0031308 ? c*12.92 : 1.055*pow(c, 0.41666) - 0.055);
}

vec3 ToSrgb(vec3 c)
{
    if (scaleInLinearGamma == 0.) 
        return c;
    
    return vec3(ToSrgb1(c.r), ToSrgb1(c.g), ToSrgb1(c.b));
}

// Nearest emulated sample given floating point position and texel offset.
// Also zero's off screen.
vec3 Fetch(vec2 pos,vec2 off){
  pos=(floor(pos*SourceSize.xy+off)+vec2(0.5,0.5))/SourceSize.xy;
#ifdef SIMPLE_LINEAR_GAMMA
  return ToLinear(brightBoost * pow(texture(Source,pos.xy).rgb, vec3(2.2)));
#else
  return ToLinear(brightBoost * texture(Source,pos.xy).rgb);
#endif
}

// Distance in emulated pixels to nearest texel.
vec2 Dist(vec2 pos)
{
    pos = pos*SourceSize.xy;
    
    return -((pos - floor(pos)) - vec2(0.5));
}
    
// 1D Gaussian.
float Gaus(float pos, float scale)
{
    return exp2(scale*pow(abs(pos), shape));
}

// 3-tap Gaussian filter along horz line.
vec3 Horz3(vec2 pos, float off)
{
    vec3 b    = Fetch(pos, vec2(-1.0, off));
    vec3 c    = Fetch(pos, vec2( 0.0, off));
    vec3 d    = Fetch(pos, vec2( 1.0, off));
    float dst = Dist(pos).x;

    // Convert distance to weight.
    float scale = hardPix;
    float wb = Gaus(dst-1.0,scale);
    float wc = Gaus(dst+0.0,scale);
    float wd = Gaus(dst+1.0,scale);

    // Return filtered sample.
    return (b*wb+c*wc+d*wd)/(wb+wc+wd);
}

// 5-tap Gaussian filter along horz line.
vec3 Horz5(vec2 pos,float off){
    vec3 a = Fetch(pos,vec2(-2.0, off));
    vec3 b = Fetch(pos,vec2(-1.0, off));
    vec3 c = Fetch(pos,vec2( 0.0, off));
    vec3 d = Fetch(pos,vec2( 1.0, off));
    vec3 e = Fetch(pos,vec2( 2.0, off));
    
    float dst = Dist(pos).x;
    // Convert distance to weight.
    float scale = hardPix;
    float wa = Gaus(dst - 2.0, scale);
    float wb = Gaus(dst - 1.0, scale);
    float wc = Gaus(dst + 0.0, scale);
    float wd = Gaus(dst + 1.0, scale);
    float we = Gaus(dst + 2.0, scale);
    
    // Return filtered sample.
    return (a*wa+b*wb+c*wc+d*wd+e*we)/(wa+wb+wc+wd+we);
}
  
// 7-tap Gaussian filter along horz line.
vec3 Horz7(vec2 pos,float off)
{
    vec3 a = Fetch(pos, vec2(-3.0, off));
    vec3 b = Fetch(pos, vec2(-2.0, off));
    vec3 c = Fetch(pos, vec2(-1.0, off));
    vec3 d = Fetch(pos, vec2( 0.0, off));
    vec3 e = Fetch(pos, vec2( 1.0, off));
    vec3 f = Fetch(pos, vec2( 2.0, off));
    vec3 g = Fetch(pos, vec2( 3.0, off));

    float dst = Dist(pos).x;
    // Convert distance to weight.
    float scale = hardBloomPix;
    float wa = Gaus(dst - 3.0, scale);
    float wb = Gaus(dst - 2.0, scale);
    float wc = Gaus(dst - 1.0, scale);
    float wd = Gaus(dst + 0.0, scale);
    float we = Gaus(dst + 1.0, scale);
    float wf = Gaus(dst + 2.0, scale);
    float wg = Gaus(dst + 3.0, scale);

    // Return filtered sample.
    return (a*wa+b*wb+c*wc+d*wd+e*we+f*wf+g*wg)/(wa+wb+wc+wd+we+wf+wg);
}
  
// Return scanline weight.
float Scan(vec2 pos, float off)
{
    float dst = Dist(pos).y;

    return Gaus(dst + off, hardScan);
}
  
// Return scanline weight for bloom.
float BloomScan(vec2 pos, float off)
{
    float dst = Dist(pos).y;
    
    return Gaus(dst + off, hardBloomScan);
}

// Allow nearest three lines to effect pixel.
vec3 Tri(vec2 pos)
{
    vec3 a = Horz3(pos,-1.0);
    vec3 b = Horz5(pos, 0.0);
    vec3 c = Horz3(pos, 1.0);
    
    float wa = Scan(pos,-1.0); 
    float wb = Scan(pos, 0.0);
    float wc = Scan(pos, 1.0);
    
    return a*wa + b*wb + c*wc;
}
  
// Small bloom.
vec3 Bloom(vec2 pos)
{
    vec3 a = Horz5(pos,-2.0);
    vec3 b = Horz7(pos,-1.0);
    vec3 c = Horz7(pos, 0.0);
    vec3 d = Horz7(pos, 1.0);
    vec3 e = Horz5(pos, 2.0);

    float wa = BloomScan(pos,-2.0);
    float wb = BloomScan(pos,-1.0); 
    float wc = BloomScan(pos, 0.0);
    float wd = BloomScan(pos, 1.0);
    float we = BloomScan(pos, 2.0);

    return a*wa+b*wb+c*wc+d*wd+e*we;
}
  
// Distortion of scanlines, and end of screen alpha.
vec2 Warp(vec2 pos)
{
    pos  = pos*2.0-1.0;    
    pos *= vec2(1.0 + (pos.y*pos.y)*warpX, 1.0 + (pos.x*pos.x)*warpY);
    
    return pos*0.5 + 0.5;
}
  
// Shadow mask.
vec3 Mask(vec2 pos)
{
    vec3 mask = vec3(maskDark, maskDark, maskDark);
  
    // Very compressed TV style shadow mask.
    if (shadowMask == 1.0) 
    {
        float line = maskLight;
        float odd = 0.0;
        
        if (fract(pos.x*0.166666666) < 0.5) odd = 1.0;
        if (fract((pos.y + odd) * 0.5) < 0.5) line = maskDark;  
        
        pos.x = fract(pos.x*0.333333333);

        if      (pos.x < 0.333) mask.r = maskLight;
        else if (pos.x < 0.666) mask.g = maskLight;
        else                    mask.b = maskLight;
        mask*=line;  
    } 

    // Aperture-grille.
    else if (shadowMask == 2.0) 
    {
        pos.x = fract(pos.x*0.333333333);

        if      (pos.x < 0.333) mask.r = maskLight;
        else if (pos.x < 0.666) mask.g = maskLight;
        else                    mask.b = maskLight;
    } 

    // Stretched VGA style shadow mask (same as prior shaders).
    else if (shadowMask == 3.0) 
    {
        float x_fac = 6.0;   // controls the horizontal repeat rate
        float y_fac = 3.0;   // controls the vertical repeat rate, you can adjust this value

        pos.x += (pos.y * y_fac) * x_fac;
        pos.x = fract(pos.x * (1 / (x_fac * 2.0)));

        if      (pos.x < 0.333) mask.r = maskLight;
        else if (pos.x < 0.666) mask.g = maskLight;
        else                    mask.b = maskLight;
    }

    // VGA style shadow mask.
    else if (shadowMask == 4.0) 
    {
        pos.xy  = floor(pos.xy*vec2(1.0, 0.5));
        pos.x  += pos.y*3.0;
        pos.x   = fract(pos.x*0.166666666);

        if      (pos.x < 0.333) mask.r = maskLight;
        else if (pos.x < 0.666) mask.g = maskLight;
        else                    mask.b = maskLight;
    }

    return mask;
}

void main()
{
    vec2 pos = Warp(fragTexCoord.xy*(TextureSize.xy/InputSize.xy))*(InputSize.xy/TextureSize.xy);
    vec3 outColor = Tri(pos);

#ifdef DO_BLOOM
    //Add Bloom
    outColor.rgb += Bloom(pos)*bloomAmount;
#endif

    if (shadowMask > 0.0)
        outColor.rgb *= Mask(gl_FragCoord.xy * 1.000001);
    
#ifdef GL_ES    /* TODO/FIXME - hacky clamp fix */
    vec2 bordertest = (pos);
    if ( bordertest.x > 0.0001 && bordertest.x < 0.9999 && bordertest.y > 0.0001 && bordertest.y < 0.9999)
        outColor.rgb = outColor.rgb;
    else
        outColor.rgb = vec3(0.0);
#endif
    finalColor = vec4(ToSrgb(outColor.rgb), 1.0);
} 
