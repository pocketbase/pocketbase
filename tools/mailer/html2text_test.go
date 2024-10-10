package mailer

import (
	"testing"
)

func TestHTML2Text(t *testing.T) {
	scenarios := []struct {
		html     string
		expected string
	}{
		{
			"",
			"",
		},
		{
			"ab  c",
			"ab c",
		},
		{
			"<!-- test html comment -->",
			"",
		},
		{
			"<!-- test html comment -->   a   ",
			"a",
		},
		{
			"<span>a</span>b<span>c</span>",
			"abc",
		},
		{
			`<a href="a/b/c">test</span>`,
			"[test](a/b/c)",
		},
		{
			`<a href="">test</span>`,
			"[test]",
		},
		{
			"<span>a</span>  <span>b</span>",
			"a b",
		},
		{
			"<span>a</span>   b   <span>c</span>",
			"a b c",
		},
		{
			"<span>a</span>   b   <div>c</div>",
			"a b \r\nc",
		},
		{
			`
				<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
				<html xmlns="http://www.w3.org/1999/xhtml">
				<head>
				    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
				    <meta name="viewport" content="width=device-width,initial-scale=1" />
				    <style>
				        body {
				            padding: 0;
				        }
				    </style>
				</head>
				<body>
					<!-- test html comment -->
					<style>
					    body {
					        padding: 0;
					    }
					</style>
				    <div class="wrapper">
				        <div class="content">
				            <p>Lorem ipsum</p>
				            <p>Dolor sit amet</p>
				            <p>
				            	<a href="a/b/c">Verify</a>
				            </p>
				            <br>
				            <p>
				            	<a href="a/b/c"><strong>Verify2.1</strong> <strong>Verify2.2</strong></a>
				            </p>
				            <br>
				            <br>
				            <div>
				            	<div>
				            		<div>
							            <ul>
							            	<li>ul.test1</li>
							            	<li>ul.test2</li>
							            	<li>ul.test3</li>
							            </ul>
							            <ol>
							            	<li>ol.test1</li>
							            	<li>ol.test2</li>
							            	<li>ol.test3</li>
							            </ol>
				            		</div>
				            	</div>
				            </div>
				            <select>
				            	<option>Option 1</option>
				            	<option>Option 2</option>
				            </select>
				            <textarea>test</textarea>
				            <input type="text" value="test" />
				            <button>test</button>
				            <p>
				                Thanks,<br/>
				                PocketBase team
				            </p>
				        </div>
				    </div>
				</body>
				</html>
			`,
			"Lorem ipsum \r\nDolor sit amet \r\n[Verify](a/b/c)  \r\n[Verify2.1 Verify2.2](a/b/c)  \r\n\r\n- ul.test1 \r\n- ul.test2 \r\n- ul.test3  \r\n- ol.test1 \r\n- ol.test2 \r\n- ol.test3         \r\nThanks,\r\nPocketBase team",
		},
	}

	for i, s := range scenarios {
		result, err := html2Text(s.html)
		if err != nil {
			t.Errorf("(%d) Unexpected error %v", i, err)
		}

		if result != s.expected {
			t.Errorf("(%d) Expected \n(%q)\n%v,\n\ngot:\n\n(%q)\n%v", i, s.expected, s.expected, result, result)
		}
	}
}
