# oauth使用手册

在oauth服务搭建好之后，使用步骤如下：
例如：后端地址为http://127.0.0.1:1102
     前端地址为http://127.0.0.1:1101

1. 第三方页面，携带我们颁发的client_id等信息，跳转到我们的前端，例如：http://127.0.0.1:1101/oauth?client_id=c83b3a16&response_type=code&scope=read&redirect_uri=http://localhost:19090/oauth/redirect
2. 用户在我们的前端输入账号密码，随后向后端的oauth login接口发起请求，例如：http://127.0.0.1:1102/oauth/login，body体中携带账号密码
3. 后端账号密码等信息校验通过后，将userid存放在请求头中，然后返回请求成功信息，
4. 前端受到成功信息后，随即携带第三方页面的client_id等信息，向后端的authorize接口发起请求，例如：http://127.0.0.1:1102/authorize?client_id=c83b3a16&response_type=code&scope=read&redirect_uri=http://localhost:19090/oauth/redirect
5. 后端authorize接口校验相关信息后，会调转到第三方页面提供的redirect_url的地址，并且加上一个code（授权码），例如：http://localhost:19090/oauth/redirect?code=SYPIGSWSNJKS2KEDDZ07KG
6. 第三方页面拿到code后，在其后端，携带code，client_id，client_secret等信息，向我方后端的token接口发起请求，例如：http://127.0.0.1:1102/token?client_id=c83b3a16&code=SYPIGSWSNJKS2KEDDZ07KG&client_secret=3a19c419&scope=read&redirect_uri=http://localhost:19090/oauth/redirect&grant_type=authorization_code
7. token接口在验证信息通过后，就会返回一段json数据，也就是令牌等信息
8. 接下来第三方页面只需要携带这个令牌，向我方发起请求，我方就能从令牌中获取到身份信息，从而响应



注意：

1. 第三方页面的redirect_url参数，必须和在我方注册的主页域名一直，否则无法通过