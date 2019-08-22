package web

const port = "9000"

const htmlTemplate = `
<!doctype html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <title></title>
</head>

<body>

    <div class="container">
        {{ range $key, $value := . }}
        <div class="row justify-content-md-center">
            <div class="card col-md-8">
                <div class="card-body">
                    <span class="badge badge-info">{{ $value.SpecNamespace }}</span>
                    <span class="badge badge-primary">{{ $value.SpecName }}</span>
                    <span class="badge badge-secondary">{{ $value.SpecType }}</span>
                    <div class="clearfix">&nbsp;</div>
                    <h5 class="card-title">
                        {{ if $value.AnnotatedName }}{{ $value.AnnotatedName }}{{ else }}{{ $value.SpecName }}{{ end }}
                    </h5>
                    <div class="media">
                        {{ if $value.AnnotatedURL }}
                        <img src="http://getfavicons.com/api/?url={{ $value.AnnotatedURL }}&size=32" class="mr-3"
                            alt="...">
                        {{ end }}
                        <p class="card-text media-body">
                            {{ $value.AnnotatedDescription }}
                            {{ if $value.AnnotatedURL }}
                            <a
                                href="{{ $value.AnnotatedURL }}" target="_blank">...</a>
                            {{ end }}
                        </p>
                    </div>
                    <div class="clearfix">&nbsp;</div>
                    {{ range $url := $value.SpecURL }}
                    <a href="{{ $url }}"
                        class="btn btn-link" target="_blank">{{ $url }}</a>
                    {{ end }}
                </div>
            </div>
        </div>
        <div class="clearfix">&nbsp;</div>
        {{ end }}
    </div>

</body>

</html>
`
