package webserver

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
                    <h5 class="card-title">
                        {{ if $value.AnnotatedName }}{{ $value.AnnotatedName }}{{ else }}{{ $value.SpecName }}{{ end }}
                    </h5>
                    <div class="media">
                        {{ if $value.AnnotatedURL }}
                        <img src="http://getfavicons.com/api/?url={{ $value.AnnotatedURL }}&size=32" class="mr-3"
                            alt="...">
                        {{ end }}
                        <p class="card-text media-body">
                            <span class="badge badge-warning">{{ $value.SpecNamespace }}</span>
                            {{ $value.AnnotatedDescription }}
                        </p>
                    </div>
                    <div class="clearfix">&nbsp;</div>
                    <a href="{{ $value.SpecURL }}" class="btn btn-outline-primary">Service</a>
                    {{ if $value.AnnotatedURL }}
                    <a href="{{ $value.AnnotatedURL }}" class="btn btn-outline-secondary">About</a>
                    {{ end }}
                </div>
            </div>
        </div>
        {{ end }}
    </div>

</body>

</html>
`
