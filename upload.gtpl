<html>
<head>
       <title>Test Bed - Large File Upload Service</title>
</head>
<body>
<h1>Large File Upload Service</h1>
<form enctype="multipart/form-data" action="http://127.0.0.1:9090/upload" method="post">
    <input type="file" name="uploadFile" /><br/><br/>
    <input type="checkbox" name="private" value="private" checked />
    <label for="private">Private upload (obfuscate filename)</label><br/><br/>
    <input type="submit" value="upload" />
</form>
</body>
</html>