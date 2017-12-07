const FieldTypes = [
    'short', // 0
    'paragraph', // 1
    'choices', // 2
    'dropdown', // 3
    'checkboxes', // 4
    'linear', // 5
    'title', // 6
    'grid', // 7
    'section', // 8
    'date', // 9
    'time', // 10
    'image', // 11
    'video', // 12
]

Handlebars.registerHelper('ifType', function(field, typename, options) {
    if (field.typeid === FieldTypes.indexOf(typename)) {
        return options.fn(this)
    }
    return ''
})

Handlebars.registerHelper('fieldtype', function(typeid) {
    return FieldTypes[typeid]
})

Handlebars.registerHelper('legend', function(array, last) {
    if (last) {
        return array[array.length-1].label
    } else {
        return array[0].label
    }
})

Handlebars.registerHelper('datePlaceholder', function () {
    console.log('EHI')
    return new Date().toLocaleDateString()
})

Handlebars.registerHelper('timePlaceholder', function () {
    return new Date().toLocaleTimeString()
})

let bootstrapForm = Handlebars.compile(`
<form action="https://docs.google.com{{path}}/d/{{action}}/formResponse"
      target="_self"
      id="bootstrapForm"
      method="POST">

    {{#if title}}
    <fieldset>
        <h2>{{ title }}<br><small>{{desc}}</small></h2>
    </fieldset>
    {{/if}}
    {{#each fields as |f|}}
    <!-- emptyline -->
    <!-- emptyline -->
    <!-- Field type: "{{ fieldtype f.typeid }}" id: "{{f.id}}" -->
    <fieldset>
        <legend for="{{f.id}}">{{f.label}}</legend>
        <div class="form-group">
            {{#if f.desc}}
            <p class="help-block">{{ f.desc }}</p>
            {{/if}}

            {{#ifType f 'short'}}
            <input id="{{f.widgets.0.id}}" type="text" name="entry.{{f.widgets.0.id}}" class="form-control" {{#if f.widgets.0.required}}required{{/if}}>
            {{/ifType}}

            {{#ifType f 'paragraph'}}
            <textarea id="{{f.widgets.0.id}}" name="entry.{{f.widgets.0.id}}" class="form-control" {{#if f.widgets.0.required}}required{{/if}}></textarea>
            {{/ifType}}

            {{#ifType f 'choices'}}
            {{#each f.widgets.0.options as |c|}}
            <div class="radio">
                {{#if c.custom}}
                <label>
                    <input type="radio" name="entry.{{f.widgets.0.id}}" value="__other_option__" {{#if f.widgets.0.required}}required{{/if}}>
                </label>
                <input type="text" name="entry.{{f.widgets.0.id}}.other_option_response" placeholder="custom value">
                {{else}}
                <label>
                    <input type="radio" name="entry.{{f.widgets.0.id}}" value="{{c.label}}" {{#if f.widgets.0.required}}required{{/if}}>
                    {{c.label}}
                </label>
                {{/if}}
            </div>
            {{/each}}
            {{/ifType}}

            {{#ifType f 'checkboxes'}}
            {{#each f.widgets.0.options as |c|}}
            <div class="checkbox">
                {{#if c.custom}}
                <label>
                    <input type="checkbox" name="entry.{{f.widgets.0.id}}" value="__other_option__" {{#if f.widgets.0.required}}required{{/if}}>
                </label>
                <input type="text" name="entry.{{f.widgets.0.id}}.other_option_response" placeholder="custom value">
                {{else}}
                <label>
                    <input type="checkbox" name="entry.{{f.widgets.0.id}}" value="{{c.label}}" {{#if f.widgets.0.required}}required{{/if}}>
                    {{c.label}}
                </label>
                {{/if}}
            </div>
            {{/each}}
            {{/ifType}}

            {{#ifType f 'dropdown'}}
            <select id="{{f.id}}" name="entry.{{f.widgets.0.id}}" class="form-control">
                {{#unless f.widgets.0.required}}
                <option value=""></option>
                {{/unless}}
                {{#each f.widgets.0.options as |c|}}
                <option value="{{c.label}}">{{c.label}}</option>
                {{/each}}
            </select>
            {{/ifType}}

            {{#ifType f 'linear'}}
            <div>
            {{#each f.widgets.0.options as |c|}}
            <label class="radio-inline">
                <input type="radio" name="entry.{{f.widgets.0.id}}" value="{{c.label}}" {{#if f.widgets.0.required}}required{{/if}}>
                {{c.label}}
            </label>
            {{/each}}
            </div>
            <div>
                <div>{{ legend f.widgets.0.options 0 }}: {{ f.widgets.0.legend.first }}</div>
                <div>{{ legend f.widgets.0.options 1 }}: {{ f.widgets.0.legend.last }}</div>
            </div>
            {{/ifType}}

            {{#ifType f 'grid'}}
            {{#each f.widgets as |w|}}
            <div>
                <span>{{w.name}}: </span>
                {{#each columns as |c|}}
                <label class="radio-inline">
                    <input type="radio" name="entry.{{w.id}}" value="{{c.label}}" {{#if w.required}}required{{/if}}>
                    {{c.label}}
                </label>
                {{/each}}
            </div>
            {{/each}}
            {{/ifType}}

            {{#ifType f 'title' }}
            {{/ifType}}

            {{#ifType f 'section' }}
            {{/ifType}}

            {{#ifType f 'date' }}
            <input type="date" id="{{f.widgets.0.id}}_date" placeholder="{{ datePlaceholder }}" class="form-control" {{#if f.widgets.0.required}}required{{/if}}>
            {{#if f.widgets.0.options.time}}
            <input type="time" id="{{f.widgets.0.id}}_time" placeholder="{{ timePlaceholder }}" class="form-control" {{#if f.widgets.0.required}}required{{/if}}>
            {{/if}}
            {{/ifType}}

            {{#ifType f 'time'}}
            <input type="time" id="{{f.widgets.0.id}}" placeholder="{{ timePlaceholder }}" class="form-control" {{#if f.widgets.0.required}}required{{/if}}>
            {{/ifType}}

            {{#ifType f 'image'}}
            {{#if f.widgets.0.src}}
                <img src="{{f.widgets.0.src}}" style="max-width: 100%;">
            {{/if}}
            {{/ifType}}

            {{#ifType f 'video'}}
            {{#if f.widgets.0.src}}
                <iframe src="{{f.widgets.0.src}}" style="width: 320px; height: 180px;"></iframe>
            {{/if}}
            {{/ifType}}

        </div>
    </fieldset>
    {{/each}}

    <!-- emptyline -->
    <input type="hidden" name="fvv" value="1">
    <input type="hidden" name="fbzx" value="{{fbzx}}">

    <!-- emptyline -->
    <input class="btn btn-primary" type="submit" value="Submit">
</form>
`)

let bootstrapJs = Handlebars.compile(`
// This script requires jQuery and jquery-form plugin
// You can use these ones from Cloudflare CDN:
// <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.1/jquery.min.js" integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4=" crossorigin="anonymous"></script>
// <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.form/4.2.2/jquery.form.min.js" integrity="sha256-2Pjr1OlpZMY6qesJM68t2v39t+lMLvxwpa8QlRjJroA=" crossorigin="anonymous"></script>
//
$('#bootstrapForm').submit(function (event) {
    event.preventDefault()
    var extraData = {}
    {{#each fields as |f|}}
    {{#ifType f 'date'}}
    {
        /* Parsing input date id={{f.widgets.0.id}} */
        var dateField = $("#{{f.widgets.0.id}}_date").val()
        var timeField = $("#{{f.widgets.0.id}}_time").val()
        let d = new Date(dateField)

        if (!isNaN(d.getTime())) {
            extraData["entry.{{f.widgets.0.id}}_year"] = d.getFullYear()
            extraData["entry.{{f.widgets.0.id}}_month"] = d.getMonth() + 1
            extraData["entry.{{f.widgets.0.id}}_day"] = d.getUTCDate()
        }

        if (timeField && timeField.split(':').length >= 2) {
            let values = timeField.split(':')
            extraData["entry.{{f.widgets.0.id}}_hour"] = values[0]
            extraData["entry.{{f.widgets.0.id}}_minute"] = values[1]
        }
    }
    {{/ifType}}

    {{#ifType f 'time'}}
    {
        // Parsing input time id={{f.widgets.0.id}}
        var field = $("#{{f.widgets.0.id}}").val()
        if (field) {
            var values = field.split(':')
            extraData["entry.{{f.widgets.0.id}}_hour"] = values[0]
            extraData["entry.{{f.widgets.0.id}}_minute"] = values[1]
            extraData["entry.{{f.widgets.0.id}}_second"] = values[2]
        }
    }
    {{/ifType}}
    {{/each}}

    $('#bootstrapForm').ajaxSubmit({
        data: extraData,
        error: function () {
            // Google Docs won't allow reading the response because of CORS, so this is handled as a failure.
            alert('Form Submitted. Thanks.')
            // You can also redirect the user to a custom thank-you page:
            // window.location = 'http://www.mydomain.com/thankyoupage.html'
        }
    })
})
`)

