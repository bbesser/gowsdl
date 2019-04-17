// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name | makePublic | addNamespacePrefix}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$type}} []{{toGoType .List.ItemType }}
	{{else if ne .Union.MemberTypes ""}}
		type {{$type}} string
	{{else if .Union.SimpleType}}
		type {{$type}} string
	{{else}}
		type {{$type}} {{toGoType .Restriction.Base}}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic | addNamespacePrefix}} {{$type}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{ else }}
			{{ .Name | makeFieldPublic}} string ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}}{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Name | makePublic | addNamespacePrefix}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` // foo3
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic | addNamespacePrefix}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Ref | toGoType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` // foo2
				{{else}}
					{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` // foo1
				{{end}}
			{{else}}
				{{template "ComplexTypeInline" .}}
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceReservedWords .Name | makeFieldPublic | addNamespacePrefix}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` // foo5 {{end}}
		{{end}}
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}
	{{ $ignored := setCurrentTargetNamespace $targetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{$name := .Name}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{with .ComplexType}}
				type {{$name | replaceReservedWords | makePublic | addNamespacePrefix}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + ` // bar1
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
					{{end}}
				}
			{{end}}
		{{else}}
			{{ $nsName := $name | replaceReservedWords | makePublic }}
			{{ $nsType := toGoType .Type | removePointerFromType }}
			type {{$nsName}} struct {
				XMLName xml.Name ` + "`xml:\"web:{{$nsName}}\"`" + `
				XmlNS   string   ` + "`xml:\"xmlns:web,attr\"`" + `
				{{$nsType}} ` + "`xml:\"{{$nsName}}\"`" + `
			}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic | addNamespacePrefix}}
		type {{$name}} struct {
			{{$typ := findNameByType .Name}}
			{{$typWithNamespace := findNameByType .Name | addNamespacePrefix}}
			{{if ne $name $typWithNamespace}}
				// XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$typ}}\"`" + ` // bar2 {{ $typWithNamespace }}
			{{end}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
				{{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
			{{end}}
		}
	{{end}}
{{end}}
`
