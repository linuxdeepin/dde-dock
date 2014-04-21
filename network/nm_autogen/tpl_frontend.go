package main

const tplFrontEnd = `
import QtQuick 2.1
import Deepin.Widgets 1.0

BaseEditSection {
    section: "{{.Name}}"
    
    header.sourceComponent: EditDownArrowHeader{
        text: dsTr("{{.DisplayName}}")
    }

    content.sourceComponent: Column {
        {{genFrontEndWidgetInfo $key.Type $key.Value}}
    }
}
`
