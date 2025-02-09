/*
-----------------------------------------------------------
| Multi-Stage Custom Renderer for                         |
-----------------------------------------------------------

Stage 1: Initial Render (TargetTex)
  - Draws sprites at a fixed resolution.
  - Uses Point Sampling for pixel-perfect scaling.

Stage 2: Intermediary Upscale (CeilTex)
  - Scales up TargetTex using integer scaling.
  - Purpose: Scale as close to output resolution while maintaining sharp edges.
  - Uses Point Sampling for crisp edges.

Stage 3: Final Render (To Screen)
  - Scales down CeilTex to fit screen size.
  - Purpose: Mask subpixel inconsistencies.
  - Uses Bilinear Sampling for smoother edges.

*/

package render

import (
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/settings"
	"cowboy-gorl/pkg/util"
	"math"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"
)


type RenderStage struct {
    Target *rl.RenderTexture2D
    WSCamera *rl.Camera2D
    SSCamera *rl.Camera2D
    trueWSCamera rl.Camera2D
    trueSSCamera rl.Camera2D
    Shader rl.Shader
    ClearColor rl.Color
    ClearBefore bool
}

func (rs *RenderStage) Begin() {
    Rs.CurrentStage = rs
    setCamTargetInternal(Rs.current_ws_cam_pos)
    rl.BeginTextureMode(*rs.Target)
    if rs.WSCamera != nil {
        rl.BeginMode2D(rs.trueWSCamera)
    }
}

func (rs *RenderStage) Prepare() {
    if rs.ClearBefore {
        rl.BeginTextureMode(*rs.Target)
        rl.ClearBackground(rs.ClearColor)
        rl.EndTextureMode()

        if rs.WSCamera != nil {
            rs.trueWSCamera = rl.NewCamera2D(rs.WSCamera.Offset, rs.WSCamera.Target, rs.WSCamera.Rotation, rs.WSCamera.Zoom)
        }
        if rs.SSCamera != nil {
            traumaSq := (float64(Rs.cameraTrauma) * float64(Rs.cameraTrauma))
            rotOffs := Rs.shakePerlin.Noise1D(
                rl.GetTime()+333*8) * (traumaSq) * float64(Rs.RenderResolution.X)

            offsetX := Rs.shakePerlin.Noise1D(
                rl.GetTime()      *8) * traumaSq * float64(Rs.RenderResolution.X) * 3
            offsetY := Rs.shakePerlin.Noise1D(
                (rl.GetTime()+666)*8) * traumaSq * float64(Rs.RenderResolution.Y) * 3
            posOffs := rl.Vector2Add(
                rs.SSCamera.Offset,
                rl.NewVector2(float32(offsetX), float32(offsetY)))

            m := rl.NewVector2(float32(rl.GetScreenWidth()) * 0.5, float32(rl.GetScreenHeight()) * 0.5)
            posOffs = rl.Vector2Add(posOffs, m)
            rs.trueSSCamera = rl.NewCamera2D(posOffs, m, float32(rotOffs), rs.SSCamera.Zoom)
        }
    }
}

func (rs *RenderStage) Continue() {
    Rs.CurrentStage = rs
    rl.BeginTextureMode(*rs.Target)
    if rs.WSCamera != nil {
        rl.BeginMode2D(rs.trueWSCamera)
    }
}

func (rs *RenderStage) End() {
	if rs.WSCamera != nil {
		rl.EndMode2D()
	}
    rl.EndTextureMode()
}

func GetOrCreateWSCamera(target *rl.RenderTexture2D) *rl.Camera2D {
	if cam, exists := Rs.ws_cameras[target]; exists {
		return cam
	}

    // the zoom helps us to scale everyting to the same relative size.
    // otherwise, if we would render at a high resolution, objects would appear
    // tiny, and would be wrongly positioned
    zoom := float32(target.Texture.Width) / float32(Rs.RenderResolution.X) 

	newCam := &rl.Camera2D{Zoom: zoom, Offset: rl.NewVector2(float32(target.Texture.Width)/2, float32(target.Texture.Height)/2)}
	Rs.ws_cameras[target] = newCam
	return newCam
}

func GetOrCreateSSCamera(target *rl.RenderTexture2D) *rl.Camera2D {
	if cam, exists := Rs.ss_cameras[target]; exists {
		return cam
	}

	newCam := &rl.Camera2D{Zoom: 1.0}
	Rs.ss_cameras[target] = newCam
	return newCam
}

func RegisterRenderTarget(name string, width int32, height int32, filter rl.TextureFilterMode) *rl.RenderTexture2D {
	rt := rl.LoadRenderTexture(width, height)
	rl.SetTextureFilter(rt.Texture, filter)
	Rs.targets[name] = &rt
	return &rt
}

func GetRenderTarget(name string) *rl.RenderTexture2D {
	return Rs.targets[name]
}



type crt_shader_data struct {
	Blendbuffer         rl.RenderTexture2D
	Accumulationbuffer rl.RenderTexture2D
	Blurbuffer         rl.RenderTexture2D

	Crtshader           rl.Shader
	Crtshader_loc_time  int32
	Blurshader          rl.Shader
	Blurshader_loc_blur int32
	Accumulateshader    rl.Shader
	Blendshader         rl.Shader
}

type RenderState struct {
    PrimaryStage *RenderStage
    FxStage *RenderStage
	GuiTargetTex       rl.RenderTexture2D
	CeilTex            rl.RenderTexture2D
	RenderResolution   rl.Vector2
	RenderScale        rl.Vector2
	MinScale           float32
	MinScaleOrig       float32 // the original render->screen factor, before changing screen size with fullscreen or similar

    crt_dat crt_shader_data 

    CurrentStage *RenderStage

    ws_cameras map[*rl.RenderTexture2D]*rl.Camera2D
    ss_cameras map[*rl.RenderTexture2D]*rl.Camera2D

    targets map[string]*rl.RenderTexture2D

    current_render_target *rl.RenderTexture2D
	current_ws_cam_pos      rl.Vector2
	current_ws_cam_target   rl.Vector2

    cam_clamp_bounds rl.Rectangle

	tex1locs map[rl.Shader]int32

	lastScreenHeight int32

    cameraTrauma float32
    shakePerlin *perlin.Perlin
}

var Rs RenderState

func Init(render_width int, render_height int) {
	Rs = RenderState{
		RenderResolution: rl.NewVector2(
			float32(render_width),
			float32(render_height)),
        ws_cameras: make(map[*rl.RenderTexture2D]*rl.Camera2D),
        ss_cameras: make(map[*rl.RenderTexture2D]*rl.Camera2D),
        targets: make(map[string]*rl.RenderTexture2D),
        shakePerlin: perlin.NewPerlin(1, 1, 1, 0),
	}

	// this is used to detect if the screen is resized
	Rs.lastScreenHeight = int32(rl.GetScreenHeight())

	recalcScaleFactor()
	Rs.MinScaleOrig = Rs.MinScale

	// create the primary render texture. all sprites will be drawn directly
	// to this texture.
    RegisterRenderTarget("Primary", int32(Rs.RenderResolution.X), int32(Rs.RenderResolution.Y), rl.FilterPoint)
    Rs.PrimaryStage = &RenderStage{
        Target: GetRenderTarget("Primary"),
        WSCamera: GetOrCreateWSCamera(GetRenderTarget("Primary")),
        SSCamera: GetOrCreateSSCamera(GetRenderTarget("Primary")),
        ClearBefore: true,
        ClearColor: rl.Blank,
    }


	Rs.GuiTargetTex = rl.LoadRenderTexture(
		int32(Rs.RenderResolution.X),
		int32(Rs.RenderResolution.Y))
	rl.SetTextureFilter(Rs.GuiTargetTex.Texture, rl.FilterPoint)

	// create a secondary render texture, that is the next integer scaling step
	// larger than the window resolution.
	ceilX := int32(Rs.RenderResolution.X) * int32(math.Ceil(float64(Rs.MinScale)))
	ceilY := int32(Rs.RenderResolution.Y) * int32(math.Ceil(float64(Rs.MinScale)))
	Rs.CeilTex = rl.LoadRenderTexture(int32(ceilX), int32(ceilY))
	rl.SetTextureFilter(Rs.CeilTex.Texture, rl.FilterBilinear)

    RegisterRenderTarget("Fx", int32(Rs.RenderResolution.X*4), int32(Rs.RenderResolution.Y*4), rl.FilterBilinear)
    Rs.FxStage = &RenderStage{
        Target: GetRenderTarget("Fx"),
        WSCamera: GetOrCreateWSCamera(GetRenderTarget("Fx")),
        SSCamera: GetOrCreateSSCamera(GetRenderTarget("Fx")),
        ClearBefore: true,
        ClearColor: rl.Blank,
    }

	// create a number of helper buffers for performing multiple render passes
	// to achieve the accumulation effect of the crt shader
	Rs.crt_dat.Blendbuffer = rl.LoadRenderTexture(ceilX, ceilY)
	rl.SetTextureFilter(Rs.crt_dat.Blendbuffer.Texture, rl.FilterBilinear)

	Rs.crt_dat.Accumulationbuffer = rl.LoadRenderTexture(ceilX, ceilY)
	rl.SetTextureFilter(Rs.crt_dat.Accumulationbuffer.Texture, rl.FilterBilinear)

	Rs.crt_dat.Blurbuffer = rl.LoadRenderTexture(ceilX, ceilY)
	rl.SetTextureFilter(Rs.crt_dat.Blurbuffer.Texture, rl.FilterBilinear)

	// load the shaders and get their uniform locations if needed
	// crt shader
	Rs.crt_dat.Crtshader = rl.LoadShader("", "shaders/crt-matthias.glsl")
	if Rs.crt_dat.Crtshader.ID == 0 {
		logging.Error("Failed to load CRT shader!")
	}

	Rs.crt_dat.Crtshader_loc_time = rl.GetShaderLocation(Rs.crt_dat.Crtshader, "TIME")
	if Rs.crt_dat.Crtshader_loc_time == 0 {
		logging.Error("Failed to find shader uniform location for TIME in crt-matthias.")
	}

	// blur shader
	Rs.crt_dat.Blurshader = rl.LoadShader("", "shaders/crt-matthias-blur.glsl")
	if Rs.crt_dat.Blurshader.ID == 0 {
		logging.Error("Failed to load CRT blur shader!")
	}
	Rs.crt_dat.Blurshader_loc_blur = rl.GetShaderLocation(Rs.crt_dat.Blurshader, "BLUR")
	if Rs.crt_dat.Blurshader_loc_blur == 0 {
		logging.Error("Failed to find shader uniform location for BLUR in crt-matthias-blur.")
	}
	rl.SetShaderValue(Rs.crt_dat.Blurshader, Rs.crt_dat.Blurshader_loc_blur, []float32{0.0003, 0.0004}, rl.ShaderUniformVec2)

	// accumulate shader
	Rs.crt_dat.Accumulateshader = rl.LoadShader("", "shaders/crt-matthias-accumulate.glsl")
	if Rs.crt_dat.Accumulateshader.ID == 0 {
		logging.Error("Failed to load CRT accumulate shader!")
	}

	// blend shader
	Rs.crt_dat.Blendshader = rl.LoadShader("", "shaders/crt-matthias-blend.glsl")
	if Rs.crt_dat.Blendshader.ID == 0 {
		logging.Error("Failed to load CRT blend shader!")
	}

	// use a map to store the texture1 locations
	Rs.tex1locs = make(map[rl.Shader]int32)
	Rs.tex1locs[Rs.crt_dat.Accumulateshader] = rl.GetShaderLocation(Rs.crt_dat.Accumulateshader, "texture1")
	Rs.tex1locs[Rs.crt_dat.Blendshader] = rl.GetShaderLocation(Rs.crt_dat.Blendshader, "texture1")

    // begin the primary stage here already, so packages can reference the
    // CurrentStage in their Initialization (for example to set the camera
    // target.)
    Rs.PrimaryStage.Begin()

	logging.Info("Custom Rendering Environment initialized.")
}

func Deinit() {
	rl.UnloadRenderTexture(Rs.crt_dat.Blendbuffer) 
	rl.UnloadRenderTexture(Rs.crt_dat.Accumulationbuffer)
	rl.UnloadRenderTexture(Rs.crt_dat.Blurbuffer)

	rl.UnloadShader(Rs.crt_dat.Crtshader)
	rl.UnloadShader(Rs.crt_dat.Blurshader)
	rl.UnloadShader(Rs.crt_dat.Accumulateshader)
	rl.UnloadShader(Rs.crt_dat.Blendshader)
}

// Camera functions
func setCamTargetInternal(target rl.Vector2) {

	Rs.CurrentStage.SSCamera.Target = target

	Rs.CurrentStage.WSCamera.Target.X = float32(int32(target.X))
	Rs.CurrentStage.SSCamera.Target.X -= Rs.CurrentStage.WSCamera.Target.X
	Rs.CurrentStage.SSCamera.Target.X *= Rs.MinScaleOrig

	Rs.CurrentStage.WSCamera.Target.Y = float32(int32(target.Y))
	Rs.CurrentStage.SSCamera.Target.Y -= Rs.CurrentStage.WSCamera.Target.Y
	Rs.CurrentStage.SSCamera.Target.Y *= Rs.MinScaleOrig
	Rs.CurrentStage.SSCamera.Target.Y = Rs.CurrentStage.SSCamera.Target.Y * -1
}

func GetWSCameraTarget() rl.Vector2 {
	return Rs.CurrentStage.WSCamera.Target
}

func GetWSCameraOffset() rl.Vector2 {
	return Rs.CurrentStage.WSCamera.Offset
}

// Returns the the difference between the camera offset and the center of the
// screen.
func GetWSCameraCenterOffset() rl.Vector2 {
	offs_from_center := rl.Vector2Subtract(
		GetWSCameraOffset(),
		rl.Vector2Scale(Rs.RenderResolution, 0.5),
	)
	return offs_from_center
}

func GetSSCameraTarget() rl.Vector2 {
	return Rs.CurrentStage.SSCamera.Target
}

func GetSSCameraOffset() rl.Vector2 {
	return Rs.CurrentStage.SSCamera.Offset
}

func SetCameraClampBounds(bounds rl.Rectangle) {
    Rs.cam_clamp_bounds = bounds
}

func SetCameraTargetSmooth(target rl.Vector2) {
    currentRenderResolution := rl.NewVector2(float32(Rs.CurrentStage.Target.Texture.Width), float32(Rs.CurrentStage.Target.Texture.Height))
    
    if Rs.cam_clamp_bounds != (rl.Rectangle{}) {
        target = util.Vector2Clamp(
            target,
            // stay away half a screen from the upper left bounds
            rl.Vector2Add(rl.Vector2Scale(currentRenderResolution, 0.5), rl.NewVector2(Rs.cam_clamp_bounds.X, Rs.cam_clamp_bounds.Y)),
            // stay away half a screen from the lower right bounds
            rl.Vector2Subtract(rl.NewVector2(Rs.cam_clamp_bounds.Width, Rs.cam_clamp_bounds.Height), rl.Vector2Scale(currentRenderResolution, 0.5)),
        )
    }
    Rs.current_ws_cam_target = target
}

func SetCameraTarget(target rl.Vector2) {
    currentRenderResolution := rl.NewVector2(
        float32(Rs.CurrentStage.Target.Texture.Width),
        float32(Rs.CurrentStage.Target.Texture.Height))
    
    if Rs.cam_clamp_bounds != (rl.Rectangle{}) {
        target = util.Vector2Clamp(
            target,
            // stay away half a screen from the upper left bounds
            rl.Vector2Add(rl.Vector2Scale(currentRenderResolution, 0.5), rl.NewVector2(Rs.cam_clamp_bounds.X, Rs.cam_clamp_bounds.Y)),
            // stay away half a screen from the lower right bounds
            rl.Vector2Subtract(rl.NewVector2(Rs.cam_clamp_bounds.Width, Rs.cam_clamp_bounds.Height), rl.Vector2Scale(currentRenderResolution, 0.5)),
        )
    }
    Rs.current_ws_cam_target = target
    Rs.current_ws_cam_pos = target
}

func update_ws_cam_lerp() {
	lerpFactor := float32(0.05) * rl.GetFrameTime() * 294 // adjust this value as needed

	// Linear interpolation
	Rs.current_ws_cam_pos.X = Rs.current_ws_cam_pos.X + (Rs.current_ws_cam_target.X-Rs.current_ws_cam_pos.X)*lerpFactor
	Rs.current_ws_cam_pos.Y = Rs.current_ws_cam_pos.Y + (Rs.current_ws_cam_target.Y-Rs.current_ws_cam_pos.Y)*lerpFactor
}

// ScreenToWorldPoint converts a screen-space point to a world-space point
func ScreenToWorldPoint(point rl.Vector2) rl.Vector2 {
    // TODO, raylib has own ScreenToWorld, .. functions, try them out
	v := rl.Vector2Zero()
	v.X = (point.X-Rs.RenderResolution.X/2.0)/Rs.CurrentStage.WSCamera.Zoom + Rs.current_ws_cam_target.X
	v.Y = (point.Y-Rs.RenderResolution.Y/2.0)/Rs.CurrentStage.WSCamera.Zoom + Rs.current_ws_cam_target.Y
	return v
}

// WorldToScreenPoint converts a world-space point to a screen-space point
func WorldToScreenPoint(point rl.Vector2) rl.Vector2 {
	v := rl.Vector2Zero()
	v.X = (point.X-Rs.current_ws_cam_pos.X)*Rs.CurrentStage.WSCamera.Zoom + Rs.RenderResolution.X/2.0
	v.Y = (point.Y-Rs.current_ws_cam_pos.Y)*Rs.CurrentStage.WSCamera.Zoom + Rs.RenderResolution.Y/2.0
	return v
}

func BeginCustomRenderWorldspace() {

    update_ws_cam_lerp()
    Rs.PrimaryStage.Prepare()
    Rs.FxStage.Prepare()


	if Rs.lastScreenHeight != int32(rl.GetScreenHeight()) {
		recalcScaleFactor()
		Rs.lastScreenHeight = int32(rl.GetScreenHeight())
	}

	// adjust mouse coordinates so that they match with the render resolution,
	// not the screen size.
	// NOTE: We might be able to move the following into recalcScaleFactor():
	rl.SetMouseOffset(int(-(float32(rl.GetScreenWidth())-Rs.RenderResolution.X*Rs.MinScale)*0.5),
		int(-(float32(rl.GetScreenHeight())-Rs.RenderResolution.Y*Rs.MinScale)*0.5))
	rl.SetMouseScale(1/Rs.MinScale, 1/Rs.MinScale)

	// begin rendering to the primary render texture
    Rs.PrimaryStage.Begin()
}

func PauseTargetTex() {
    Rs.PrimaryStage.End()
}

func ContinueTargetTex() {
    Rs.PrimaryStage.Begin()
}

func BeginCustomRenderScreenspace() {
    Rs.PrimaryStage.End()
	rl.BeginTextureMode(Rs.GuiTargetTex)
}

func EndCustomRender() {
	rl.EndTextureMode()

	rl.BeginTextureMode(Rs.CeilTex)

	rl.BeginMode2D(Rs.PrimaryStage.trueSSCamera)
	rl.DrawTexturePro(Rs.PrimaryStage.Target.Texture,
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.PrimaryStage.Target.Texture.Width),
			Height: float32(Rs.PrimaryStage.Target.Texture.Height),
		},
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.CeilTex.Texture.Width),
			Height: float32(Rs.CeilTex.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	rl.EndMode2D()

    rl.BeginMode2D(Rs.FxStage.trueSSCamera)
	rl.DrawTexturePro(Rs.FxStage.Target.Texture,
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.FxStage.Target.Texture.Width),
			Height: float32(Rs.FxStage.Target.Texture.Height),
		},
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.CeilTex.Texture.Width),
			Height: float32(Rs.CeilTex.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	rl.EndMode2D()

	// also render the gui to the the ceiltex
	rl.DrawTexturePro(Rs.GuiTargetTex.Texture,
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.GuiTargetTex.Texture.Width),
			Height: float32(Rs.GuiTargetTex.Texture.Height),
		},
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.CeilTex.Texture.Width),
			Height: float32(Rs.CeilTex.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	rl.EndTextureMode()

    //renderPass(&Rs.FxTex, nil, &Rs.CeilTex, nil, 1, -1)

	// NOTE: Not sure if the shader should also render over the GUI, or if
	// we need to separate World and GUI into individual render textures.
	if settings.CurrentSettings().EnableCrtEffect {
		rl.SetShaderValue(Rs.crt_dat.Crtshader, Rs.crt_dat.Crtshader_loc_time, []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
		// TODO: add a toggle setting between the full shader with accumulate
		// and a "lite" version without accumulation
		renderPass(&Rs.crt_dat.Accumulationbuffer, nil, &Rs.crt_dat.Blurbuffer, &Rs.crt_dat.Blurshader, 1, 1)
		renderPass(&Rs.CeilTex, &Rs.crt_dat.Blurbuffer, &Rs.crt_dat.Accumulationbuffer, &Rs.crt_dat.Accumulateshader, 1, 1)
		renderPass(&Rs.CeilTex, &Rs.crt_dat.Accumulationbuffer, &Rs.crt_dat.Blendbuffer, &Rs.crt_dat.Blendshader, 1, 1)
		renderPass(&Rs.crt_dat.Blendbuffer, &Rs.crt_dat.Blurbuffer, &Rs.CeilTex, &Rs.crt_dat.Crtshader, 1, 1)
	}

	// render the oversize render texture to the actual screen.
	rl.DrawTexturePro(Rs.CeilTex.Texture,
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.CeilTex.Texture.Width),
			Height: float32(Rs.CeilTex.Texture.Height),
		},
		rl.Rectangle{
			// the position calculations for X and Y are in place, so the
			// texture is rendered in the middle of the screen, in case the
			// aspect ratio does not match
			X:      (float32(rl.GetScreenWidth()) - Rs.RenderResolution.X*Rs.MinScale) * 0.5,
			Y:      (float32(rl.GetScreenHeight()) - Rs.RenderResolution.Y*Rs.MinScale) * 0.5,
			Width:  Rs.RenderResolution.X * Rs.MinScale,
			Height: Rs.RenderResolution.Y * Rs.MinScale,
		},
		rl.Vector2{X: 0, Y: 0}, 0, rl.White,
	)

	rl.SetMouseScale(1, 1)
	rl.SetMouseOffset(0, 0)
}

func recalcScaleFactor() {
	Rs.RenderScale = rl.Vector2{
		X: float32(rl.GetScreenWidth()) / Rs.RenderResolution.X,
		Y: float32(rl.GetScreenHeight()) / Rs.RenderResolution.Y,
	}
	Rs.MinScale = float32(math.Min(
		float64(Rs.RenderScale.X),
		float64(Rs.RenderScale.Y)))
}

func renderPass(tex0 *rl.RenderTexture2D, tex1 *rl.RenderTexture2D, dest *rl.RenderTexture2D, shader *rl.Shader, mirrorX, mirrorY float32) {
	rl.BeginTextureMode(*dest)
	if shader != nil {
		rl.BeginShaderMode(*shader)
		if tex1 != nil {
			tex1_loc := Rs.tex1locs[*shader]
			rl.SetShaderValueTexture(*shader, tex1_loc, tex1.Texture)
		}
	}
	rl.DrawTexturePro(tex0.Texture,
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  mirrorX * float32(Rs.CeilTex.Texture.Width),
			Height: mirrorY * -float32(Rs.CeilTex.Texture.Height),
		},
		rl.Rectangle{
			X:      0.0,
			Y:      0.0,
			Width:  float32(Rs.CeilTex.Texture.Width),
			Height: float32(Rs.CeilTex.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0}, 0, rl.White,
	)
	if shader != nil {
		rl.EndShaderMode()
	}
	rl.EndTextureMode()
}
