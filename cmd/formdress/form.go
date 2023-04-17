package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FieldType int

const (
	FieldShort      FieldType = 0
	FieldParagraph            = 1
	FieldChoices              = 2
	FieldDropdown             = 3
	FieldCheckboxes           = 4
	FieldLinear               = 5
	FieldTitle                = 6
	FieldGrid                 = 7
	FieldSection              = 8
	FieldDate                 = 9
	FieldTime                 = 10
	FieldImage                = 11
	FieldVideo                = 12
	FieldUpload               = 13
)

type Widget map[string]interface{}
type Option map[string]interface{}

type Field struct {
	ID     int       `json:"id"`
	Label  string    `json:"label"`
	Desc   string    `json:"desc"`
	TypeID FieldType `json:"typeid"`

	Widgets []Widget `json:"widgets"`
}

type Fields []Field

type Form struct {
	Title        string `json:"title"`
	Header       string `json:"header"`
	Desc         string `json:"desc"`
	Path         string `json:"path"`
	Action       string `json:"action"`
	Fbzx         string `json:"fbzx"`
	SectionCount int    `json:"sectionCount"`
	AskEmail     bool   `json:"askEmail"`

	Fields Fields `json:"fields"`
}

func toInt(i interface{}) int {
	integer, ok := i.(int)
	if ok {
		return integer
	}

	number, ok := i.(json.Number)
	if n, err := number.Int64(); ok && err == nil {
		return int(n)
	}

	return 0
}

func toString(i interface{}) string {
	number, ok := i.(json.Number)
	if ok {
		return number.String()
	}

	s, ok := i.(string)
	if ok {
		return s
	}

	return ""
}

func toBool(i interface{}) bool {
	boolean, ok := i.(bool)
	if ok {
		return boolean
	}

	number, ok := i.(json.Number)
	if n, err := number.Int64(); ok && err == nil {
		return n != 0
	}

	return false
}

func toSlice(i interface{}) []interface{} {
	slice, ok := i.([]interface{})
	if ok {
		return slice
	}
	return nil
}

func NewFieldFromData(data []interface{}) Field {
	f := Field{
		ID:     toInt(data[0]),
		Label:  toString(data[1]),
		Desc:   toString(data[2]),
		TypeID: FieldType(toInt(data[3])),
	}

	switch f.TypeID {
	case FieldShort:
		fallthrough
	case FieldParagraph:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),
		}}

	case FieldChoices:
		fallthrough
	case FieldCheckboxes:
		fallthrough
	case FieldDropdown:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		options := toSlice(widget[1])

		opts := []Option{}
		for _, opt := range options {
			o := toSlice(opt)

			option := Option{
				"label": toString(o[0]),
			}
			// Handle the case for missing information in option object
			if len(o) > 2 {
				option["href"] = toString(o[2])
			}
			if len(o) > 4 {
				option["custom"] = toBool(o[4])
			}

			opts = append(opts, option)
		}

		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),
			"options":  opts,
		}}

	case FieldLinear:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		legend := toSlice(widget[3])
		options := toSlice(widget[1])

		opts := []Option{}
		for _, opt := range options {
			o := toSlice(opt)
			opts = append(opts, Option{
				"label": toString(o[0]),
			})
		}

		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),
			"options":  opts,
			"legend": Option{
				"first": toString(legend[0]),
				"last":  toString(legend[1]),
			},
		}}

	case FieldGrid:
		widgets := toSlice(data[4])
		f.Widgets = []Widget{}
		for _, widget := range widgets {
			w := toSlice(widget)
			columns := toSlice(w[1])

			cols := []Option{}
			for _, col := range columns {
				c := toSlice(col)
				cols = append(cols, Option{"label": c[0]})
			}
			f.Widgets = append(f.Widgets, Widget{
				"id":       toString(w[0]),
				"required": toBool(w[2]),
				"name":     toString(toSlice(w[3])[0]),

				"columns": cols,
			})
		}
	case FieldDate:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		options := toSlice(widget[7])

		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),

			"options": Option{
				"time": toBool(options[0]),
				"year": toBool(options[1]),
			},
		}}
	case FieldTime:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		options := toSlice(widget[6])

		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),

			"options": Option{
				"duration": toBool(options[0]),
			},
		}}

	case FieldVideo:
		extra := toSlice(data[6])
		opts := toSlice(extra[2])
		f.Widgets = []Widget{{
			"id": toString(extra[0]),
			"res": Option{
				"w":        toInt(opts[0]),
				"h":        toInt(opts[1]),
				"showText": toBool(opts[2]),
			},
		}}

	case FieldImage:
		extra := toSlice(data[6])
		opts := toSlice(extra[2])
		f.Widgets = []Widget{{
			"id": toString(extra[0]),
			"res": Option{
				"w":        toInt(opts[0]),
				"h":        toInt(opts[1]),
				"showText": f.Desc != "",
			},
		}}

	case FieldUpload:
		widgets := toSlice(data[4])
		widget := toSlice(widgets[0])
		options := toSlice(widget[10])
		f.Widgets = []Widget{{
			"id":       toString(widget[0]),
			"required": toBool(widget[2]),
			"options": Option{
				"types":          toSlice(options[1]),
				"maxUploads":     toInt(options[2]),
				"maxSizeInBytes": toString(options[3]),
			},
		}}

	case FieldSection:
		fallthrough
	case FieldTitle:

	}
	return f
}

func NewFieldsFromData(data []interface{}) Fields {
	f := make(Fields, 0, 0)
	for _, d := range data {
		field := NewFieldFromData(toSlice(d))
		f = append(f, field)
	}
	return f
}

func (f *Form) UnmarshalJSON(b []byte) error {
	data := make([]interface{}, 0, 0)
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	f.Title = toString(data[3])
	f.Path = toString(data[2])
	f.Action = toString(data[14])

	extraData := toSlice(data[1])

	f.Fields = NewFieldsFromData(toSlice(extraData[1]))
	f.SectionCount = 1
	for _, field := range f.Fields {
		if field.TypeID == FieldSection {
			f.SectionCount++
		}
	}

	f.Desc = toString(extraData[0])
	f.Header = toString(extraData[8])

	if otherExtraData := toSlice(extraData[10]); otherExtraData != nil && len(otherExtraData) > 4 {
		f.AskEmail = toInt(otherExtraData[4]) == 1
	}

	return nil
}

var InvalidForm = errors.New("Invalid Form")

func ExtractImages(doc *goquery.Document, form *Form) {
	// Fetch images src from the document
	for _, w := range form.Fields {
		if w.TypeID == FieldImage {
			imgEl := doc.Find("[data-item-id='" + strconv.Itoa(w.ID) + "'] img")
			w.Widgets[0]["src"] = imgEl.AttrOr("src", "")
		} else if w.TypeID == FieldVideo {
			iframeEl := doc.Find("[data-item-id='" + strconv.Itoa(w.ID) + "'] iframe")
			w.Widgets[0]["src"] = iframeEl.AttrOr("src", "")
		}

	}
}

func FormExtract(content io.Reader) (*Form, error) {
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return nil, err
	}

	var script *goquery.Selection
	_ = doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "var FB_PUBLIC_LOAD_DATA_") {
			script = s
			return false
		}
		return true
	})

	if script == nil {
		return nil, InvalidForm
	}

	fbzx, ok := doc.Find("[name=\"fbzx\"]").Attr("value")
	if !ok {
		return nil, InvalidForm
	}

	s := script.Text()
	s = strings.Replace(s, "var FB_PUBLIC_LOAD_DATA_ =", "", -1)
	s = strings.Replace(s, ";", "", -1)
	s = strings.TrimSpace(s)

	form := &Form{}
	json.Unmarshal([]byte(s), form)
	form.Fbzx = fbzx

	ExtractImages(doc, form)

	return form, nil
}
