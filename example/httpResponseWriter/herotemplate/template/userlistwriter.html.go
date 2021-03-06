//由hero生成的代码
//不要编辑
package template

import (
	"io"
	"github.com/shiyanhui/hero"
)

func UserListToWriter(userList []string, w io.Writer) (int, error) {
	_buffer := hero.GetBuffer()
	defer hero.PutBuffer(_buffer)
	_buffer.WriteString(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>

    <body>
        `)
	for _, user := range userList {
		_buffer.WriteString(`
        <ul>
            `)
		_buffer.WriteString(`<li>
    `)
		hero.EscapeHTML(user, _buffer)
		_buffer.WriteString(`
</li>
`)
		_buffer.WriteString(`
        </ul>
    `)
	}
	_buffer.WriteString(`
    </body>
</html>
`)
	return w.Write(_buffer.Bytes())
}