module github.com/rubensayshi/yubitoast

go 1.12

require (
	fyne.io/fyne v1.2.2
	github.com/deckarep/gosx-notifier v0.0.0-20180201035817-e127226297fb // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20191125211704-12ad95a8df72
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/hpcloud/tail v1.0.0
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8
	github.com/kr/pretty v0.2.0 // indirect
	github.com/martinlindhe/notify v0.0.0-20181008203735-20632c9a275a
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/pkg/errors v0.8.1
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	golang.org/x/crypto v0.0.0-20200214034016-1d94cc7ab1c6 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace fyne.io/fyne => ../fyne

replace github.com/go-gl/glfw/v3.3/glfw => ../glfw/v3.3/glfw
