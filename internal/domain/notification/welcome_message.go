package notification

var WelcomeMessage = `<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; padding: 20px;">

<h2 style="color:#4A90E2;">Hello World from HitPipe! 💙</h2>

<p>Это HTML-письмо с <b>жирным</b> текстом, <span style="color:red;">цветом</span> и ссылками.</p>

<img src="https://picsum.photos/300" 
     alt="Random Image"
     style="border-radius:10px; margin: 20px 0;">

<a href="http://localhost:8080/verify?token=%s" 
   style="display:inline-block; padding:12px 20px; background:#4A90E2; color:white; text-decoration:none; border-radius:8px;">
   Перейти на сайт
</a>

</body>
</html>`
