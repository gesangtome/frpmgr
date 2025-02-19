package ui

import (
	"fmt"
	"github.com/koho/frpmgr/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
)

type SimpleSection struct {
	view     *walk.Dialog
	sections []*config.Section
	title    string
	prefix   string
	types    []string
	port     uint16
}

func NewSimpleSectionDialog(title string, prefix string, types []string, port uint16) *SimpleSection {
	v := new(SimpleSection)
	v.title = title
	v.prefix = prefix
	v.types = types
	v.port = port
	v.sections = make([]*config.Section, 0)
	return v
}

func (t *SimpleSection) View() Dialog {
	var acceptPB, cancelPB *walk.PushButton
	var remotePortEdit *walk.LineEdit
	var localIPEdit *walk.LineEdit
	var db *walk.DataBinder
	icon, _ := loadLogoIcon(32)
	return Dialog{
		Icon:          icon,
		AssignTo:      &t.view,
		Title:         "添加 " + t.title,
		Layout:        VBox{},
		Font:          Font{Family: "微软雅黑", PointSize: 9},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo: &db,
			DataSource: &struct {
				Port string
				Addr string
			}{"", "127.0.0.1"},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "远程端口:"},
					LineEdit{AssignTo: &remotePortEdit, Text: Bind("Port", Regexp{"^\\d+$"})},
					Label{Text: "本地地址:"},
					LineEdit{AssignTo: &localIPEdit, Text: Bind("Addr", Regexp{".+"})},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "确定", AssignTo: &acceptPB, OnClicked: func() {
						if remotePortEdit.Text() == "" || localIPEdit.Text() == "" {
							return
						}
						for _, proto := range t.types {
							sect := config.Section{
								Name:       fmt.Sprintf("%s-%s-%s", t.prefix, proto, remotePortEdit.Text()),
								Type:       proto,
								LocalIP:    localIPEdit.Text(),
								LocalPort:  strconv.FormatInt(int64(t.port), 10),
								RemotePort: remotePortEdit.Text(),
							}
							t.sections = append(t.sections, &sect)
						}
						t.view.Accept()
					}},
					PushButton{Text: "取消", AssignTo: &cancelPB, OnClicked: func() { t.view.Cancel() }},
				},
			},
		},
	}
}

func (t *SimpleSection) Run(owner walk.Form) (int, error) {
	return t.View().Run(owner)
}
