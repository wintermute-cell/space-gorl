<!-- LTeX: language=en-US -->
# Scripts
> [!NOTE]
> Scripts are generally not part of the shipped game and only intended to aid
> during development. They will not be bug-free, and their output should always
> be checked for correctness.

## sync_settings.py
Based on the `GameSettings` struct found in `settings.go`, edits files to
generate fallback settings in `settings.go` and fills out the `settings.json`
file in `assets`. 

The values written as comments are used as default values:
```go
TargetFps    int  `json:"targetFps"`    // 144
```

Usage example:
```bash
python scripts/sync_settings.py [path/to/settings.go path/to/settings.json]
```

If no path is given, the following defaults apply: 
```
./pkg/settings/settings.go
./assets/settings.json
```

## create_scene.sh & create_entity.sh
Please see [Using scripts for scene and entity creation](https://github.com/wintermute-cell/cowboy-gorl/blob/gui/documentation/creating-scenes-and-entities.md#recommended-using-scripts-for-creation)
