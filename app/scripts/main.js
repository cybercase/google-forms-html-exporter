$(function () {
    // enable bootstrap tooltips
    $('[data-toggle="tooltip"]').tooltip({trigger: 'manual'})
})

$('#gifdiv').click(function () {
    if ($(this).find('img').attr('src') == 'images/preview.jpg') {
        $(this).find('img').attr('src', 'images/preview.gif');
    } else {
        $(this).find('img').attr('src', 'images/preview.jpg');
    }
});

$('#example-action').click((event) => {
    event.preventDefault()
    let exampleUrl = $('#example-url').attr('href');
    $('#input-url').val(exampleUrl)
})

$('#input-url').keypress((event) => {
    if (event.keyCode === 13) {
        $('#button-fetch').click()
    }
})

$('#button-fetch').click((event) => {
    event.stopPropagation()

    let $inputUrl = $('#input-url')
    function showTooltip(message) {
        try {
        $inputUrl
            .focus()
            .attr('title', message)
            .tooltip('fixTitle')
            .tooltip('show')
        setTimeout(() => $inputUrl.tooltip('hide'), 2500)
        } catch (e) { alert(e)}
    }

    function setCodeAsInnerText(selector, text) {
        text = text
            .split('\n')
            .filter(l => !!l.trim())  // Remove whitelines
            .map(l => l.trim() === '<!-- emptyline -->' ? '' : l)  // Transform comment to emptyline
            .join('\n')

        let $el = $(selector)
        $el.text(text)
        hljs.highlightBlock($el.get(0))
    }

    let url = $inputUrl.val()

    if (!url) {
        return showTooltip('Don\'t forget the URL!')
    }

    let $btnFetch = $('#button-fetch')
    $btnFetch.button('loading')

    $.get(`${window.config.serverAddress}/formdress?url=${url}`)
    .fail((response) => {
        if (response.status === 0) {
            showTooltip('Sorry, service is unavailable at the moment')
        }
        else if (response.responseJSON) {
            showTooltip(response.responseJSON.Error)
        }
    })
    .done((context) => {
        // Bootstrap Form
        let bootstrapCodeForm = bootstrapForm(context)
        setCodeAsInnerText('#target-bootstrap-html', bootstrapCodeForm)

        let bootstrapCodeJs = bootstrapJs(context)
        setCodeAsInnerText('#target-bootstrap-js', bootstrapCodeJs)

        // Exec Bootstrap JS
        $('#target-demo').html(bootstrapCodeForm)
        eval(bootstrapCodeJs)

        $('#main-area').removeClass('hidden')
        $('.marketing-area').addClass('hidden')
    })
    .always(() => $btnFetch.button('reset'))

    return;
})
