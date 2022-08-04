module github.com/HelloWorksGroup/KOOK2QQ-bot

go 1.18

require (
	github.com/Nigh/MiraiGo-Template-Mod v0.0.0-20220804125811-8dc160f3287f
	github.com/jpillora/overseer v1.1.6
	github.com/lonelyevil/khl v0.0.27
	github.com/lonelyevil/khl/log_adapter/plog v0.0.27
	github.com/phuslu/log v1.0.81
	github.com/spf13/viper v1.12.0
	local/khlcard v0.0.0-00010101000000-000000000000
	local/rt v0.0.0-00010101000000-000000000000
)

require (
	github.com/Baozisoftware/qrcode-terminal-go v0.0.0-20170407111555-c0650d8dff0f // indirect
	github.com/Mrs4s/MiraiGo v0.0.0-20220720124026-5c0e2c5773de // indirect
	github.com/RomiChan/protobuf v0.1.1-0.20220624030127-3310cba9dbc0 // indirect
	github.com/RomiChan/syncx v0.0.0-20220404072119-d7ea0ae15a4c // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/bits-and-blooms/bitset v1.3.0 // indirect
	github.com/bits-and-blooms/bloom/v3 v3.2.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/fumiama/imgsz v0.0.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/s3 v1.1.4 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/tidwall/gjson v1.14.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20220803195053-6e608f9ce704 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)

replace local/khlcard => ./kcard

replace local/rt => ./qq
