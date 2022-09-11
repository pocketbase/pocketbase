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
            padding: 30px 20px;
        }
        p {
            display: block;
            margin: 10px 0;
            font-family: Source Sans Pro, sans-serif, emoji;
        }
        strong {
            font-weight: bold;
        }
        em, i {
            font-style: italic;
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
            line-height: 45px;
            width: 100%;
            max-width: 100%;
            text-align: center;
            padding: 0 30px;
            margin: 10px 0;
            font-family: Source Sans Pro, sans-serif, emoji;;
            font-size: 14px;
            font-weight: bold;
            border-radius: 3px;
            box-sizing: border-box;
        }
        .wrapper {
            display: block;
            width: 460px;
            max-width: 100%;
            margin: auto;
            font-family: inherit;
        }
        .content {
            display: block;
            width: 100%;
            padding: 10px 20px;
            font-family: inherit;
            box-sizing: border-box;
            background: #fff;
            border-radius: 10px;
            -webkit-box-shadow: 0px 2px 30px 0px rgb(0,0,0,0.05);
            -moz-box-shadow: 0px 2px 30px 0px rgb(0,0,0,0.05);
            box-shadow: 0px 2px 30px 0px rgb(0,0,0,0.05);
        }
        .footer {
            display: block;
            width: 100%;
            text-align: center;
            margin: 10px 0;
            color:  #666f75;
            font-size: 11px;
            font-family: inherit;
        }
    </style>
</head>
<body>
    <div class="wrapper">
        <div class="content">
            {{template "content" .}}
        </div>
    </div>
</body>
</html>
`
