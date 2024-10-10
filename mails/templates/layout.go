package templates

const Layout = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1" />
    <style>
        body, html {
            padding: 0;
            margin: 0;
            border: 0;
            color: #16161a;
            background: #fff;
            font-size: 14px;
            line-height: 20px;
            font-weight: normal;
            font-family: Source Sans Pro, sans-serif, emoji;
        }
        body {
            padding: 20px 30px;
        }
        strong {
            font-weight: bold;
        }
        em, i {
            font-style: italic;
        }
        p {
            display: block;
            margin: 10px 0;
            font-family: inherit;
        }
        small {
            font-size: 12px;
            line-height: 16px;
        }
        hr {
            display: block;
            height: 1px;
            border: 0;
            width: 100%;
            background: #e1e6ea;
            margin: 10px 0;
        }
        a {
            color: inherit;
        }
        .hidden {
            display: none !important;
        }
        .btn {
            display: inline-block;
            vertical-align: top;
            border: 0;
            cursor: pointer;
            color: #fff !important;
            background: #16161a !important;
            text-decoration: none !important;
            line-height: 40px;
            width: auto;
            min-width: 150px;
            text-align: center;
            padding: 0 20px;
            margin: 5px 0;
            font-family: Source Sans Pro, sans-serif, emoji;;
            font-size: 14px;
            font-weight: bold;
            border-radius: 6px;
            box-sizing: border-box;
        }
    </style>
</head>
<body>
    {{template "content" .}}
</body>
</html>
`
