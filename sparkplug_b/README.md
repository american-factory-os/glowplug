# Sparkplug B Protocol Buffer

Glowplug includes the official Sparkplug B `sparkplug_b_go.proto` file from the [Eclipse Tahu](https://github.com/eclipse/tahu/blob/master/sparkplug_b/sparkplug_b_c_sharp.proto) project, licensed under the Eclipse Public License 2.0.

To build:

```bash
protoc --proto_path=sparkplug_b --go_out=sparkplug --go_opt=paths=source_relative sparkplug_b/sparkplug_b_go.proto
```

