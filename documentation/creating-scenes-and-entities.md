<!-- LTeX: language=en-US -->
# How to create new Scenes and Entities
For both `scenes` and `entities`, there is a template file:

```bash
pkg/scenes/scene_template.go
pkg/entities/entity_template.go
```

To create a new `scene` or `entity`, copy over the template to a new file:
```bash
cp pkg/scenes/scene_template.go pkg/scenes/scene_my_new_scene.go
```

And replace every occurrence of the word `Template` in that new file.

## (Recommended) Using scripts for creation
Since the whole copy and replace procedure can get repetitive after a while,
there exist two scripts, `scripts/create_scene.sh` and `scripts/create_entity.sh`.

### For scenes:
```bash
./scripts/create_scene.sh <scene_name_in_snake_case>
```
This command will generate `scene_<scene_name_in_snake_case>.go` under
`pkg/scenes`, and automatically fill in the scene name in the new file.

### For entities:
```bash
./scripts/create_entity.sh <entity_name_in_snake_case>
```
This will create `entity_<entity_name_in_snake_case>.go` in the `pkg/entities`
directory, and automatically fill in the entity name in the new file.

> [!NOTE]
> Always ensure you use snake_case for naming as shown in the placeholders above.
