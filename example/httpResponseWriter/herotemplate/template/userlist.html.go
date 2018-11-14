//由hero生成的代码
//不要编辑
package template

import (
	"bytes"
	"github.com/shiyanhui/hero"
)

func UserList(userList []string, buffer *bytes.Buffer) {
	buffer.WriteString(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>
    <body>
        `)
	for _, user := range userList {
		buffer.WriteString(`
        <ul>
            `)
		buffer.WriteString(`<li>
    `)
		hero.EscapeHTML(user, buffer)
		buffer.WriteString(`
</li>
`)
		buffer.WriteString(`
        </ul>
    `)
	}
	buffer.WriteString(`
    </body>
</html>
`)
}